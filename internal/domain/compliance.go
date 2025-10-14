package domain

import (
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CustomerActionType string

const (
	CustomerActionTypeKYCUpdate          CustomerActionType = "KYC_UPDATE"
	CustomerActionTypeDocumentSubmission CustomerActionType = "DOCUMENT_SUBMISSION"
	CustomerActionTypePortfolioReview    CustomerActionType = "PORTFOLIO_REVIEW"
	CustomerActionTypeRiskAssessment     CustomerActionType = "RISK_ASSESSMENT"
	CustomerActionTypeRebalancing        CustomerActionType = "REBALANCING"
)

type ActionStatus string

const (
	ActionStatusPending    ActionStatus = "PENDING"
	ActionStatusInProgress ActionStatus = "IN_PROGRESS"
	ActionStatusCompleted  ActionStatus = "COMPLETED"
	ActionStatusCancelled  ActionStatus = "CANCELLED"
)

type ComplianceTaskType string

const (
	ComplianceTaskTypeKYCReview            ComplianceTaskType = "KYC_REVIEW"
	ComplianceTaskTypeAMLCheck             ComplianceTaskType = "AML_CHECK"
	ComplianceTaskTypeRiskAssessment       ComplianceTaskType = "RISK_ASSESSMENT"
	ComplianceTaskTypeDocumentVerification ComplianceTaskType = "DOCUMENT_VERIFICATION"
)

type CustomerAction struct {
	ID                string             `gorm:"primarykey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CustomerID        string             `gorm:"type:uuid;not null" json:"customerId"`
	Type              CustomerActionType `gorm:"not null" json:"type"`
	Title             string             `gorm:"not null" json:"title"`
	Description       *string            `json:"description"`
	DueDate           *time.Time         `json:"dueDate"`
	Priority          TaskPriority       `gorm:"not null" json:"priority"`
	Status            ActionStatus       `gorm:"not null" json:"status"`
	AssignedToName    *string            `json:"assignedToName"`
	EstimatedTime     *string            `json:"estimatedTime"`
	RequiredDocuments pq.StringArray     `gorm:"type:text[]" json:"requiredDocuments"`
	CompletedSteps    int                `gorm:"default:0" json:"completedSteps"`
	TotalSteps        int                `gorm:"not null" json:"totalSteps"`
	CreatedAt         time.Time          `json:"createdAt"`
	UpdatedAt         time.Time          `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt     `gorm:"index" json:"-"`

	// Relationships
	AssignedTo *TeamMember `gorm:"foreignKey:AssignedToName;references:Name" json:"assignedTo,omitempty"`
	Customer   Customer    `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (CustomerAction) TableName() string {
	return "customer_actions"
}

type ComplianceTask struct {
	ID             string             `gorm:"primarykey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CustomerID     string             `gorm:"type:uuid;not null" json:"customerId"`
	Type           ComplianceTaskType `gorm:"not null" json:"type"`
	Title          string             `gorm:"not null" json:"title"`
	Description    *string            `json:"description"`
	DueDate        *time.Time         `json:"dueDate"`
	Status         TaskStatus         `gorm:"not null" json:"status"`
	AssignedToName *string            `json:"assignedToName"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt     `gorm:"index" json:"-"`

	// Relationships
	AssignedTo *TeamMember `gorm:"foreignKey:AssignedToName;references:Name" json:"assignedTo,omitempty"`
	Customer   Customer    `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (ComplianceTask) TableName() string {
	return "compliance_tasks"
}

type ProductRecommendation struct {
	ID             string          `gorm:"primarykey;type:uuid;default:uuid_generate_v4()" json:"id"`
	ProductID      string          `gorm:"not null" json:"productId"`
	ProductName    string          `gorm:"not null" json:"productName"`
	ProductType    string          `gorm:"not null" json:"productType"`
	ExpectedReturn decimal.Decimal `gorm:"type:decimal(20,2);not null" json:"expectedReturn"`
	RiskScore      int             `gorm:"not null" json:"riskScore"`
	MatchScore     int             `gorm:"not null" json:"matchScore"`
	Reasoning      string          `gorm:"not null" json:"reasoning"`
	CustomerID     string          `gorm:"type:uuid;not null" json:"customerId"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`

	// Relationships
	Customer Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (ProductRecommendation) TableName() string {
	return "product_recommendations"
}
