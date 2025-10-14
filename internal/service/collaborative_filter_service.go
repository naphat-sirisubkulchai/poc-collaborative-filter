package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"poc-collaborative-filter/internal/config"
	"poc-collaborative-filter/internal/domain"
	"poc-collaborative-filter/internal/repository"
	"poc-collaborative-filter/internal/utils"

	"go.uber.org/zap"
)

type CollaborativeFilterService interface {
	CalculateSimilarity(customer1, customer2 *domain.Customer) *domain.SimilarityScore
	FindSimilarCustomers(ctx context.Context, customerID string, limit int) ([]*SimilarCustomer, error)
	GenerateRecommendations(ctx context.Context, customerID string, recType domain.RecommendationType) ([]*domain.Recommendation, error)
}

type collaborativeFilterService struct {
	customerRepo       repository.CustomerRepository
	recommendationRepo repository.RecommendationRepository
	cfg                *config.Config
	strategySelector   *domain.StrategySelector
	diversifier        *Diversifier
	logger             *zap.Logger
}

type SimilarCustomer struct {
	Customer *domain.Customer
	Score    *domain.SimilarityScore
}

func NewCollaborativeFilterService(
	customerRepo repository.CustomerRepository,
	recommendationRepo repository.RecommendationRepository,
	cfg *config.Config,
	logger *zap.Logger,
) CollaborativeFilterService {
	return &collaborativeFilterService{
		customerRepo:       customerRepo,
		recommendationRepo: recommendationRepo,
		cfg:                cfg,
		strategySelector:   domain.NewStrategySelector(),
		diversifier:        NewDiversifier(domain.DefaultRecommendationConfig()),
		logger:             logger,
	}
}

// CalculateSimilarity computes similarity score between two customers based on:
// - Personal information (age, occupation, income, etc.)
// - Survey data (investment preferences, sectors, etc.)
// - Activity tracking (searches, views, transactions, etc.)
func (s *collaborativeFilterService) CalculateSimilarity(customer1, customer2 *domain.Customer) *domain.SimilarityScore {
	weights := domain.Weights{
		Personal:       domain.PersonalWeight,
		Survey:         domain.SurveyWeight,
		ActivityAccess: domain.ActivityAccessWeight,
		Transaction:    domain.TransactionWeight,
	}

	personalSim := s.calculatePersonalSimilarity(customer1, customer2)
	surveySim := s.calculateSurveySimilarity(customer1, customer2)
	activitySim := s.calculateActivitySimilarity(customer1, customer2)
	transactionSim := s.calculateTransactionSimilarity(customer1, customer2)

	totalScore := (personalSim * weights.Personal) +
		(surveySim * weights.Survey) +
		(activitySim * weights.ActivityAccess) +
		(transactionSim * weights.Transaction)

	return &domain.SimilarityScore{
		Customer1ID: customer1.ID,
		Customer2ID: customer2.ID,
		Score:       totalScore,
		Weights:     weights,
	}
}

// calculatePersonalSimilarity compares personal attributes
func (s *collaborativeFilterService) calculatePersonalSimilarity(c1, c2 *domain.Customer) float64 {
	score := 0.0
	matches := 0.0
	total := 0.0

	// Age similarity (within 10 years = high similarity)
	age1 := s.getNumericProperty(c1, "Age")
	age2 := s.getNumericProperty(c2, "Age")
	if age1 > 0 && age2 > 0 {
		ageDiff := math.Abs(age1 - age2)
		ageScore := math.Max(0, 1.0-(ageDiff/domain.AgeNormalizationFactor))
		score += ageScore
		total += 1.0
	}

	// Income similarity
	income1 := s.getNumericProperty(c1, "Annual Income (THB)")
	income2 := s.getNumericProperty(c2, "Annual Income (THB)")
	if income1 > 0 && income2 > 0 {
		incomeRatio := math.Min(income1, income2) / math.Max(income1, income2)
		score += incomeRatio
		total += 1.0
	}

	// Investment experience similarity
	exp1 := s.getNumericProperty(c1, "Investment Experience (Years)")
	exp2 := s.getNumericProperty(c2, "Investment Experience (Years)")
	if exp1 > 0 && exp2 > 0 {
		expDiff := math.Abs(exp1 - exp2)
		expScore := math.Max(0, 1.0-(expDiff/domain.ExperienceNormalizationFactor))
		score += expScore
		total += 1.0
	}

	// Occupation match
	if c1.GetPropertyValue("Occupation") == c2.GetPropertyValue("Occupation") {
		matches += 1.0
	}
	total += 1.0

	// Risk profile match
	if c1.RiskProfile == c2.RiskProfile {
		score += 1.0
	}
	total += 1.0

	if total > 0 {
		return (score + matches) / total
	}
	return 0.0
}

