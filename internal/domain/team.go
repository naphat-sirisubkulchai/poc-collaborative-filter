package domain

import (
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TeamMemberStatus string

const (
	TeamMemberStatusActive   TeamMemberStatus = "ACTIVE"
	TeamMemberStatusInactive TeamMemberStatus = "INACTIVE"
	TeamMemberStatusOnLeave  TeamMemberStatus = "ON_LEAVE"
)

type PerformanceRating string

const (
	PerformanceRatingExcellent    PerformanceRating = "EXCELLENT"
	PerformanceRatingGood         PerformanceRating = "GOOD"
	PerformanceRatingAverage      PerformanceRating = "AVERAGE"
	PerformanceRatingBelowAverage PerformanceRating = "BELOW_AVERAGE"
	PerformanceRatingPoor         PerformanceRating = "POOR"
)

type TeamMember struct {
	ID                    string            `gorm:"primarykey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name                  string            `gorm:"uniqueIndex;not null" json:"name"`
	Email                 string            `gorm:"uniqueIndex;not null" json:"email"`
	Phone                 *string           `json:"phone"`
	Role                  string            `gorm:"not null" json:"role"`
	Department            string            `gorm:"not null" json:"department"`
	ManagedAUM            decimal.Decimal   `gorm:"type:decimal(20,2);default:0" json:"managedAUM"`
	ClientsManaged        int               `gorm:"default:0" json:"clientsManaged"`
	Status                TeamMemberStatus  `gorm:"not null" json:"status"`
	JoinDate              time.Time         `gorm:"not null" json:"joinDate"`
	Specialization        pq.StringArray    `gorm:"type:text[]" json:"specialization"`
	PerformanceRating     PerformanceRating `gorm:"not null" json:"performanceRating"`
	LastActivity          *time.Time        `json:"lastActivity"`
	TotalLeads            int               `gorm:"default:0" json:"totalLeads"`
	MonthlyAUMGoal        decimal.Decimal   `gorm:"type:decimal(20,2);default:0" json:"monthlyAUMGoal"`
	QuarterlyAUMGoal      decimal.Decimal   `gorm:"type:decimal(20,2);default:0" json:"quarterlyAUMGoal"`
	YearlyAUMGoal         decimal.Decimal   `gorm:"type:decimal(20,2);default:0" json:"yearlyAUMGoal"`
	MonthlyCommission     decimal.Decimal   `gorm:"type:decimal(20,2);default:0" json:"monthlyCommission"`
	QuarterlyCommission   decimal.Decimal   `gorm:"type:decimal(20,2);default:0" json:"quarterlyCommission"`
	YearlyCommission      decimal.Decimal   `gorm:"type:decimal(20,2);default:0" json:"yearlyCommission"`
	GoalAchievementRate   float64           `gorm:"default:0" json:"goalAchievementRate"`
	MonthlyLeadGoal       int               `gorm:"default:0" json:"monthlyLeadGoal"`
	QuarterlyLeadGoal     int               `gorm:"default:0" json:"quarterlyLeadGoal"`
	YearlyLeadGoal        int               `gorm:"default:0" json:"yearlyLeadGoal"`
	ConvertedLeads        int               `gorm:"default:0" json:"convertedLeads"`
	CreatedAt             time.Time         `json:"createdAt"`
	UpdatedAt             time.Time         `json:"updatedAt"`
	DeletedAt             gorm.DeletedAt    `gorm:"index" json:"-"`
}

func (TeamMember) TableName() string {
	return "team_members"
}
