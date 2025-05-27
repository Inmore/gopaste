package sqlite

import (
	"database/sql"
	"time"

	"github.com/inmore/gopaste/internal/model"
	"github.com/inmore/gopaste/internal/storage"
)

var _ storage.Storage = (*Store)(nil)

type Store struct{ db *sql.DB }

func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	return s, s.migrate()
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`Create table if not exists pastes(
		id text primary key,
		content text,
		expires_at integer
	);`)
	return err
}

func (s *Store) Save(p *model.Paste) error {
	_, err := s.db.Exec(
		`insert or replace into pastes(id, content, expires_at) values(?,?,?);`,
		p.ID, p.Content, p.ExpiresAt.Unix(),
	)
	return err
}

func (s *Store) Load(id string) (*model.Paste, error) {
	row := s.db.QueryRow(`select content, expires_at from pastes where id=?;`, id)

	var content string
	var expires int64
	if err := row.Scan(&content, &expires); err != nil {
		return nil, err
	}
	if time.Now().Unix() > expires {
		return nil, sql.ErrNoRows
	}
	return &model.Paste{
		ID:        id,
		Content:   content,
		ExpiresAt: time.Unix(expires, 0),
	}, nil
}

func (s *Store) DeleteExpired() (int, error) {
	res, err := s.db.Exec(`delete from pastes where expires_at < ?;`, time.Now().Unix())
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}

func (s *Store) Close() error { return s.db.Close() }
