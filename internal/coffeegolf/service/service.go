package service

import (
	"database/sql"

	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
)

const defaultTouramentLength = 10

type Service struct {
	db      *sql.DB
	queries *database.Queries
}

func New(db *sql.DB, queries *database.Queries) *Service {
	return &Service{
		db:      db,
		queries: queries,
	}
}
