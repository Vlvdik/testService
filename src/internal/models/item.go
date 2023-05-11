package models

import "time"

type Item struct {
	ID          int       `json:"id"`
	CampaignID  int       `json:"campaign_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"created_at"`
}
