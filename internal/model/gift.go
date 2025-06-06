package model

import "time"

type Gift struct {
	ID               string    `json:"id" db:"id"`                                           // Поле ID → столбец "id"
	StarCount        int       `json:"star_count" db:"star_count"`                           // Поле StarCount → столбец "star_count"
	UpgradeStarCount int       `json:"upgrade_star_count,omitempty" db:"upgrade_star_count"` // Необязательное поле
	TotalCount       int       `json:"total_count,omitempty" db:"total_count"`               // Необязательное поле
	RemainingCount   int       `json:"remaining_count,omitempty" db:"remaining_count"`       // Необязательное поле
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}
