package models

type Message struct {
	MessageID int    `db:"message_id"`
	AuthorID  int    `db:"author_id"`
	Text      string `db:"text"`
	PubDate   int64  `db:"pub_date"`
	Flagged   int    `db:"flagged"`
}
