package domain

import (
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type LeadPriority string

const (
	LeadPriorityHigh   LeadPriority = "HIGH"
	LeadPriorityMedium LeadPriority = "MEDIUM"
	LeadPriorityLow    LeadPriority = "LOW"
)

type LeadStatus string

const (
	LeadStatusNew       LeadStatus = "NEW"
	LeadStatusContacted LeadStatus = "CONTACTED"
	LeadStatusQualified LeadStatus = "QUALIFIED"
	LeadStatusConverted LeadStatus = "CONVERTED"
	LeadStatusLost      LeadStatus = "LOST"
)

type Lead struct {
	ID             string          `gorm:"primarykey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CustomerName   string          `gorm:"not null" json:"customerName"`
	Segment        string          `gorm:"not null" json:"segment"`
	AUM            decimal.Decimal `gorm:"type:decimal(20,2);not null" json:"aum"`
	Priority       LeadPriority    `gorm:"not null" json:"priority"`
	Status         LeadStatus      `gorm:"not null" json:"status"`
	AssignedRMName *string         `json:"assignedRMName"`
	ContactMethod  *string         `json:"contactMethod"`
	LastContact    *time.Time      `json:"lastContact"`
	CustomerType   pq.StringArray  `gorm:"type:text[]" json:"customerType"`
	Campaign       *string         `json:"campaign"`
	Email          *string         `json:"email"`
	Phone          *string         `json:"phone"`
	CustomerID     *string         `gorm:"type:uuid" json:"customerId"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`

	// Relationships
	AssignedRM *TeamMember `gorm:"foreignKey:AssignedRMName;references:Name" json:"assignedRM,omitempty"`
	Customer   *Customer   `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (Lead) TableName() string {
	return "leads"
}
