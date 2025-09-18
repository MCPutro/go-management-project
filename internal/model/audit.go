package model

import "time"

type Audit struct {
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy int64      `json:"created_by"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy int64      `json:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
