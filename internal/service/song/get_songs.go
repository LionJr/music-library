package song

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"music-library/internal/models"
	"net/http"
	"strconv"
)

// GetSongs                godoc
// @Summary                Get songs
// @Description            Get songs by group and song with pagination, default pagination value will be 3
// @Tags                   Song
// @Accept                 json
// @Produce                json
// @Param   	           group   query     string  false       "page number in pagination"
// @Param  		           song    query     string  false       "number of elements in one page"
// @Param   	           page    query     int     false       "page number in pagination"
// @Param  		           limit   query     int     false       "number of elements in one page"
// @Success      		   200    {object}  models.GetSongsResponse
// @Failure      		   400    {object}  models.ErrorResponse
// @Failure      		   400    {object}  models.ErrorResponse
// @Failure      		   409    {object}  models.ErrorResponse
// @Failure      		   500    {object}  models.ErrorResponse
// @Router       		   /songs [get]
func (s *Service) GetSongs(ctx *gin.Context) {
	group := ctx.Param("group")
	song := ctx.Param("song")

	group = sanitizeForSQL(group)
	song = sanitizeForSQL(song)

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = models.DefaultPaginationPage
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil || limit < 1 {
		limit = models.DefaultPaginationSize
	}

	songs, totalSongCount, err := s.Repo.GetSongs(ctx, group, song, page, limit)
	if err != nil {
		s.Logger.Info("song.GetSongs", zap.Error(err))
		sendErrorResponse(ctx, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := models.GetSongsResponse{
		Songs:          songs,
		TotalSongCount: totalSongCount,
		Page:           page,
	}

	sendSuccessResponse(ctx, resp, http.StatusOK)
}
