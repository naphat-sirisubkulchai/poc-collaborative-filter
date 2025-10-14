package service

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"poc-collaborative-filter/internal/domain"
)

// Diversifier applies diversification to recommendations
type Diversifier struct {
	config domain.RecommendationConfig
}

// NewDiversifier creates a new diversifier with given config
func NewDiversifier(config domain.RecommendationConfig) *Diversifier {
	return &Diversifier{config: config}
}

// DiversifyRecommendations applies diversity, novelty, and popularity factors
func (d *Diversifier) DiversifyRecommendations(
	recommendations []*domain.Recommendation,
) []*domain.Recommendation {

	if len(recommendations) == 0 {
		return recommendations
	}

	var diversified []*domain.Recommendation
	seenCategories := make(map[string]int)
	seenTypes := make(map[string]int)

	for _, rec := range recommendations {
		// Extract category/type from item name or metadata
		category := d.extractCategory(rec)
		itemType := string(rec.Type)

		// Calculate diversity penalty
		categoryPenalty := float64(seenCategories[category]) * 0.15
		typePenalty := float64(seenTypes[itemType]) * 0.10
		totalPenalty := categoryPenalty + typePenalty

		// Calculate novelty boost (for newer recommendations)
		noveltyBoost := d.calculateNoveltyBoost(rec)

		// Calculate popularity factor (if available in metadata)
		popularityScore := d.calculatePopularityScore(rec)

		// Save original score
		originalScore := rec.Score

		// Adjust score with all factors
		adjustedScore := (originalScore * d.config.SimilarityWeight) -
			(totalPenalty * d.config.DiversityWeight) +
			(noveltyBoost * d.config.NoveltyWeight) +
			(popularityScore * d.config.PopularityWeight)

		// Ensure score doesn't go negative
		if adjustedScore < 0 {
			adjustedScore = originalScore * 0.1 // Keep it low but not zero
		}

		rec.Score = adjustedScore
		diversified = append(diversified, rec)

		// Track seen categories/types
		seenCategories[category]++
		seenTypes[itemType]++
	}

	// Re-sort by adjusted scores
	sort.Slice(diversified, func(i, j int) bool {
		return diversified[i].Score > diversified[j].Score
	})

	return diversified
}

// extractCategory extracts category from recommendation
func (d *Diversifier) extractCategory(rec *domain.Recommendation) string {
	// Try to extract from metadata first
	if rec.Metadata != nil {
		var metadata map[string]interface{}
		if err := json.Unmarshal(rec.Metadata, &metadata); err == nil {
			if category, ok := metadata["category"].(string); ok {
				return category
			}
		}
	}

	// Fallback: extract from item name
	// e.g., "Technology Sector" -> "Technology"
	parts := strings.Split(rec.ItemName, " ")
	if len(parts) > 0 {
		return parts[0]
	}

	return "UNKNOWN"
}

// calculateNoveltyBoost gives higher scores to newer recommendations
func (d *Diversifier) calculateNoveltyBoost(rec *domain.Recommendation) float64 {
	// Items created in last 30 days get novelty boost
	daysSinceCreated := time.Since(rec.CreatedAt).Hours() / 24

	if daysSinceCreated <= 7 {
		return 1.0 // Very new (1 week)
	} else if daysSinceCreated <= 30 {
		// Linear decay from 1.0 to 0.3 over 30 days
		return 1.0 - ((daysSinceCreated - 7) / 30.0 * 0.7)
	}

	return 0.0 // Older than 30 days - no boost
}

// calculatePopularityScore extracts popularity from metadata
func (d *Diversifier) calculatePopularityScore(rec *domain.Recommendation) float64 {
	if rec.Metadata == nil {
		return 0.0
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(rec.Metadata, &metadata); err != nil {
		return 0.0
	}

	// Check various popularity indicators
	if basedOnUsers, ok := metadata["basedOnUsers"].(float64); ok {
		// More similar customers = more popular
		// Normalize: 1 user = 0.1, 10 users = 1.0
		return (basedOnUsers / 10.0) * 0.5
	}

	return 0.0
}

// AddCategoryToMetadata adds category information to recommendation metadata
func (d *Diversifier) AddCategoryToMetadata(rec *domain.Recommendation, category string) error {
	var metadata map[string]interface{}

	if rec.Metadata != nil {
		if err := json.Unmarshal(rec.Metadata, &metadata); err != nil {
			metadata = make(map[string]interface{})
		}
	} else {
		metadata = make(map[string]interface{})
	}

	metadata["category"] = category

	updatedMetadata, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	rec.Metadata = updatedMetadata
	return nil
}
