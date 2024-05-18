package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nobitayon/memorization-app/account/handler/model"
	"github.com/nobitayon/memorization-app/account/handler/model/apperrors"
)

type SignUpReq struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

func (h *Handler) Signup(c *gin.Context) {
	var req SignUpReq

	if ok := bindData(c, &req); !ok {
		return
	}

	u := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.UserService.Signup(c, u)
	if err != nil {
		log.Printf("failed to sign up user: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
}
