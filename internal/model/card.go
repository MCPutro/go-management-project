package model

type Card struct {
	ID       int64  `json:"id"`
	ListID   int64  `json:"list_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Position int    `json:"position"`
	Audit           // 👈 EMBED AUDIT STRUCT
}
