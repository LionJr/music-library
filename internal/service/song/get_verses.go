package song

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"music-library/internal/models"
	"net/http"
	"strconv"
)

// GetVerses               godoc
// @Summary                Get verses of song
// @Description            Get verses of song with pagination, default pagination value will be 3
// @Tags                   Song
// @Accept                 json
// @Produce                json
// @Param   	           id      path      int     true          "song id"
// @Param   	           page    query     int     false         "page number in pagination"
// @Param  		           limit   query     int     false         "number of elements in one page"
// @Success      		   200    {object}  models.GetSongVerseResponse
// @Failure      		   400    {object}  models.ErrorResponse
// @Failure      		   404    {object}  models.ErrorResponse
// @Failure      		   500    {object}  models.ErrorResponse
// @Router       		   /songs/{id}/verses [get]
func (s *Service) GetVerses(ctx *gin.Context) {
	idParam := ctx.Param("id")
	songId, err := strconv.Atoi(idParam)
	if err != nil || songId <= 0 {
		s.Logger.Info("song.GetVerse: ", zap.String("id", idParam))
		sendErrorResponse(ctx, "invalid song id", http.StatusBadRequest)
		return
	}

	exists, err := s.Repo.SongExists(ctx, songId)
	if err != nil {
		s.Logger.Info("song.GetVerse: ", zap.Error(err))
		sendErrorResponse(ctx, "internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		sendErrorResponse(ctx, "song does not exist", http.StatusNotFound)
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = models.DefaultPaginationPage
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil || limit < 1 {
		limit = models.DefaultPaginationSize
	}

	verses, totalVerseCount, err := s.Repo.GetSongVerses(ctx, songId, page, limit)
	if err != nil {
		s.Logger.Info("song.GetVerse", zap.Error(err))
		sendErrorResponse(ctx, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := models.GetSongVerseResponse{
		Verses:          verses,
		TotalVerseCount: totalVerseCount,
		Page:            page,
	}

	sendSuccessResponse(ctx, resp, http.StatusOK)
}
