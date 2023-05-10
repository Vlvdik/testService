package models

type Item struct {
	ID          int    `json:"id"`
	CampaignID  int    `json:"campaign_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Removed     bool   `json:"removed"`
	CreatedAt   string `json:"created_at"`
}
