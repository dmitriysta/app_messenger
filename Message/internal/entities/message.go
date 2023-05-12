package entities

import "time"

type Message struct {
	Id        int
	UserID    int
	ChannelID int
	Content   string
	CreatedAt time.Time
}
