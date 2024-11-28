package song

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"music-library/internal/models"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const layout = "02.01.2006"

// Edit godoc
// @Summary     Update song
// @Description Update song properties by song id
// @Tags        Song
// @Accept      json
// @Produce     json
// @Param       req  body     models.EditSongRequest true "Song field(s) need to be updated"
// @Param       id   path     integer                true "Song id"
// @Success     200  {object} models.SuccessResponse "Song successfully updated"
// @Failure     400  {object} models.ErrorResponse
// @Failure     404  {object} models.ErrorResponse
// @Failure     409  {object} models.ErrorResponse
// @Failure     500  {object} models.ErrorResponse
// @Router      /songs/{id} [patch]
func (s *Service) Edit(ctx *gin.Context) {
	idParam := ctx.Param("id")
	songId, err := strconv.Atoi(idParam)
	if err != nil || songId <= 0 {
		s.Logger.Info("song.Edit: ", zap.String("id", idParam))
		sendErrorResponse(ctx, "invalid song id", http.StatusBadRequest)
		return
	}

	var req models.EditSongRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		s.Logger.Info("song.Edit: unmarshal request body", zap.Error(err))
		sendErrorResponse(ctx, "invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := s.Repo.SongExists(ctx, songId)
	if err != nil {
		s.Logger.Info("song.Edit: ", zap.Error(err))
		sendErrorResponse(ctx, "internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		sendErrorResponse(ctx, "song does not exist", http.StatusNotFound)
		return
	}

	validationResult := validateInput(&req)
	if len(validationResult) > 0 {
		message := strings.Join(validationResult, "; ")
		sendErrorResponse(ctx, message, http.StatusBadRequest)
		return
	}

	err = s.Repo.Edit(ctx, songId, &req)
	if err != nil {
		s.Logger.Error("song.Edit", zap.Error(err))
		sendErrorResponse(ctx, "internal server error", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(ctx, "OK", http.StatusOK)
}

func validateInput(input *models.EditSongRequest) []string {
	validationErrors := make([]string, 0)

	if input.GroupName != nil {
		*input.GroupName = sanitizeForSQL(*input.GroupName)
	}

	if input.SongName != nil {
		*input.SongName = sanitizeForSQL(*input.SongName)
	}

	if input.ReleaseDate != nil {
		_, err := time.Parse(layout, *input.ReleaseDate)
		if err != nil {
			validationErrors = append(validationErrors, "invalid release date")
		}
	}

	if input.Link != nil {
		parsedURL, err := url.ParseRequestURI(*input.Link)
		if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
			validationErrors = append(validationErrors, "invalid link")
		}
	}

	if input.Verse != nil {
		if input.Verse.Index <= 0 {
			validationErrors = append(validationErrors, "invalid verse index")
		}

		input.Verse.Text = sanitizeForSQL(input.Verse.Text)
	}

	return validationErrors
}
