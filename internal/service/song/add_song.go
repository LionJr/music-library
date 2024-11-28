package song

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"music-library/internal/models"
	"net/http"
)

// Add                     godoc
// @Summary                Adding a new song
// @Description            Adding a new song if it is not already existing one
// @Tags                   Song
// @Accept                 json
// @Produce                json
// @Param req              body   models.NewSongRequest true  "song information to add"
// @Success      		   200    {object}  models.NewSongResponse
// @Failure      		   400    {object}  models.ErrorResponse
// @Failure      		   409    {object}  models.ErrorResponse
// @Failure      		   500    {object}  models.ErrorResponse
// @Router       		   /songs [post]
func (s *Service) Add(ctx *gin.Context) {
	var req models.NewSongRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.Logger.Info("song.Add: unmarshal request body", zap.Error(err))
		sendErrorResponse(ctx, "invalid request body", http.StatusBadRequest)
		return
	}

	groupName := sanitizeForSQL(req.GroupName)
	songName := sanitizeForSQL(req.SongName)

	apiURL := fmt.Sprintf("%s?group=%s&song=%s", s.config.ExternalAPI.URL, groupName, songName)
	songInfoResp, err := http.Get(apiURL)
	if err != nil || songInfoResp.StatusCode != http.StatusOK {
		s.Logger.Info("song.Add: fetch song details error", zap.Error(err), zap.Int("status_code", songInfoResp.StatusCode))
		sendErrorResponse(ctx, "internal server error", http.StatusInternalServerError)
		return
	}
	defer songInfoResp.Body.Close()

	var song models.Song
	if err = jsoniter.NewDecoder(songInfoResp.Body).Decode(&song); err != nil {
		s.Logger.Info("song.Add: parse api response error", zap.Error(err))
		sendErrorResponse(ctx, "internal server error", http.StatusInternalServerError)
		return
	}

	song.GroupName, song.SongName = groupName, songName

	songId, err := s.Repo.Add(ctx, &song)
	if err != nil {
		s.Logger.Info("song.Add: ", zap.Error(err))
		if err.Error() == "song already exists" {
			sendErrorResponse(ctx, "song already exists", http.StatusConflict)
		} else {
			sendErrorResponse(ctx, "song add error", http.StatusInternalServerError)
		}
		return
	}

	resp := models.NewSongResponse{
		Message: "Song successfully added",
		SongID:  songId,
	}

	sendSuccessResponse(ctx, resp, http.StatusOK)
}
