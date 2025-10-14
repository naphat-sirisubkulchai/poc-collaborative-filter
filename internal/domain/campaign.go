package domain

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CampaignStatus string

const (
	CampaignStatusDraft     CampaignStatus = "DRAFT"
	CampaignStatusActive    CampaignStatus = "ACTIVE"
	CampaignStatusPaused    CampaignStatus = "PAUSED"
	CampaignStatusCompleted CampaignStatus = "COMPLETED"
	CampaignStatusCancelled CampaignStatus = "CANCELLED"
)

type CampaignChannel string

const (
	CampaignChannelEmail       CampaignChannel = "EMAIL"
	CampaignChannelSMS         CampaignChannel = "SMS"
	CampaignChannelDirectMail  CampaignChannel = "DIRECT_MAIL"
	CampaignChannelDigital     CampaignChannel = "DIGITAL"
	CampaignChannelEvents      CampaignChannel = "EVENTS"
	CampaignChannelSocialMedia CampaignChannel = "SOCIAL_MEDIA"
)

type Campaign struct {
	ID             string          `gorm:"primarykey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name           string          `gorm:"not null" json:"name"`
	Objective      string          `gorm:"not null" json:"objective"`
	Status         CampaignStatus  `gorm:"not null" json:"status"`
	StartDate      time.Time       `gorm:"not null" json:"startDate"`
	EndDate        *time.Time      `json:"endDate"`
	Owner          string          `gorm:"not null" json:"owner"`
	Channel        CampaignChannel `gorm:"not null" json:"channel"`
	TargetAudience int             `gorm:"not null" json:"targetAudience"`
	LeadsGenerated int             `gorm:"default:0" json:"leadsGenerated"`
	ConversionRate float64         `gorm:"default:0" json:"conversionRate"`
	Budget         decimal.Decimal `gorm:"type:decimal(20,2);not null" json:"budget"`
	Spent          decimal.Decimal `gorm:"type:decimal(20,2);default:0" json:"spent"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (Campaign) TableName() string {
	return "campaigns"
}
