package models

import (
	"time"

	"gorm.io/gorm"
)

type Jobs struct {
	gorm.Model
	EmployerID   uint
	Title        string `gorm:"not null"`
	Type         string
	Description  string `gorm:"not null"`
	WorkWhenever bool
	WorkShift    string
	Visa         bool
	Experience   string
	StartDate    time.Time
	ExactSalary  uint
	RangeSalary  RangeSalary `gorm:"embedded;embeddedPrefix:rangesalary_"`
	ExpireAt     time.Time

	EnterpriseID      uint
	EnterpriseName    string
	EnterpriseAddress string

	Crawl         bool `gorm:"default:false"`
	JobURL        string
	JobSourceName string
}

type RangeSalary struct {
	From uint
	To   uint
}

func MigrateJobs(db *gorm.DB) error {
	err := db.AutoMigrate(&Jobs{})
	return err
}
