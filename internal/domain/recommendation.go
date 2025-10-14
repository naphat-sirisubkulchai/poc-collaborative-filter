package domain

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RecommendationType string

const (
	RecommendationTypeProduct  RecommendationType = "PRODUCT"
	RecommendationTypeCustomer RecommendationType = "CUSTOMER"
	RecommendationTypeFund     RecommendationType = "FUND"
	RecommendationTypeSector   RecommendationType = "SECTOR"
)

type RecommendationStatus string

const (
	RecommendationStatusActive    RecommendationStatus = "ACTIVE"
	RecommendationStatusDismissed RecommendationStatus = "DISMISSED"
	RecommendationStatusAccepted  RecommendationStatus = "ACCEPTED"
	RecommendationStatusExpired   RecommendationStatus = "EXPIRED"
)

// Recommendation represents a recommendation made to a customer
type Recommendation struct {
	ID               string               `gorm:"primarykey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CustomerID       string               `gorm:"type:uuid;not null;index" json:"customerId"`
	Type             RecommendationType   `gorm:"not null" json:"type"`
	ItemID           string               `gorm:"not null" json:"itemId"`
	ItemName         string               `gorm:"not null" json:"itemName"`
	Score            float64              `gorm:"not null" json:"score"`
	Reason           string               `json:"reason"`
	Status           RecommendationStatus `gorm:"default:'ACTIVE'" json:"status"`
	Metadata         datatypes.JSON       `gorm:"type:jsonb" json:"metadata"`
	SimilarCustomers pq.StringArray       `gorm:"type:text[]" json:"similarCustomers"`
	CreatedAt        time.Time            `json:"createdAt"`
	UpdatedAt        time.Time            `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt       `gorm:"index" json:"-"`

	// Relationships
	Customer Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

// TableName specifies the table name for Recommendation model
func (Recommendation) TableName() string {
	return "recommendations"
}

// SimilarityScore represents the similarity between two customers
type SimilarityScore struct {
	Customer1ID string  `json:"customer1Id"`
	Customer2ID string  `json:"customer2Id"`
	Score       float64 `json:"score"`
	Weights     Weights `json:"weights"`
}

type Weights struct {
	Personal       float64 `json:"personal"`
	Survey         float64 `json:"survey"`
	ActivityAccess float64 `json:"activityAccess"`
	Transaction    float64 `json:"transaction"`
}

// RecommendationMeta holds metadata about how the recommendation was generated
type RecommendationMeta struct {
	Algorithm     string            `json:"algorithm"`
	Parameters    map[string]string `json:"parameters,omitempty"`
	Confidence    float64           `json:"confidence"`
	BasedOnUsers  int               `json:"basedOnUsers"`
	SimilarityAvg float64           `json:"similarityAvg"`
}
