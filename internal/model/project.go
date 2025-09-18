package model

type Project struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Audit              // 👈 EMBED AUDIT STRUCT
}
