package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller example
type Controller struct {
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

type PingResponse struct {
	Message string `json:"message"`
}

// Ping godoc
// @Summary Simple Ping/Pong protocol
// @Description Receive a Ping request and send a Pong response
// @Produce  json
// @Success 200 {object} controller.PingResponse
// @Router /ping [get]
func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, PingResponse{Message: "pong"})
}
