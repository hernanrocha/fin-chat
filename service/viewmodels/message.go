package viewmodels

import "time"

type MessageView struct {
	ID        uint      `json:"id"`
	Text      string    `json:"text"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	RoomID    uint      `json:"room_id"`
}

type CreateMessageRequest struct {
	Text string `json:"text"`
}

type CreateMessageResponse struct {
	MessageView
}

type ListMessageResponse struct {
	Messages []MessageView `json:"messages"`
}
