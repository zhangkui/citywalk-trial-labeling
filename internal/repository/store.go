package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store { return &Store{db: db} }

func (s *Store) DB() *sql.DB { return s.db }

func (s *Store) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *Store) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *Store) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *Store) Begin(ctx context.Context) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, nil)
}

type UserRecord struct {
	ID           int64
	Username     string
	PasswordHash string
	Nickname     sql.NullString
	Avatar       sql.NullString
	Bio          sql.NullString
	City         sql.NullString
	Status       int8
	LastLoginAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Roles        []string
	Permissions  []string
}

type RouteRecord struct {
	ID            int64
	Title         string
	Description   sql.NullString
	ThemeID       sql.NullInt64
	City          sql.NullString
	StartLat      sql.NullFloat64
	StartLng      sql.NullFloat64
	EndLat        sql.NullFloat64
	EndLng        sql.NullFloat64
	TotalDistance sql.NullInt64
	Duration      sql.NullInt64
	Difficulty    sql.NullInt64
	CoverImage    sql.NullString
	CreatedBy     int64
	Status        int8
	FavoriteCount int64
	Rating        float64
	RatingCount   int64
	ViewCount     int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     sql.NullTime
}

type StoryRecord struct {
	ID           int64
	Title        string
	Content      string
	CoverImage   sql.NullString
	LandmarkIDs  sql.NullString
	RouteID      sql.NullInt64
	CreatedBy    int64
	Status       int8
	LikeCount    int64
	ViewCount    int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}

type LandmarkRecord struct {
	ID           int64
	Name         string
	Address      sql.NullString
	Lat          float64
	Lng          float64
	CategoryID   sql.NullInt64
	Description  sql.NullString
	CoverImage   sql.NullString
	CreatedBy    int64
	Status       int8
	Rating       float64
	RatingCount  int64
	ViewCount    int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}

type CommentRecord struct {
	ID         int64
	TargetType string
	TargetID   int64
	Content    string
	ParentID   int64
	CreatedBy  int64
	LikeCount  int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  sql.NullTime
}

type EventRecord struct {
	ID                  int64
	Title               string
	RouteID             sql.NullInt64
	MeetupAddress       sql.NullString
	MeetupLat           sql.NullFloat64
	MeetupLng           sql.NullFloat64
	MeetupTime          time.Time
	StartTime           time.Time
	EndTime             time.Time
	MaxParticipants     int64
	CurrentParticipants  int64
	Fee                 int64
	Description         sql.NullString
	CreatedBy           int64
	Status              int8
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type ParticipationRecord struct {
	ID        int64
	EventID   int64
	UserID    int64
	Name      string
	Remark    sql.NullString
	Status    int8
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BadgeRecord struct {
	ID             int64
	Name           string
	Description    sql.NullString
	Icon           sql.NullString
	ConditionType  string
	ConditionValue  sql.NullInt64
	CreatedAt      time.Time
}

func boolToInt(v bool) int8 {
	if v {
		return 1
	}
	return 0
}

func int8ToBool(v int8) bool { return v != 0 }

func idsCSV(ids []int64) string {
	if len(ids) == 0 {
		return ""
	}
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("%d", id))
	}
	return strings.Join(parts, ",")
}

