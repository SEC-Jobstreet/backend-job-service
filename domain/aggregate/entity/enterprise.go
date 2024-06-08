package entity

import "github.com/google/uuid"

type Enterprise struct {
	ID uuid.UUID `json:"id"`

	Name      string `json:"name"`
	Country   string `json:"country"`
	Address   string `json:"address"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Field     string `json:"field"`
	Size      string `json:"size"`
	Url       string `json:"url"`
	License   string `json:"license"`

	EmployerID   string `json:"employer_id"`
	EmployerRole string `json:"employer_role"`
}
