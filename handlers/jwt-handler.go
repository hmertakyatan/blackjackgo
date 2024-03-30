package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hmertakyatan/blackjackgo/services"
)

type JwtHandler struct {
	jwtService *services.JwtService
}

func NewJwtHandler(jwtService *services.JwtService) *JwtHandler {
	return &JwtHandler{
		jwtService: jwtService,
	}
}

func (jh *JwtHandler) GenerateTokenFromPayload(c *gin.Context) {
	var payload map[string]interface{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
		return
	}

	generatedToken, err := jh.jwtService.GenerateToken(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
		return
	}
	c.JSON(http.StatusOK, generatedToken)

}
