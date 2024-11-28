package song

import (
	"context"
	"music-library/internal/models"
)

type Repo interface {
	Add(ctx context.Context, song *models.Song) (int, error)
	Delete(ctx context.Context, id int) error
	Edit(ctx context.Context, id int, input *models.EditSongRequest) error
	GetSongs(ctx context.Context, group, song string, page, limit int) ([]models.Song, int, error)
	GetSongVerses(ctx context.Context, songId, page, limit int) ([]models.Verse, int, error)
	SongExists(ctx context.Context, id int) (bool, error)
	VerseExists(ctx context.Context, songId, index int) (bool, error)
}
