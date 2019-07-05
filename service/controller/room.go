package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/hernanrocha/fin-chat/service/models"
	"github.com/hernanrocha/fin-chat/service/viewmodels"
)

// RoomController ...
type RoomController struct {
	db *gorm.DB
}

// NewRoomController ...
func NewRoomController() *RoomController {
	return &RoomController{
		db: models.GetDB(),
	}
}

// ListRooms godoc
// @Summary List Rooms
// @Description List Rooms in database
// @Tags Rooms
// @Param Authorization header string true "JWT Token"
// @Produce  json
// @Success 200 {object} viewmodels.ListRoomResponse
// @Router /api/v1/rooms [get]
func (c *RoomController) ListRooms(ctx *gin.Context) {
	var rooms []models.Room
	if err := c.db.Find(&rooms).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomList := make([]viewmodels.RoomView, len(rooms))
	for i, r := range rooms {
		roomList[i] = viewmodels.RoomView{
			ID:   r.ID,
			Name: r.Name,
		}
	}

	response := &viewmodels.ListRoomResponse{
		Rooms: roomList,
	}

	ctx.JSON(http.StatusOK, response)
}

// CreateRoom godoc
// @Summary Create Room
// @Description Create Room in database
// @Tags Rooms
// @Param Authorization header string true "JWT Token"
// @Param user body viewmodels.CreateRoomRequest true "Room Data"
// @Produce  json
// @Success 200 {object} viewmodels.CreateRoomResponse
// @Router /api/v1/rooms [post]
func (c *RoomController) CreateRoom(ctx *gin.Context) {
	var json viewmodels.CreateRoomRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room := &models.Room{
		Name: json.Name,
	}

	if err := c.db.Create(room).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := &viewmodels.CreateRoomResponse{
		viewmodels.RoomView{
			ID:   room.ID,
			Name: room.Name,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// GetRoom godoc
// @Summary Get Room
// @Description Get Room by ID
// @Tags Rooms
// @Param Authorization header string true "JWT Token"
// @Param id path int true "Room ID"
// @Produce  json
// @Success 200 {object} viewmodels.GetRoomResponse
// @Router /api/v1/rooms/{id} [get]
func (c *RoomController) GetRoom(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	var room models.Room

	if err := c.db.Where("id = ?", id).Find(&room).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := &viewmodels.GetRoomResponse{
		viewmodels.RoomView{
			ID:   room.ID,
			Name: room.Name,
		},
	}

	ctx.JSON(http.StatusOK, response)
}
