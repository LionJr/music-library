package song

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// Delete godoc
// @Summary      	     Remove song from music library
// @Description  	     Remove song from music library by song id
// @Tags         	     Song
// @Accept       	     json
// @Produce      	     json
// @Param 			     id 	             path      integer                true   "song id"
// @Success      	     200  		         {object}  string
// @Failure      	     400  			     {object}  models.ErrorResponse
// @Failure      	     404  			     {object}  models.ErrorResponse
// @Failure      	     500  			     {object}  models.ErrorResponse
// @Router       	     /songs/{id} [delete]
func (s *Service) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	songId, err := strconv.Atoi(idParam)
	if err != nil || songId <= 0 {
		s.Logger.Info("song.Delete: ", zap.String("id", idParam))
		sendErrorResponse(ctx, "invalid song id", http.StatusBadRequest)
		return
	}

	exists, err := s.Repo.SongExists(ctx, songId)
	if err != nil {
		s.Logger.Info("song.Delete: ", zap.Error(err))
		sendErrorResponse(ctx, "internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		sendErrorResponse(ctx, "song does not exist", http.StatusNotFound)
		return
	}

	err = s.Repo.Delete(ctx, songId)
	if err != nil {
		s.Logger.Info("song.Delete: ", zap.Error(err))
		sendErrorResponse(ctx, "internal server error", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(ctx, "Successfully deleted", http.StatusOK)
}
