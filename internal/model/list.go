package model

type List struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	Name      string `json:"name"`
	Position  int    `json:"position"`
	Audit            // 👈 EMBED AUDIT STRUCT
}