// calculateSurveySimilarity compares survey responses using Jaccard/Dice similarity
func (s *collaborativeFilterService) calculateSurveySimilarity(c1, c2 *domain.Customer) float64 {
	score := 0.0
	total := 0.0

	// Investment types match using Jaccard
	types1 := s.getStringProperty(c1, "Interesting Investment Types")
	types2 := s.getStringProperty(c2, "Interesting Investment Types")
	if types1 != "" && types2 != "" {
		score += utils.JaccardSimilarity(types1, types2)
		total += 1.0
	}

	// Industry sectors match using Jaccard
	sectors1 := s.getStringProperty(c1, "Interesting Industry Sectors")
	sectors2 := s.getStringProperty(c2, "Interesting Industry Sectors")
	if sectors1 != "" && sectors2 != "" {
		score += utils.JaccardSimilarity(sectors1, sectors2)
		total += 1.0
	}

	// Regional preferences using Dice coefficient (more forgiving)
	regions1 := s.getStringProperty(c1, "Regional Preferences")
	regions2 := s.getStringProperty(c2, "Regional Preferences")
	if regions1 != "" && regions2 != "" {
		score += utils.DiceCoefficientSimilarity(regions1, regions2)
		total += 1.0
	}

	// ESG interest match
	esg1 := s.getStringProperty(c1, "ESG Interest")
	esg2 := s.getStringProperty(c2, "ESG Interest")
	if esg1 != "" && esg2 != "" {
		if esg1 == esg2 {
			score += 1.0
		}
		total += 1.0
	}

	// Crypto interest match
	crypto1 := s.getStringProperty(c1, "Crypto Interest")
	crypto2 := s.getStringProperty(c2, "Crypto Interest")
	if crypto1 != "" && crypto2 != "" {
		if crypto1 == crypto2 {
			score += 1.0
		}
		total += 1.0
	}

	if total > 0 {
		return score / total
	}
	return 0.0
}

// calculateActivitySimilarity compares activity patterns with fallback for cold start
func (s *collaborativeFilterService) calculateActivitySimilarity(c1, c2 *domain.Customer) float64 {
	score := 0.0
	total := 0.0

	// Search activity similarity
	searches1 := s.getNumericProperty(c1, "Total Searches")
	searches2 := s.getNumericProperty(c2, "Total Searches")
	if searches1 > 0 || searches2 > 0 {
		ratio := s.calculateRatioSimilarity(searches1, searches2)
		score += ratio
		total += 1.0
	}

	// Views similarity
	views1 := s.getNumericProperty(c1, "Total Views")
	views2 := s.getNumericProperty(c2, "Total Views")
	if views1 > 0 || views2 > 0 {
		ratio := s.calculateRatioSimilarity(views1, views2)
		score += ratio
		total += 1.0
	}

	// Factsheet reads similarity
	factsheet1 := s.getNumericProperty(c1, "Factsheet Reads")
	factsheet2 := s.getNumericProperty(c2, "Factsheet Reads")
	if factsheet1 > 0 || factsheet2 > 0 {
		ratio := s.calculateRatioSimilarity(factsheet1, factsheet2)
		score += ratio
		total += 1.0
	}

	// Share activity similarity
	shares1 := s.getNumericProperty(c1, "Share Count")
	shares2 := s.getNumericProperty(c2, "Share Count")
	if shares1 > 0 || shares2 > 0 {
		ratio := s.calculateRatioSimilarity(shares1, shares2)
		score += ratio
		total += 1.0
	}

	// If no activity at all, return neutral score instead of 0
	if total == 0 {
		return domain.ColdStartNeutralSimilarity
	}

	return score / total
}

