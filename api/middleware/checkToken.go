package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"qbot/pkg/utils"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("(authorization is empty)")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 20004,
				"msg":  err.Error(),
			})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 20004,
				"msg":  err.Error(),
			})

			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("the authorization header format is not supported")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 20004,
				"msg":  err.Error(),
			})

			return
		}
		accessToken := fields[1]
		payload, err := utils.NewJWT().ParseToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 20004,
				"msg":  err.Error(),
			})
			return
		}
		c.Set(AuthorizationPayloadKey, payload)
		c.Next()
	}
}
