package viewmodels

type RoomView struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type RegisterRequest struct {
	UserView
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	UserView
}

type ListRoomResponse struct {
	Rooms []RoomView `json:"rooms"`
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

type CreateRoomResponse struct {
	RoomView
}

type GetRoomResponse struct {
	RoomView
}

type ListMessageResponse struct {
	Messages []MessageView `json:"messages"`
}
