package domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "TODO"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusCancelled  TaskStatus = "CANCELLED"
)

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "LOW"
	TaskPriorityMedium TaskPriority = "MEDIUM"
	TaskPriorityHigh   TaskPriority = "HIGH"
	TaskPriorityUrgent TaskPriority = "URGENT"
)

type Task struct {
	ID              int            `gorm:"primarykey;autoIncrement" json:"id"`
	Title           string         `gorm:"not null" json:"title"`
	Description     *string        `json:"description"`
	Status          TaskStatus     `gorm:"not null" json:"status"`
	Priority        TaskPriority   `gorm:"not null" json:"priority"`
	DueDate         *time.Time     `json:"dueDate"`
	AssigneeName    *string        `json:"assigneeName"`
	AppointmentTime *string        `json:"appointmentTime"`
	MeetingLink     datatypes.JSON `gorm:"type:jsonb" json:"meetingLink"`
	Location        datatypes.JSON `gorm:"type:jsonb" json:"location"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Assignee *TeamMember `gorm:"foreignKey:AssigneeName;references:Name" json:"assignee,omitempty"`
}

func (Task) TableName() string {
	return "tasks"
}
