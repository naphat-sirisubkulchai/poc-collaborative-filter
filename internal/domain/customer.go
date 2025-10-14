package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CustomerSegment string

const (
	CustomerSegmentRetail   CustomerSegment = "RETAIL"
	CustomerSegmentAffluent CustomerSegment = "AFFLUENT"
	CustomerSegmentHNI      CustomerSegment = "HNI"
	CustomerSegmentUHNI     CustomerSegment = "UHNI"
)

type CustomerStatus string

const (
	CustomerStatusActive   CustomerStatus = "ACTIVE"
	CustomerStatusInactive CustomerStatus = "INACTIVE"
	CustomerStatusProspect CustomerStatus = "PROSPECT"
	CustomerStatusDormant  CustomerStatus = "DORMANT"
)

type CustomerRiskProfile string

const (
	CustomerRiskProfileConservative CustomerRiskProfile = "CONSERVATIVE"
	CustomerRiskProfileModerate     CustomerRiskProfile = "MODERATE"
	CustomerRiskProfileBalanced     CustomerRiskProfile = "BALANCED"
	CustomerRiskProfileGrowth       CustomerRiskProfile = "GROWTH"
	CustomerRiskProfileAggressive   CustomerRiskProfile = "AGGRESSIVE"
)

type CustomerTraits struct {
	Personality    []string        `json:"personality,omitempty"`
	Preferences    []string        `json:"preferences,omitempty"`
	Notes          string          `json:"notes,omitempty"`
	Tags           []string        `json:"tags,omitempty"`
	Properties     []TraitProperty `json:"properties,omitempty"`
	RelatedPersons []RelatedPerson `json:"relatedPersons,omitempty"`
}

type TraitProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"` // text, number, date, boolean
}

type RelatedPerson struct {
	Name         string  `json:"name"`
	Relationship string  `json:"relationship"`
	Phone        *string `json:"phone,omitempty"`
	Email        *string `json:"email,omitempty"`
	Notes        *string `json:"notes,omitempty"`
}

type Customer struct {
	ID                      string              `gorm:"primarykey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name                    string              `gorm:"not null" json:"name"`
	Title                   *string             `json:"title"`
	Email                   string              `gorm:"uniqueIndex;not null" json:"email"`
	Phone                   *string             `json:"phone"`
	Segment                 CustomerSegment     `gorm:"not null" json:"segment"`
	AUM                     decimal.Decimal     `gorm:"type:decimal(20,2);not null" json:"aum"`
	AccountType             pq.StringArray      `gorm:"type:text[]" json:"accountType"`
	AssignedRMName          *string             `json:"assignedRMName"`
	LastContact             *time.Time          `json:"lastContact"`
	Status                  CustomerStatus      `gorm:"not null" json:"status"`
	RiskProfile             CustomerRiskProfile `gorm:"not null" json:"riskProfile"`
	JoinDate                time.Time           `gorm:"not null" json:"joinDate"`
	TotalProducts           int                 `gorm:"default:0" json:"totalProducts"`
	Address                 *string             `json:"address"`
	RelationshipManagerName *string             `json:"relationshipManagerName"`
	PerformanceYTD          decimal.Decimal     `gorm:"type:decimal(20,2);default:0" json:"performanceYTD"`
	Holdings                pq.StringArray      `gorm:"type:text[]" json:"holdings"`
	ContactHistory          int                 `gorm:"default:0" json:"contactHistory"`
	Traits                  datatypes.JSON      `gorm:"type:jsonb" json:"traits"`
	CreatedAt               time.Time           `json:"createdAt"`
	UpdatedAt               time.Time           `json:"updatedAt"`
	DeletedAt               gorm.DeletedAt      `gorm:"index" json:"-"`

	// Relationships
	AssignedRM             *TeamMember        `gorm:"foreignKey:AssignedRMName;references:Name" json:"assignedRM,omitempty"`
	RelationshipManager    *TeamMember        `gorm:"foreignKey:RelationshipManagerName;references:Name" json:"relationshipManager,omitempty"`
	Recommendations        []Recommendation   `gorm:"foreignKey:CustomerID" json:"recommendations,omitempty"`
	ProductRecommendations []ProductRecommendation `gorm:"foreignKey:CustomerID" json:"productRecommendations,omitempty"`
	ComplianceTasks        []ComplianceTask   `gorm:"foreignKey:CustomerID" json:"complianceTasks,omitempty"`
	CustomerActions        []CustomerAction   `gorm:"foreignKey:CustomerID" json:"customerActions,omitempty"`
}

func (Customer) TableName() string {
	return "customers"
}

// SetTraits sets the customer traits with proper JSON marshaling
func (c *Customer) SetTraits(traits CustomerTraits) error {
	traitsJSON, err := json.Marshal(traits)
	if err != nil {
		return err
	}
	c.Traits = traitsJSON
	return nil
}

// GetTraits retrieves and unmarshals the customer traits
func (c *Customer) GetTraits() (*CustomerTraits, error) {
	if c.Traits == nil || len(c.Traits) == 0 {
		return nil, nil
	}

	var traits CustomerTraits
	if err := json.Unmarshal(c.Traits, &traits); err != nil {
		return nil, err
	}
	return &traits, nil
}

// UpdateTraits partially updates specific fields in traits
func (c *Customer) UpdateTraits(update func(*CustomerTraits) error) error {
	traits, err := c.GetTraits()
	if err != nil {
		return err
	}
	if traits == nil {
		traits = &CustomerTraits{}
	}

	if err := update(traits); err != nil {
		return err
	}

	return c.SetTraits(*traits)
}

// AddTraitProperty adds or updates a trait property
func (c *Customer) AddTraitProperty(key, value, propType string) error {
	return c.UpdateTraits(func(traits *CustomerTraits) error {
		found := false
		for i, prop := range traits.Properties {
			if prop.Key == key {
				traits.Properties[i].Value = value
				traits.Properties[i].Type = propType
				found = true
				break
			}
		}
		if !found {
			traits.Properties = append(traits.Properties, TraitProperty{
				Key:   key,
				Value: value,
				Type:  propType,
			})
		}
		return nil
	})
}

// GetTraitProperty gets a specific property by key
func (c *Customer) GetTraitProperty(key string) (*TraitProperty, error) {
	traits, err := c.GetTraits()
	if err != nil {
		return nil, err
	}
	if traits == nil {
		return nil, fmt.Errorf("customer has no traits")
	}

	for _, prop := range traits.Properties {
		if prop.Key == key {
			return &prop, nil
		}
	}

	return nil, fmt.Errorf("property not found: %s", key)
}

// GetPropertyValue retrieves a property value by key from traits
func (c *Customer) GetPropertyValue(key string) string {
	traits, err := c.GetTraits()
	if err != nil || traits == nil {
		return ""
	}

	for _, prop := range traits.Properties {
		if prop.Key == key {
			return prop.Value
		}
	}
	return ""
}
