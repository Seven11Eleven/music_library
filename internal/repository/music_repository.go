package repository

import (
	"context"
	"fmt"
	"github.com/Seven11Eleven/music_library/internal/domain/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"strings"
)

type musicRepository struct {
	pool *pgxpool.Pool
}

func (m musicRepository) GetMusic(ctx context.Context, musicName, groupName string) (*models.Music, error) {
	query := `
		SELECT 
			m.id, m.title, m.group_name, m.release_date, m.link, 
			v.verse_text, v.verse_number
		FROM 
			music m
		LEFT JOIN 
			verses v ON m.id = v.music_id
		WHERE 
			m.title = $1 AND m.group_name = $2
		ORDER BY 
			v.verse_number;
	`

	rows, err := m.pool.Query(ctx, query, musicName, groupName)
	if err != nil {
		log.Printf("Error querying music: %v", err)
		return nil, err
	}
	defer rows.Close()

	var music models.Music
	var verses []models.Verse
	isFirstRow := true

	for rows.Next() {
		var verse models.Verse
		if isFirstRow {
			err := rows.Scan(&music.ID, &music.SongName, &music.GroupName, &music.ReleaseDate, &music.Link, &verse.Text, &verse.Number)
			if err != nil {
				log.Printf("Error scanning row: %v", err)
				return nil, err
			}
			isFirstRow = false
		} else {
			err := rows.Scan(nil, nil, nil, nil, nil, &verse.Text, &verse.Number)
			if err != nil {
				log.Printf("Error scanning verse row: %v", err)
				return nil, err
			}
		}
		verses = append(verses, verse)
	}

	if isFirstRow {
		return nil, nil
	}

	music.Verses = verses

	return &music, nil
}

