package models

const (
	DefaultPaginationPage = 1
	DefaultPaginationSize = 3
)

type Song struct {
	ID          int    `json:"id" db:"id"`
	GroupName   string `json:"group_name" db:"group_name"`
	SongName    string `json:"song_name" db:"song_name"`
	ReleaseDate string `json:"releaseDate" db:"release_date"`
	Link        string `json:"link" db:"link"`
	Text        string `json:"text"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}

type Verse struct {
	Id     int    `json:"id" db:"id"`
	SongId int    `json:"song_id" db:"song_id"`
	Index  int    `json:"index" db:"verse_index"`
	Text   string `json:"text" db:"text"`
}

type VerseToUpdate struct {
	Index int    `json:"index" db:"verse_index"`
	Text  string `json:"text" db:"text"`
}

type NewSongRequest struct {
	GroupName string `json:"group"`
	SongName  string `json:"song"`
}

type NewSongResponse struct {
	Message string `json:"message"`
	SongID  int    `json:"song_id"`
}

type EditSongRequest struct {
	GroupName   *string        `json:"group"`
	SongName    *string        `json:"song"`
	ReleaseDate *string        `json:"release_date"`
	Link        *string        `json:"link"`
	Verse       *VerseToUpdate `json:"verse"`
}

type GetSongVerseResponse struct {
	Verses          []Verse `json:"verses"`
	TotalVerseCount int     `json:"total_verse_count"`
	Page            int     `json:"page"`
}

type GetSongsResponse struct {
	Songs          []Song `json:"songs"`
	TotalSongCount int    `json:"total_song_count"`
	Page           int    `json:"page"`
}
