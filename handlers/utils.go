package handlers

import (
	"database/sql"
	"strings"

	db "ctrlstock/db"
)

// Server agrupa las dependencias que necesitan los handlers.
type Server struct {
	Queries *db.Queries
}

func toNullString(s string) sql.NullString {
	s = strings.TrimSpace(s)
	return sql.NullString{String: s, Valid: s != ""}
}

func toNullInt32(n *int32) sql.NullInt32 {
	if n == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *n, Valid: true}
}
