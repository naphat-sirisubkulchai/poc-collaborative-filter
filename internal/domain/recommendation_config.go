package domain

// RecommendationConfig configures the recommendation algorithm weights
type RecommendationConfig struct {
	SimilarityWeight float64 // 0.6 - Focus on similarity
	DiversityWeight  float64 // 0.25 - Encourage exploration
	NoveltyWeight    float64 // 0.10 - Show new/trending items
	PopularityWeight float64 // 0.05 - Consider popularity
}

// DefaultRecommendationConfig returns the default configuration
func DefaultRecommendationConfig() RecommendationConfig {
	return RecommendationConfig{
		SimilarityWeight: 0.6,
		DiversityWeight:  0.25,
		NoveltyWeight:    0.10,
		PopularityWeight: 0.05,
	}
}

// AggressiveDiversityConfig favors diversity over similarity
func AggressiveDiversityConfig() RecommendationConfig {
	return RecommendationConfig{
		SimilarityWeight: 0.4,
		DiversityWeight:  0.40,
		NoveltyWeight:    0.15,
		PopularityWeight: 0.05,
	}
}

// ConservativeConfig favors similarity (less exploration)
func ConservativeConfig() RecommendationConfig {
	return RecommendationConfig{
		SimilarityWeight: 0.75,
		DiversityWeight:  0.15,
		NoveltyWeight:    0.05,
		PopularityWeight: 0.05,
	}
}
