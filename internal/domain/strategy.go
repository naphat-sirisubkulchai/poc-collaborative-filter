package domain

type RecommendationStrategy string

const (
	StrategyCollaborative RecommendationStrategy = "COLLABORATIVE"
	StrategyContentBased  RecommendationStrategy = "CONTENT_BASED"
	StrategyHybrid        RecommendationStrategy = "HYBRID"
	StrategyPopularity    RecommendationStrategy = "POPULARITY"
)

// StrategySelector determines which recommendation strategy to use based on customer data
type StrategySelector struct {
	ActivityThreshold  float64
	ProfileThreshold   int
	TransactionMinimum int
}

func NewStrategySelector() *StrategySelector {
	return &StrategySelector{
		ActivityThreshold:  MinActivityScoreForHybrid,
		ProfileThreshold:   MinProfilePropertiesForStrategy,
		TransactionMinimum: MinTransactionsForCollaborative,
	}
}

// SelectStrategy determines the best strategy for a customer
func (ss *StrategySelector) SelectStrategy(customer *Customer) RecommendationStrategy {
	traits, _ := customer.GetTraits()
	if traits == nil {
		return StrategyPopularity
	}

	activityScore := ss.calculateActivityScore(customer)
	profileCompleteness := len(traits.Properties)

	// New customer with minimal data → Content-based or Popularity
	if activityScore < ss.ActivityThreshold && profileCompleteness < ss.ProfileThreshold {
		if profileCompleteness > MinProfilePropertiesForContent {
			return StrategyContentBased // Use profile data
		}
		return StrategyPopularity // Fall back to popular items
	}

	// Rich activity data → Pure Collaborative
	if activityScore > MinActivityScoreForCollaborative {
		return StrategyCollaborative
	}

	// Middle ground → Hybrid (combine both)
	return StrategyHybrid
}

func (ss *StrategySelector) calculateActivityScore(customer *Customer) float64 {
	score := 0.0
	traits, err := customer.GetTraits()
	if err != nil || traits == nil {
		return 0.0
	}

	// Helper to get numeric property value
	getNum := func(key string) float64 {
		for _, prop := range traits.Properties {
			if prop.Key == key && prop.Type == "number" {
				var val float64
				if n, err := parseFloat(prop.Value); err == nil {
					val = n
				}
				return val
			}
		}
		return 0.0
	}

	score += getNum("Total Searches") * SearchWeight
	score += getNum("Total Views") * ViewWeight
	score += getNum("Factsheet Reads") * FactsheetWeight
	score += getNum("Buy Count") * BuyWeight
	score += getNum("Sell Count") * SellWeight
	score += getNum("Switch Count") * SwitchWeight

	return score
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := parseFloatHelper(s, &f)
	return f, err
}

func parseFloatHelper(s string, f *float64) (int, error) {
	// Simple float parsing - in production use strconv.ParseFloat
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			*f = *f*10 + float64(s[i]-'0')
			n++
		} else if s[i] == '.' {
			break
		}
	}
	return n, nil
}
