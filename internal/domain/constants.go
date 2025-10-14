package domain

// Strategy thresholds for recommendation algorithm selection
const (
	// MinActivityScoreForCollaborative is the minimum activity score required to use collaborative filtering
	MinActivityScoreForCollaborative = 20.0

	// MinActivityScoreForHybrid is the minimum activity score required to use hybrid filtering
	MinActivityScoreForHybrid = 5.0

	// MinProfilePropertiesForContent is the minimum number of profile properties for content-based filtering
	MinProfilePropertiesForContent = 5

	// MinProfilePropertiesForStrategy is the minimum profile completeness threshold
	MinProfilePropertiesForStrategy = 10

	// MinTransactionsForCollaborative is the minimum number of transactions for collaborative filtering
	MinTransactionsForCollaborative = 3
)

// Activity scoring weights for calculating customer engagement
const (
	// SearchWeight is the weight for search activity
	SearchWeight = 0.1

	// ViewWeight is the weight for view activity
	ViewWeight = 0.1

	// FactsheetWeight is the weight for factsheet reading activity
	FactsheetWeight = 0.2

	// BuyWeight is the weight for buy transactions (highest weight)
	BuyWeight = 1.0

	// SellWeight is the weight for sell transactions
	SellWeight = 0.8

	// SwitchWeight is the weight for switch transactions
	SwitchWeight = 0.6
)

// Similarity calculation constants
const (
	// AgeNormalizationFactor normalizes age differences (in years)
	AgeNormalizationFactor = 50.0

	// ExperienceNormalizationFactor normalizes investment experience differences (in years)
	ExperienceNormalizationFactor = 20.0

	// ColdStartNeutralSimilarity is the default similarity for cold start customers with no activity
	ColdStartNeutralSimilarity = 0.5

	// ColdStartLowSimilarity is the similarity when one customer has activity and another doesn't
	ColdStartLowSimilarity = 0.2
)

// Recommendation scoring constants
const (
	// ContentBasedConfidence is the default confidence score for content-based recommendations
	ContentBasedConfidence = 0.8

	// PopularityBaseScore is the starting score for popularity-based recommendations
	PopularityBaseScore = 0.6

	// PopularityScoreDecrement is the score decrement per rank in popularity list
	PopularityScoreDecrement = 0.1

	// PopularityMinScore is the minimum score for popularity-based recommendations
	PopularityMinScore = 0.1

	// HybridContentReductionFactor is the reduction factor for content-based items in hybrid mode
	HybridContentReductionFactor = 0.7
)

// Default similarity weights
const (
	// PersonalWeight is the weight for personal information similarity
	PersonalWeight = 0.25

	// SurveyWeight is the weight for survey data similarity
	SurveyWeight = 0.35

	// ActivityAccessWeight is the weight for activity and access patterns
	ActivityAccessWeight = 0.20

	// TransactionWeight is the weight for transaction behavior
	TransactionWeight = 0.20
)
