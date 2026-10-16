package handlers

import (
	"database/sql"

	db "ctrlstock/db"
)

// Server agrupa las dependencias que necesitan los handlers.
type Server struct {
	Queries *db.Queries
}

func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func toNullInt32(n *int32) sql.NullInt32 {
	if n == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *n, Valid: true}
}
