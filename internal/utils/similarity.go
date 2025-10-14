package utils

import (
	"math"
	"strings"
)

// JaccardSimilarity calculates the Jaccard index between two comma-separated strings
// Jaccard = |A ∩ B| / |A ∪ B|
func JaccardSimilarity(s1, s2 string) float64 {
	if s1 == "" && s2 == "" {
		return 1.0
	}
	if s1 == "" || s2 == "" {
		return 0.0
	}

	set1 := stringToSet(s1)
	set2 := stringToSet(s2)

	intersection := 0
	union := make(map[string]bool)

	// Calculate intersection and union
	for k := range set1 {
		union[k] = true
		if set2[k] {
			intersection++
		}
	}

	for k := range set2 {
		union[k] = true
	}

	if len(union) == 0 {
		return 0.0
	}

	return float64(intersection) / float64(len(union))
}

// DiceCoefficientSimilarity calculates Dice coefficient (F1-like measure)
// Dice = 2 * |A ∩ B| / (|A| + |B|)
// More forgiving than Jaccard for partial matches
func DiceCoefficientSimilarity(s1, s2 string) float64 {
	if s1 == "" && s2 == "" {
		return 1.0
	}
	if s1 == "" || s2 == "" {
		return 0.0
	}

	set1 := stringToSet(s1)
	set2 := stringToSet(s2)

	intersection := 0
	for k := range set1 {
		if set2[k] {
			intersection++
		}
	}

	if len(set1)+len(set2) == 0 {
		return 0.0
	}

	return 2.0 * float64(intersection) / float64(len(set1)+len(set2))
}

// CosineSimilarity calculates cosine similarity between two sets with weights
// Useful for weighted similarities
func CosineSimilarity(set1, set2 map[string]float64) float64 {
	if len(set1) == 0 || len(set2) == 0 {
		return 0.0
	}

	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0

	for k, v1 := range set1 {
		if v2, exists := set2[k]; exists {
			dotProduct += v1 * v2
		}
		magnitude1 += v1 * v1
	}

	for _, v2 := range set2 {
		magnitude2 += v2 * v2
	}

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))
}

func stringToSet(s string) map[string]bool {
	set := make(map[string]bool)
	items := strings.Split(s, ",")
	for _, item := range items {
		trimmed := strings.TrimSpace(strings.ToLower(item))
		if trimmed != "" {
			set[trimmed] = true
		}
	}
	return set
}
