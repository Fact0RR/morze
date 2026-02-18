package domain

type Contact struct {
    ContactID   int      `json:"contact_id"`
    UserID      int      `json:"user_id"`
    Data        string   `json:"data"`
    Additionals []string `json:"additionals,omitempty"`
}