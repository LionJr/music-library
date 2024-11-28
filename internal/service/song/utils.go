package song

import (
	"github.com/gin-gonic/gin"
	"music-library/internal/models"
	"regexp"
)

func sendErrorResponse(ctx *gin.Context, msg string, status int) {
	resp := models.ErrorResponse{
		Message: msg,
	}

	ctx.JSON(status, resp)
}

func sendSuccessResponse(ctx *gin.Context, data interface{}, status int) {
	ctx.JSON(status, data)
}

func sanitizeForSQL(input string) string {
	re := regexp.MustCompile(`[^\w\s.,а-яА-Я]`)
	return re.ReplaceAllString(input, "")
}
