CREATE TABLE songs (
   id SERIAL PRIMARY KEY,
   group_name VARCHAR(100) NOT NULL,
   song_name VARCHAR(200) NOT NULL,
   release_date VARCHAR(50) NOT NULL,
   link VARCHAR(255) NOT NULL,
   created_at TIMESTAMP DEFAULT NOW(),
   updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE song_verses (
    id SERIAL PRIMARY KEY,
    song_id INTEGER REFERENCES songs(id) ON DELETE CASCADE,
    verse_index INTEGER NOT NULL,
    text TEXT NOT NULL
);