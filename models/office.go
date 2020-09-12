package models

// Office struct used for storing offices
type Office struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsGlobal bool   `json:"is_global"`
	IsActive bool   `json:"is_active"`
}
