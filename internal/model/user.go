package model

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"` // jangan dikirim ke frontend
	Audit           // 👈 EMBED AUDIT STRUCT
}
