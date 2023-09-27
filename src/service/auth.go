package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"me.daily/src/bundle"
	"me.daily/src/token"
)

func (s *Service) checkAuth(c *gin.Context) {
	path := c.Request.URL.Path
	if c.Request.Method == "POST" {
		switch path {
		case "/api/login":
			fallthrough
		case "/api/user":
			c.Next()
			return
		}
	}

	authStr, err := c.Cookie("Authorization")

	if err != nil {
		c.JSON(http.StatusOK, bundle.ErrorResponse{
			Code: bundle.CodeToken,
		})
		c.Abort()
	} else {
		auth, err := token.PareToken(authStr)
		if err != nil {
			c.JSON(http.StatusOK, bundle.ErrorResponse{
				Code: bundle.CodeToken,
			})
			c.Abort()
		} else {
			if _, ok := s.c.Get(strconv.Itoa(auth.UserId)); ok {
				c.Set("user_id", auth.UserId)
				c.Next()
			} else {
				c.JSON(http.StatusOK, bundle.ErrorResponse{
					Code: bundle.CodeToken,
				})
				c.Abort()
			}
		}
	}
}
