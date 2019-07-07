package viewmodels

type RoomView struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ListRoomResponse struct {
	Rooms []RoomView `json:"rooms"`
}

type CreateRoomRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateRoomResponse struct {
	RoomView
}

type GetRoomResponse struct {
	RoomView
}
