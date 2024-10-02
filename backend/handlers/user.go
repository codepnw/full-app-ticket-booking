package handlers

import (
	"net/http"
	"strconv"

	"github.com/codepnw/ticket-api/models"
	userService "github.com/codepnw/ticket-api/services/user"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userSrv userService.IUserService
}

func NewUserHandler(userSrv userService.IUserService) *userHandler {
	return &userHandler{userSrv: userSrv}
}

func (h *userHandler) SignupUser(c *gin.Context) {
	req := new(models.UserRegisterReq)

	if err := c.ShouldBindJSON(req); err != nil {
		errorBadRequest(c, err.Error())
		return
	}

	result, err := h.userSrv.CreateUser(req)
	if err != nil {
		errorInternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *userHandler) SignupAdmin(c *gin.Context) {
	req := new(models.UserRegisterReq)

	if err := c.ShouldBindJSON(req); err != nil {
		errorBadRequest(c, err.Error())
		return
	}

	result, err := h.userSrv.CreateAdmin(req)
	if err != nil {
		errorInternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *userHandler) GetProfile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("user_id"))

	user, err := h.userSrv.GetProfile(uint(id))
	if err != nil {
		errorInternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}