// calculateRatioSimilarity calculates similarity with fallback for zero values (cold start)
func (s *collaborativeFilterService) calculateRatioSimilarity(v1, v2 float64) float64 {
	// Both zero = neutral similarity (not perfect match)
	if v1 == 0 && v2 == 0 {
		return domain.ColdStartNeutralSimilarity
	}
	// One zero = low but not zero similarity
	if v1 == 0 || v2 == 0 {
		return domain.ColdStartLowSimilarity
	}
	// Both have values = ratio-based similarity
	return math.Min(v1, v2) / math.Max(v1, v2)
}

// calculateTransactionSimilarity compares transaction behavior
func (s *collaborativeFilterService) calculateTransactionSimilarity(c1, c2 *domain.Customer) float64 {
	score := 0.0
	total := 0.0

	// Buy count similarity
	buy1 := s.getNumericProperty(c1, "Buy Count")
	buy2 := s.getNumericProperty(c2, "Buy Count")
	if buy1 > 0 && buy2 > 0 {
		ratio := math.Min(buy1, buy2) / math.Max(buy1, buy2)
		score += ratio
		total += 1.0
	}

	// Sell count similarity
	sell1 := s.getNumericProperty(c1, "Sell Count")
	sell2 := s.getNumericProperty(c2, "Sell Count")
	if sell1 > 0 && sell2 > 0 {
		ratio := math.Min(sell1, sell2) / math.Max(sell1, sell2)
		score += ratio
		total += 1.0
	}

	// Switch count similarity
	switch1 := s.getNumericProperty(c1, "Switch Count")
	switch2 := s.getNumericProperty(c2, "Switch Count")
	if switch1 > 0 && switch2 > 0 {
		ratio := math.Min(switch1, switch2) / math.Max(switch1, switch2)
		score += ratio
		total += 1.0
	}

	// Fund detail access similarity
	fund1 := s.getNumericProperty(c1, "Fund Detail Access")
	fund2 := s.getNumericProperty(c2, "Fund Detail Access")
	if fund1 > 0 && fund2 > 0 {
		ratio := math.Min(fund1, fund2) / math.Max(fund1, fund2)
		score += ratio
		total += 1.0
	}

	if total > 0 {
		return score / total
	}
	return 0.0
}

