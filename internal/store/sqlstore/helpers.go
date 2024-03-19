package sqlstore

import (
	"database/sql"

	"github.com/rautaruukkipalich/go_notes/internal/model"
)

func parseRowsToSliceNotes(rows *sql.Rows) ([]model.Note, error) {
	notes := []model.Note{}

	for rows.Next() {
		var note model.Note
		if err := rows.Scan(
			&note.ID, &note.AuthorID, &note.Body,
			&note.IsPublic, &note.CreatedAt, &note.UpdatedAt,
		); err != nil {
			return notes, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}
