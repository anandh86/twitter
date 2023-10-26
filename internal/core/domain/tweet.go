package domain

type Tweet struct {
	ID     int    `json:"id"`
	Body   string `json:"body"`
	Author int    `json:"author_id"`
}