func (m musicRepository) SaveMusic(ctx context.Context, music *models.Music) (*models.Music, error) {
	log.Infof("Saving new music: %s by %s", music.SongName, music.GroupName)

	var musicID int
	query := `
	INSERT INTO music (release_date, title, group_name, link) 
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	tx, err := m.pool.Begin(ctx)
	if err != nil {
		log.Errorf("Error beginning transaction: %v", err)
		return nil, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil && err != pgx.ErrTxClosed {
			log.Warnf("Error rolling back transaction: %v", err)
		}
	}(tx, ctx)

	err = tx.QueryRow(ctx, query, music.ReleaseDate, music.SongName, music.GroupName, music.Link).Scan(&musicID)
	if err != nil {
		log.Errorf("Error saving music: %v", err)
		return nil, err
	}

	log.Infof("Music saved with ID: %d", musicID)

	values := []string{}
	args := []interface{}{}
	placeholdersIndex := 1

	for _, verse := range music.Verses {
		values = append(values, fmt.Sprintf("($%d,$%d,$%d)", placeholdersIndex, 1+placeholdersIndex, placeholdersIndex+2))
		args = append(args, musicID, verse.Text, placeholdersIndex/3+1)
		placeholdersIndex += 3
	}

	query = fmt.Sprintf("INSERT INTO verses (music_id, verse_text, verse_number) VALUES %s", strings.Join(values, ", "))

	log.Infof("Saving %d verses for music ID %d", len(music.Verses), musicID)

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		log.Errorf("Error saving music verses: %v", err)
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return nil, err
	}

	log.Infof("Music and verses saved successfully for ID: %d", musicID)
	return music, nil
}

func (m musicRepository) GetMusicsByFilters(ctx context.Context, filters models.MusicFilters, page, pageSize int) ([]models.Music, error) {
	log.Infof("Fetching music list with filters: %+v", filters)
	query := `
				SELECT 
				    	id, release_date, title, group_name, link
				FROM 
				    	music
				WHERE
				    	1=1
					AND
					    	($1::DATE IS NULL OR release_date = $1::DATE)
					AND	
					    	($2::TEXT IS NULL OR title ILIKE '%' || $2::TEXT || '%')
					AND	
					    	($3::TEXT IS NULL OR group_name ILIKE '%' || $3::TEXT || '%')
					AND
					    	($4::TEXT IS NULL OR link = $4::TEXT)
				LIMIT $5 OFFSET $6	
`
	offset := (page - 1) * pageSize

	rows, err := m.pool.Query(ctx, query, filters.ReleaseDate, filters.SongName, filters.GroupName, filters.Link, pageSize, offset)
	if err != nil {
		log.Errorf("Error fetching music list: %v", err)
		return nil, err
	}
	defer rows.Close()

	var musics []models.Music
	for rows.Next() {
		var music models.Music
		err := rows.Scan(&music.ID, &music.ReleaseDate, &music.SongName, &music.GroupName, &music.Link)
		if err != nil {
			log.Errorf("Error scanning music row: %v", err)
			return nil, err
		}
		musics = append(musics, music)
	}

	log.Infof("Successfully fetched %d music records", len(musics))
	return musics, nil
}

func (m musicRepository) GetMusicTextWithPaginationByVerse(ctx context.Context, musicID string, limit, offset int) (*models.Music, error) {
	log.Infof("Fetching verses for music ID: %s with pagination limit %d, offset %d", musicID, limit, offset)
	query := `
        SELECT 
            m.id AS music_id,
            m.title,
            m.release_date,
            m.group_name,
            m.link,
            v.verse_text,
            v.verse_number
        FROM 
            music m
        JOIN 
            verses v ON m.id = v.music_id
        WHERE 
            m.id = $1
        ORDER BY 
            v.verse_number
        LIMIT $2 OFFSET $3;
    `

	rows, err := m.pool.Query(ctx, query, musicID, limit, offset)
	if err != nil {
		log.Errorf("Error fetching verses: %v", err)
		return nil, err
	}
	defer rows.Close()

	var musicVerses models.Music
	var verses []models.Verse

	for rows.Next() {
		var verse models.Verse
		if err := rows.Scan(&musicVerses.ID, &musicVerses.SongName, &musicVerses.ReleaseDate, &musicVerses.GroupName, &musicVerses.Link, &verse.Text, &verse.Number); err != nil {
			log.Errorf("Error scanning verse row: %v", err)
			return nil, err
		}
		verses = append(verses, verse)
	}

	musicVerses.Verses = verses
	log.Infof("Successfully fetched %d verses for music ID: %s", len(verses), musicID)
	return &musicVerses, nil
}

func (m musicRepository) DeleteMusic(ctx context.Context, musicID string) error {
	log.Infof("Deleting music with ID: %s", musicID)
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		log.Errorf("Error beginning transaction for delete: %v", err)
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil && err != pgx.ErrTxClosed {
			log.Warnf("Error rolling back transaction: %v", err)
		}
	}(tx, ctx)

	query := `DELETE FROM music WHERE id = $1;`

	_, err = tx.Exec(ctx, query, musicID)
	if err != nil {
		log.Errorf("Error deleting music with ID %s: %v", musicID, err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Errorf("Error committing delete transaction: %v", err)
		return err
	}

	log.Infof("Music with ID %s deleted successfully", musicID)
	return nil
}

func (m musicRepository) UpdateMusic(ctx context.Context, music models.Music) (models.Music, error) {
	log.Infof("Updating music with ID: %s", music.ID)
	// Обновление данных музыки
	query := `UPDATE music SET `
	params := []interface{}{}
	paramCount := 1

	if music.ReleaseDate != nil {
		query += fmt.Sprintf("release_date = $%d, ", paramCount)
		paramCount++
		params = append(params, music.ReleaseDate)
	}
	if music.Link != "" {
		query += fmt.Sprintf("link = $%d, ", paramCount)
		paramCount++
		params = append(params, music.Link)
	}
	if music.SongName != "" {
		query += fmt.Sprintf("title = $%d, ", paramCount)
		paramCount++
		params = append(params, music.SongName)
	}
	if music.GroupName != "" {
		query += fmt.Sprintf("group_name = $%d, ", paramCount)
		paramCount++
		params = append(params, music.GroupName)
	}

	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, release_date, title, group_name, link", paramCount)
	params = append(params, music.ID)

	var updatedMusic models.Music

	err := m.pool.QueryRow(ctx, query, params...).Scan(
		&updatedMusic.ID,
		&updatedMusic.ReleaseDate,
		&updatedMusic.SongName,
		&updatedMusic.GroupName,
		&updatedMusic.Link,
	)
	if err != nil {
		log.Errorf("Error updating music with ID %s: %v", music.ID, err)
		return models.Music{}, err
	}

	log.Infof("Music with ID %s updated successfully", music.ID)

	if len(music.Verses) > 0 {
		for _, verse := range music.Verses {
			verseUpdateQuery := `UPDATE verses SET verse_text = $1 WHERE music_id = $2 AND verse_number = $3`
			_, err := m.pool.Exec(ctx, verseUpdateQuery, verse.Text, music.ID, verse.Number)
			if err != nil {
				log.Errorf("Error updating verse for music ID %s: %v", music.ID, err)
				return models.Music{}, err
			}
		}
	}

	versesQuery := `SELECT verse_text, verse_number FROM verses WHERE music_id = $1 ORDER BY verse_number`
	rows, err := m.pool.Query(ctx, versesQuery, updatedMusic.ID)
	if err != nil {
		log.Errorf("Error fetching updated verses for music ID %s: %v", updatedMusic.ID, err)
		return models.Music{}, err
	}
	defer rows.Close()

	var verses []models.Verse
	for rows.Next() {
		var verse models.Verse
		if err := rows.Scan(&verse.Text, &verse.Number); err != nil {
			log.Errorf("Error scanning updated verse for music ID %s: %v", updatedMusic.ID, err)
			return models.Music{}, err
		}
		verses = append(verses, verse)
	}

	updatedMusic.Verses = verses
	log.Infof("Music with ID %s and verses updated successfully", updatedMusic.ID)

	return updatedMusic, nil
}

func NewMusicRepository(pool *pgxpool.Pool) models.MusicRepository {
	log.Info("Creating new music repository")
	return &musicRepository{pool: pool}
}
