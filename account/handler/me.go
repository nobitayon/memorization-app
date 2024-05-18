package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nobitayon/memorization-app/account/handler/model"
	"github.com/nobitayon/memorization-app/account/handler/model/apperrors"
)

// user detail
func (h *Handler) Me(c *gin.Context) {

	// will be ripped out later to some middleware file
	// just for mocking purpose
	user, exists := c.Get("user")
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	uid := user.(*model.User).UID

	u, err := h.UserService.Get(c, uid)
	if err != nil {
		log.Printf("Unable to find user: %v\n%v", uid, err)
		e := apperrors.NewNotFound("user", uid.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}