// FindSimilarCustomers finds customers similar to the given customer
func (s *collaborativeFilterService) FindSimilarCustomers(ctx context.Context, customerID string, limit int) ([]*SimilarCustomer, error) {
	start := time.Now()
	s.logger.Info("Finding similar customers",
		zap.String("customer_id", customerID),
		zap.Int("limit", limit),
	)

	//Input validation
	if customerID == "" {
		s.logger.Error("Customer ID is required")
		return nil, errors.New("customer ID is required")
	}

	if limit < 0 {
		s.logger.Error("Limit must be non-negative", zap.Int("limit", limit))
		return nil, errors.New("limit must be non-negative")
	}

	// Set reasonable default
	if limit == 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100 // Prevent abuse
	}

	//Error wrapping with context
	targetCustomer, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		s.logger.Error("Failed to find target customer",
			zap.String("customer_id", customerID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to find target customer %s: %w", customerID, err)
	}

	// Get all other customers
	allCustomers, err := s.customerRepo.FindAll(ctx, 0, 0)
	if err != nil {
		s.logger.Error("Failed to fetch all customers",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to fetch all customers for similarity calculation: %w", err)
	}

	s.logger.Debug("Calculating similarity scores",
		zap.Int("total_customers", len(allCustomers)),
	)

	// Calculate similarity scores
	var similarities []*SimilarCustomer
	for _, customer := range allCustomers {
		if customer.ID == customerID {
			continue
		}

		score := s.CalculateSimilarity(targetCustomer, customer)

		// Only include customers above similarity threshold
		if score.Score >= s.cfg.CollaborativeFilter.SimilarityThreshold {
			similarities = append(similarities, &SimilarCustomer{
				Customer: customer,
				Score:    score,
			})
		}
	}

	// Sort by similarity score descending
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Score.Score > similarities[j].Score.Score
	})

	// Apply limit
	if limit > 0 && len(similarities) > limit {
		similarities = similarities[:limit]
	}

	duration := time.Since(start)
	s.logger.Info("Found similar customers",
		zap.String("customer_id", customerID),
		zap.Int("similar_count", len(similarities)),
		zap.Duration("duration", duration),
	)

	return similarities, nil
}

