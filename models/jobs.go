package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Jobs struct {
	ID uuid.UUID `gorm:"primarykey"`

	EmployerID   string `gorm:"index:employer_id"`
	Status       string `gorm:"not null; default: REVIEW"` // REVIEW, POSTED, CLOSED
	Title        string `gorm:"not null"`
	Type         string
	WorkWhenever bool
	WorkShift    string
	Description  string `gorm:"not null"`
	Visa         bool
	Experience   uint32
	StartDate    int64
	Currency     string
	ExactSalary  uint32
	RangeSalary  string
	ExpireAt     int64

	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`

	EnterpriseID      uuid.UUID
	EnterpriseName    string
	EnterpriseAddress string

	Crawl         bool `gorm:"default: false"`
	JobURL        string
	JobSourceName string
}

func MigrateJobs(db *gorm.DB) error {
	err := db.AutoMigrate(&Jobs{})
	return err
}
