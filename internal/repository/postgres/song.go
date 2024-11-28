package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"music-library/internal/models"
	"strings"
)

type SongRepository struct {
	db *sqlx.DB
}

func NewSongRepository(db *sqlx.DB) *SongRepository {
	return &SongRepository{db: db}
}

func (m *SongRepository) Add(ctx context.Context, song *models.Song) (int, error) {
	var (
		id     int
		exists bool
	)

	query := `SELECT EXISTS(SELECT id FROM songs WHERE group_name = $1 AND song_name = $2)`
	err := m.db.QueryRowContext(ctx, query, song.GroupName, song.SongName).Scan(&exists)
	if err != nil {
		return id, err
	}

	if exists {
		return id, errors.New("song already exists")
	}

	query = `INSERT INTO songs(group_name, song_name, release_date, link) VALUES ($1, $2, $3, $4) RETURNING id`
	tx, _ := m.db.Begin()
	if err := tx.QueryRowContext(ctx, query,
		song.GroupName,
		song.SongName,
		song.ReleaseDate,
		song.Link,
	).Scan(&id); err != nil {
		_ = tx.Rollback()
		return id, err
	}

	verseQuery := `INSERT INTO song_verses(song_id, verse_index, text) VALUES ($1, $2, $3)`

	verses := strings.Split(song.Text, "\n\n")
	for index := range verses {
		if _, err := tx.ExecContext(ctx, verseQuery, id, index+1, verses[index]); err != nil {
			_ = tx.Rollback()
			return id, err
		}
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return id, err
	}

	return id, nil
}

func (m *SongRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE 
			  FROM songs 
			  WHERE id = $1`
	_, err := m.db.ExecContext(ctx, query, id)
	return err
}

func (m *SongRepository) Edit(ctx context.Context, id int, input *models.EditSongRequest) error {
	var (
		conditions      []string
		verseConditions []string
		args            []interface{}
		verseArgs       []interface{}
	)

	if input.GroupName != nil {
		args = append(args, *input.GroupName)
		conditions = append(conditions, fmt.Sprintf("group_name = $%d", len(args)))
	}

	if input.SongName != nil {
		args = append(args, *input.SongName)
		conditions = append(conditions, fmt.Sprintf("song_name = $%d", len(args)))
	}

	if input.ReleaseDate != nil {
		args = append(args, *input.ReleaseDate)
		conditions = append(conditions, fmt.Sprintf("release_date = $%d", len(args)))
	}

	if input.Link != nil {
		args = append(args, *input.Link)
		conditions = append(conditions, fmt.Sprintf("link = $%d", len(args)))
	}

	if input.Verse != nil {
		verseExists, err := m.VerseExists(ctx, id, input.Verse.Index)
		if err != nil {
			return err
		}

		if !verseExists {
			return errors.New("no verse found with provided index")
		}

		verseArgs = append(verseArgs, input.Verse.Text)
		verseConditions = append(verseConditions, fmt.Sprintf("text = $%d", len(verseArgs)))
	}

	if len(args) > 0 {
		args = append(args, id)
		query := `UPDATE songs SET` + " " + strings.Join(conditions, ", ") + fmt.Sprintf(" WHERE id = $%d", len(args))

		_, err := m.db.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to execute update query for song with id - %d: %w", id, err)
		}
	}

	if len(verseArgs) > 0 {
		verseArgs = append(verseArgs, id, input.Verse.Index)
		verseQuery := `UPDATE song_verses SET` + " " + strings.Join(verseConditions, ", ") + fmt.Sprintf(" WHERE song_id = $%d AND verse_index = $%d", len(verseArgs)-1, len(verseArgs))

		_, err := m.db.ExecContext(ctx, verseQuery, verseArgs...)
		if err != nil {
			return fmt.Errorf("failed to execute verse update query for song with id - %d: %w", id, err)
		}
	}

	return nil
}

func (m *SongRepository) GetSongs(ctx context.Context, group, song string, page, limit int) ([]models.Song, int, error) {
	var (
		songs      []models.Song
		conditions []string
		args       []interface{}
		totalCount int
	)

	query := `SELECT s.id, s.group_name, s.song_name, 
                     s.release_date, s.link, s.created_at, s.updated_at 
			  FROM songs AS s`

	if group != "" {
		args = append(args, group)
		conditions = append(conditions, fmt.Sprintf("s.group_name = $%d", len(args)))
	}

	if song != "" {
		args = append(args, song)
		conditions = append(conditions, fmt.Sprintf("s.song_name = $%d", len(args)))
	}

	if len(conditions) > 0 {
		query += ` WHERE ` + strings.Join(conditions, " AND ")
	}

	offset := (page - 1) * limit
	args = append(args, limit, offset)
	query += fmt.Sprintf(" ORDER BY s.id LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	err := m.db.SelectContext(ctx, &songs, query, args...)
	if err != nil {
		return nil, totalCount, err
	}

	countQuery := `SELECT COUNT(songs.id) FROM songs`
	if len(conditions) > 0 {
		countQuery += ` AND ` + strings.Join(conditions, ` AND `)
	}

	err = m.db.GetContext(ctx, &totalCount, countQuery, args[:len(conditions)]...)
	if err != nil {
		return nil, totalCount, err
	}

	return songs, totalCount, nil
}

func (m *SongRepository) GetSongVerses(ctx context.Context, songId, page, limit int) ([]models.Verse, int, error) {
	var (
		verses     []models.Verse
		totalCount int
	)

	query := `SELECT sv.id, sv.song_id, sv.verse_index, sv.text 
              FROM song_verses AS sv
              WHERE sv.song_id = $1
              ORDER BY sv.verse_index LIMIT $2 OFFSET $3`

	offset := (page - 1) * limit

	err := m.db.SelectContext(ctx, &verses, query, songId, limit, offset)
	if err != nil {
		return nil, totalCount, err
	}

	countQuery := `SELECT COUNT(sv.id) FROM song_verses AS sv WHERE sv.song_id = $1`

	err = m.db.GetContext(ctx, &totalCount, countQuery, songId)
	if err != nil {
		return nil, totalCount, err
	}

	return verses, totalCount, nil
}

func (m *SongRepository) SongExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT id 
    				 		FROM songs 
    				 		WHERE id = $1)`
	err := m.db.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *SongRepository) VerseExists(ctx context.Context, songId, index int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT sv.id 
    				 		FROM song_verses AS sv
    				 		WHERE sv.song_id = $1 AND sv.verse_index = $2)`
	err := m.db.QueryRowContext(ctx, query, songId, index).Scan(&exists)
	return exists, err
}