// GenerateRecommendations generates recommendations based on similar customers' behavior
func (s *collaborativeFilterService) GenerateRecommendations(ctx context.Context, customerID string, recType domain.RecommendationType) ([]*domain.Recommendation, error) {
	start := time.Now()
	s.logger.Info("Generating recommendations",
		zap.String("customer_id", customerID),
		zap.String("type", string(recType)),
	)

	//Input validation
	if customerID == "" {
		s.logger.Error("Customer ID is required")
		return nil, errors.New("customer ID is required")
	}

	// Validate recommendation type
	validTypes := map[domain.RecommendationType]bool{
		domain.RecommendationTypeFund:    true,
		domain.RecommendationTypeSector:  true,
		domain.RecommendationTypeProduct: true,
	}

	if !validTypes[recType] {
		s.logger.Error("Invalid recommendation type",
			zap.String("type", string(recType)),
		)
		return nil, fmt.Errorf("invalid recommendation type: %s", recType)
	}

	//Get customer and select strategy
	targetCustomer, err := s.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		s.logger.Error("Failed to find customer",
			zap.String("customer_id", customerID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to find customer %s: %w", customerID, err)
	}

	strategy := s.strategySelector.SelectStrategy(targetCustomer)
	s.logger.Info("Strategy selected",
		zap.String("customer_id", customerID),
		zap.String("strategy", string(strategy)),
	)

	// ✅ NEW: Route to appropriate strategy
	var recommendations []*domain.Recommendation
	switch strategy {
	case domain.StrategyCollaborative:
		recommendations, err = s.generateCollaborativeRecommendations(ctx, customerID, targetCustomer, recType)

	case domain.StrategyContentBased:
		recommendations, err = s.generateContentBasedRecommendations(ctx, targetCustomer, recType)

	case domain.StrategyPopularity:
		recommendations, err = s.generatePopularityRecommendations(ctx, targetCustomer, recType)

	case domain.StrategyHybrid:
		recommendations, err = s.generateHybridRecommendations(ctx, customerID, targetCustomer, recType)

	default:
		// Fallback to collaborative
		recommendations, err = s.generateCollaborativeRecommendations(ctx, customerID, targetCustomer, recType)
	}

	duration := time.Since(start)
	if err != nil {
		s.logger.Error("Failed to generate recommendations",
			zap.String("customer_id", customerID),
			zap.String("strategy", string(strategy)),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		return nil, err
	}

	s.logger.Info("Recommendations generated successfully",
		zap.String("customer_id", customerID),
		zap.String("strategy", string(strategy)),
		zap.Int("count", len(recommendations)),
		zap.Duration("duration", duration),
	)

	return recommendations, nil
}

// generateCollaborativeRecommendations uses collaborative filtering with similar customers
func (s *collaborativeFilterService) generateCollaborativeRecommendations(
	ctx context.Context,
	customerID string,
	customer *domain.Customer,
	recType domain.RecommendationType,
) ([]*domain.Recommendation, error) {
	s.logger.Debug("Using collaborative filtering strategy",
		zap.String("customer_id", customerID),
		zap.String("type", string(recType)),
	)

	// Find similar customers
	similarCustomers, err := s.FindSimilarCustomers(ctx, customerID, 10)
	if err != nil {
		s.logger.Error("Failed to find similar customers for collaborative filtering",
			zap.String("customer_id", customerID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to find similar customers: %w", err)
	}

	// ✅ If no similar customers, fall back to content-based
	if len(similarCustomers) == 0 {
		s.logger.Warn("No similar customers found, falling back to content-based filtering",
			zap.String("customer_id", customerID),
		)
		return s.generateContentBasedRecommendations(ctx, customer, recType)
	}

	// Generate recommendations based on type
	var recommendations []*domain.Recommendation

	switch recType {
	case domain.RecommendationTypeFund:
		recommendations = s.generateFundRecommendations(customerID, similarCustomers)
	case domain.RecommendationTypeSector:
		recommendations = s.generateSectorRecommendations(customerID, similarCustomers)
	case domain.RecommendationTypeProduct:
		recommendations = s.generateProductRecommendations(customerID, similarCustomers)
	}

	// Apply diversity
	beforeDiversity := len(recommendations)
	recommendations = s.diversifier.DiversifyRecommendations(recommendations)
	s.logger.Debug("Applied diversity filter",
		zap.Int("before", beforeDiversity),
		zap.Int("after", len(recommendations)),
	)

	// Apply limit
	maxRecs := s.cfg.CollaborativeFilter.MaxRecommendations
	if len(recommendations) > maxRecs {
		recommendations = recommendations[:maxRecs]
	}

	// Save to database
	for _, rec := range recommendations {
		if err := s.recommendationRepo.Create(ctx, rec); err != nil {
			s.logger.Error("Failed to save recommendation",
				zap.String("customer_id", customerID),
				zap.String("item_id", rec.ItemID),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to save recommendation: %w", err)
		}
	}

	s.logger.Info("Collaborative recommendations generated",
		zap.String("customer_id", customerID),
		zap.Int("count", len(recommendations)),
		zap.Int("similar_customers", len(similarCustomers)),
	)

	return recommendations, nil
}

func (s *collaborativeFilterService) generateFundRecommendations(customerID string, similarCustomers []*SimilarCustomer) []*domain.Recommendation {
	// Aggregate fund preferences from similar customers
	fundScores := make(map[string]float64)
	fundCounts := make(map[string]int)
	var similarCustomerIDs []string

	for _, sim := range similarCustomers {
		similarCustomerIDs = append(similarCustomerIDs, sim.Customer.ID)

		// Get interesting investment types
		types := s.getStringProperty(sim.Customer, "Interesting Investment Types")
		if types != "" {
			typeList := strings.Split(types, ",")
			for _, t := range typeList {
				fund := strings.TrimSpace(t)
				if fund != "" {
					// Parse and weight by similarity score
					fundScores[fund] += sim.Score.Score
					fundCounts[fund]++
				}
			}
		}
	}

	// Create recommendations
	var recommendations []*domain.Recommendation
	for fund, totalScore := range fundScores {
		avgScore := totalScore / float64(fundCounts[fund])
		avgSimilarity := s.calculateAverageSimilarity(similarCustomers)

		metadata := domain.RecommendationMeta{
			Algorithm:     "User-based Collaborative Filtering",
			Confidence:    avgScore,
			BasedOnUsers:  len(similarCustomers),
			SimilarityAvg: avgSimilarity,
		}
		metadataJSON, _ := json.Marshal(metadata)

		rec := &domain.Recommendation{
			CustomerID:       customerID,
			Type:             domain.RecommendationTypeFund,
			ItemID:           fund,
			ItemName:         fund,
			Score:            avgScore,
			Reason:           "Based on similar customers with " + strconv.FormatFloat(avgSimilarity, 'f', 2, 64) + " similarity score",
			Status:           domain.RecommendationStatusActive,
			Metadata:         metadataJSON,
			SimilarCustomers: similarCustomerIDs,
		}
		recommendations = append(recommendations, rec)
	}

	return recommendations
}

func (s *collaborativeFilterService) generateSectorRecommendations(customerID string, similarCustomers []*SimilarCustomer) []*domain.Recommendation {
	// Aggregate sector preferences from similar customers
	sectorScores := make(map[string]float64)
	sectorCounts := make(map[string]int)
	var similarCustomerIDs []string

	for _, sim := range similarCustomers {
		similarCustomerIDs = append(similarCustomerIDs, sim.Customer.ID)

		// Get interesting industry sectors
		sectors := s.getStringProperty(sim.Customer, "Interesting Industry Sectors")
		if sectors != "" {
			sectorList := strings.Split(sectors, ",")
			for _, s := range sectorList {
				sector := strings.TrimSpace(s)
				if sector != "" {
					sectorScores[sector] += sim.Score.Score
					sectorCounts[sector]++
				}
			}
		}
	}

	// Create recommendations
	var recommendations []*domain.Recommendation
	for sector, totalScore := range sectorScores {
		avgScore := totalScore / float64(sectorCounts[sector])
		avgSimilarity := s.calculateAverageSimilarity(similarCustomers)

		metadata := domain.RecommendationMeta{
			Algorithm:     "User-based Collaborative Filtering",
			Confidence:    avgScore,
			BasedOnUsers:  len(similarCustomers),
			SimilarityAvg: avgSimilarity,
		}
		metadataJSON, _ := json.Marshal(metadata)

		rec := &domain.Recommendation{
			CustomerID:       customerID,
			Type:             domain.RecommendationTypeSector,
			ItemID:           sector,
			ItemName:         sector,
			Score:            avgScore,
			Reason:           "Recommended by " + strconv.Itoa(sectorCounts[sector]) + " similar customers",
			Status:           domain.RecommendationStatusActive,
			Metadata:         metadataJSON,
			SimilarCustomers: similarCustomerIDs,
		}
		recommendations = append(recommendations, rec)
	}

	return recommendations
}

func (s *collaborativeFilterService) generateProductRecommendations(customerID string, similarCustomers []*SimilarCustomer) []*domain.Recommendation {
	// Similar logic to funds but for products
	return s.generateFundRecommendations(customerID, similarCustomers)
}

// Helper functions
func (s *collaborativeFilterService) getNumericProperty(customer *domain.Customer, key string) float64 {
	value := customer.GetPropertyValue(key)
	if value == "" {
		return 0
	}

	numValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return numValue
}

func (s *collaborativeFilterService) getStringProperty(customer *domain.Customer, key string) string {
	return customer.GetPropertyValue(key)
}

// Removed - using utils.JaccardSimilarity instead

func (s *collaborativeFilterService) calculateAverageSimilarity(customers []*SimilarCustomer) float64 {
	if len(customers) == 0 {
		return 0
	}

	total := 0.0
	for _, customer := range customers {
		total += customer.Score.Score
	}
	return total / float64(len(customers))
}

// generateContentBasedRecommendations uses customer's own profile for cold start
func (s *collaborativeFilterService) generateContentBasedRecommendations(
	ctx context.Context,
	customer *domain.Customer,
	recType domain.RecommendationType,
) ([]*domain.Recommendation, error) {
	s.logger.Debug("Using content-based filtering strategy",
		zap.String("customer_id", customer.ID),
		zap.String("type", string(recType)),
	)

	var recommendations []*domain.Recommendation

	// Use customer's stated preferences from survey
	switch recType {
	case domain.RecommendationTypeFund:
		types := customer.GetPropertyValue("Interesting Investment Types")
		if types != "" {
			typeList := strings.Split(types, ",")
			for i, t := range typeList {
				if i >= s.cfg.CollaborativeFilter.MaxRecommendations {
					break
				}

				fund := strings.TrimSpace(t)
				if fund != "" {
					metadata := domain.RecommendationMeta{
						Algorithm:    "Content-Based Filtering",
						Confidence:   domain.ContentBasedConfidence,
						BasedOnUsers: 0,
					}
					metadataJSON, _ := json.Marshal(metadata)

					rec := &domain.Recommendation{
						CustomerID: customer.ID,
						Type:       recType,
						ItemID:     fund,
						ItemName:   fund,
						Score:      domain.ContentBasedConfidence,
						Reason:     fmt.Sprintf("Matches your %s risk profile and stated interest", customer.RiskProfile),
						Status:     domain.RecommendationStatusActive,
						Metadata:   metadataJSON,
					}
					recommendations = append(recommendations, rec)
				}
			}
		}

	case domain.RecommendationTypeSector:
		sectors := customer.GetPropertyValue("Interesting Industry Sectors")
		if sectors != "" {
			sectorList := strings.Split(sectors, ",")
			for i, sec := range sectorList {
				if i >= s.cfg.CollaborativeFilter.MaxRecommendations {
					break
				}

				sector := strings.TrimSpace(sec)
				if sector != "" {
					metadata := domain.RecommendationMeta{
						Algorithm:    "Content-Based Filtering",
						Confidence:   domain.ContentBasedConfidence,
						BasedOnUsers: 0,
					}
					metadataJSON, _ := json.Marshal(metadata)

					rec := &domain.Recommendation{
						CustomerID: customer.ID,
						Type:       recType,
						ItemID:     sector,
						ItemName:   sector,
						Score:      domain.ContentBasedConfidence,
						Reason:     "Based on your stated sector preferences",
						Status:     domain.RecommendationStatusActive,
						Metadata:   metadataJSON,
					}
					recommendations = append(recommendations, rec)
				}
			}
		}
	}

	// If still no recommendations, fall back to popularity
	if len(recommendations) == 0 {
		s.logger.Warn("No content-based recommendations found, falling back to popularity",
			zap.String("customer_id", customer.ID),
		)
		return s.generatePopularityRecommendations(ctx, customer, recType)
	}

	// Save to database
	for _, rec := range recommendations {
		if err := s.recommendationRepo.Create(ctx, rec); err != nil {
			s.logger.Error("Failed to save content-based recommendation",
				zap.String("customer_id", customer.ID),
				zap.String("item_id", rec.ItemID),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to save recommendation: %w", err)
		}
	}

	s.logger.Info("Content-based recommendations generated",
		zap.String("customer_id", customer.ID),
		zap.Int("count", len(recommendations)),
	)

	return recommendations, nil
}

// generatePopularityRecommendations returns popular items for extreme cold start
func (s *collaborativeFilterService) generatePopularityRecommendations(
	ctx context.Context,
	customer *domain.Customer,
	recType domain.RecommendationType,
) ([]*domain.Recommendation, error) {
	s.logger.Debug("Using popularity-based filtering strategy",
		zap.String("customer_id", customer.ID),
		zap.String("type", string(recType)),
	)

	var recommendations []*domain.Recommendation
	var popularItems []string

	// Predefined popular items (for POC)
	switch recType {
	case domain.RecommendationTypeFund:
		popularItems = []string{"Equity Funds", "Bond Funds", "Mixed Funds", "Index Funds", "ESG Funds"}
	case domain.RecommendationTypeSector:
		popularItems = []string{"Technology", "Healthcare", "Finance", "Energy", "Consumer Goods"}
	case domain.RecommendationTypeProduct:
		popularItems = []string{"Mutual Funds", "ETFs", "Stocks", "Bonds", "Derivatives"}
	}

	for i, item := range popularItems {
		if i >= s.cfg.CollaborativeFilter.MaxRecommendations {
			break
		}

		score := domain.PopularityBaseScore - (float64(i) * domain.PopularityScoreDecrement)
		if score < domain.PopularityMinScore {
			score = domain.PopularityMinScore
		}

		metadata := domain.RecommendationMeta{
			Algorithm:    "Popularity-Based",
			Confidence:   score,
			BasedOnUsers: 0,
		}
		metadataJSON, _ := json.Marshal(metadata)

		rec := &domain.Recommendation{
			CustomerID: customer.ID,
			Type:       recType,
			ItemID:     item,
			ItemName:   item,
			Score:      score,
			Reason:     "Popular choice among all customers",
			Status:     domain.RecommendationStatusActive,
			Metadata:   metadataJSON,
		}
		recommendations = append(recommendations, rec)
	}

	// Save to database
	for _, rec := range recommendations {
		if err := s.recommendationRepo.Create(ctx, rec); err != nil {
			s.logger.Error("Failed to save popularity-based recommendation",
				zap.String("customer_id", customer.ID),
				zap.String("item_id", rec.ItemID),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to save recommendation: %w", err)
		}
	}

	s.logger.Info("Popularity-based recommendations generated",
		zap.String("customer_id", customer.ID),
		zap.Int("count", len(recommendations)),
	)

	return recommendations, nil
}

// generateHybridRecommendations combines collaborative and content-based
func (s *collaborativeFilterService) generateHybridRecommendations(
	ctx context.Context,
	customerID string,
	customer *domain.Customer,
	recType domain.RecommendationType,
) ([]*domain.Recommendation, error) {
	s.logger.Debug("Using hybrid filtering strategy",
		zap.String("customer_id", customerID),
		zap.String("type", string(recType)),
	)

	// Get both collaborative and content-based
	collabRecs, _ := s.generateCollaborativeRecommendations(ctx, customerID, customer, recType)
	contentRecs, _ := s.generateContentBasedRecommendations(ctx, customer, recType)

	s.logger.Debug("Merging collaborative and content-based recommendations",
		zap.Int("collaborative_count", len(collabRecs)),
		zap.Int("content_based_count", len(contentRecs)),
	)

	// Merge recommendations
	mergedMap := make(map[string]*domain.Recommendation)

	// Add collaborative recommendations
	for _, rec := range collabRecs {
		mergedMap[rec.ItemID] = rec
	}

	// Merge content-based recommendations
	for _, rec := range contentRecs {
		if existing, exists := mergedMap[rec.ItemID]; exists {
			// Average scores for items in both
			existing.Score = (existing.Score + rec.Score) / 2.0
			existing.Reason = "Recommended by both collaborative filtering and your preferences"

			var metadata domain.RecommendationMeta
			json.Unmarshal(existing.Metadata, &metadata)
			metadata.Algorithm = "Hybrid"
			metadataJSON, _ := json.Marshal(metadata)
			existing.Metadata = metadataJSON
		} else {
			// Add with reduced weight
			rec.Score *= domain.HybridContentReductionFactor
			mergedMap[rec.ItemID] = rec
		}
	}

	// Convert to slice and sort
	var recommendations []*domain.Recommendation
	for _, rec := range mergedMap {
		recommendations = append(recommendations, rec)
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Apply limit
	maxRecs := s.cfg.CollaborativeFilter.MaxRecommendations
	if len(recommendations) > maxRecs {
		recommendations = recommendations[:maxRecs]
	}

	s.logger.Info("Hybrid recommendations generated",
		zap.String("customer_id", customerID),
		zap.Int("count", len(recommendations)),
		zap.Int("merged_items", len(mergedMap)),
	)

	return recommendations, nil
}
