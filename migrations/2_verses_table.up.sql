CREATE TABLE verses(
    id SERIAL PRIMARY KEY,
    music_id INT REFERENCES music(id) ON DELETE CASCADE ,
    verse_text TEXT NOT NULL,
    verse_number INT NOT NULL
);