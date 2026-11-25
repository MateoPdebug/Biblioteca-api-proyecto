package handler

import (
	"biblioteca-api/internal/models"
	"biblioteca-api/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

type LaptopHandler struct{ service *service.LaptopService }

func NewLaptopHandler(s *service.LaptopService) *LaptopHandler { return &LaptopHandler{s} }
func laptopError(c *gin.Context, err error) {
	status := 500
	message := "error interno"
	switch {
	case errors.Is(err, service.ErrLaptopNotFound):
		status = 404
	case errors.Is(err, service.ErrLaptopForbidden):
		status = 403
	case errors.Is(err, service.ErrLaptopInvalid):
		status = 400
	case errors.Is(err, service.ErrLaptopUnavailable), errors.Is(err, service.ErrLaptopReturned), errors.Is(err, service.ErrLaptopSerial):
		status = 409
	}
	if status != 500 {
		message = err.Error()
	}
	c.JSON(status, gin.H{"error": message})
}
func (h *LaptopHandler) Create(c *gin.Context) {
	var req models.LaptopRequest
	if c.ShouldBindJSON(&req) != nil {
		c.JSON(400, gin.H{"error": "datos inválidos"})
		return
	}
	l, err := h.service.Create(req)
	if err != nil {
		laptopError(c, err)
		return
	}
	c.JSON(201, l)
}
func (h *LaptopHandler) List(c *gin.Context) {
	l, err := h.service.List()
	if err != nil {
		laptopError(c, err)
		return
	}
	c.JSON(200, l)
}
func (h *LaptopHandler) Borrow(c *gin.Context) {
	var req models.LaptopLoanRequest
	if c.ShouldBindJSON(&req) != nil {
		c.JSON(400, gin.H{"error": "laptop_id inválido"})
		return
	}
	l, err := h.service.Borrow(c.GetInt("user_id"), req.LaptopID)
	if err != nil {
		laptopError(c, err)
		return
	}
	c.JSON(201, l)
}
func (h *LaptopHandler) MyLoans(c *gin.Context) {
	l, err := h.service.MyLoans(c.GetInt("user_id"))
	if err != nil {
		laptopError(c, err)
		return
	}
	c.JSON(200, l)
}
func (h *LaptopHandler) Return(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "id inválido"})
		return
	}
	if err = h.service.Return(c.GetInt("user_id"), c.GetString("role"), id); err != nil {
		laptopError(c, err)
		return
	}
	c.JSON(200, gin.H{"message": "laptop devuelta correctamente"})
}
