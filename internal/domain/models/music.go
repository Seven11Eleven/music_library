package models

import (
	"context"
	"time"
)

type Verse struct {
	Text   string `json:"text"`
	Number int    `json:"number"`
}

type Music struct {
	ID          string     `json:"id,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	Verses      []Verse    `json:"verses,omitempty"`
	Link        string     `json:"link,omitempty"`
	SongName    string     `json:"song_name"`
	GroupName   string     `json:"group_name"`
}

type MusicQuery struct {
	GroupName string `json:"group_name"`
	SongName  string `json:"song_name"`
}

type MusicFilters struct {
	ReleaseDate *time.Time
	Link        *string
	SongName    *string
	GroupName   *string
}

type MusicRepository interface {
	SaveMusic(ctx context.Context, music *Music) (*Music, error)
	GetMusic(ctx context.Context, musicName, groupName string) (*Music, error)
	GetMusicsByFilters(ctx context.Context, filters MusicFilters, page, pageSize int) ([]Music, error)
	GetMusicTextWithPaginationByVerse(ctx context.Context, musicID string, limit, offset int) (*Music, error)
	DeleteMusic(ctx context.Context, musicID string) error
	UpdateMusic(ctx context.Context, music Music) (Music, error)
}

type MusicService interface {
	SaveMusic(ctx context.Context, music *MusicQuery) (*Music, error)
	GetMusicsByFilters(ctx context.Context, filters MusicFilters, page, pageSize int) ([]Music, error)
	GetMusicTextWithPaginationByVerse(ctx context.Context, musicID string, limit, offset int) (*Music, error)
	DeleteMusic(ctx context.Context, musicID string) error
	UpdateMusic(ctx context.Context, music Music) (Music, error)
}
