package domain

type Tweet struct {
	TweetId  int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}
