package sqlstore

import (
	"database/sql"
	"errors"
	"time"

	"github.com/rautaruukkipalich/go_notes/internal/model"
)

func (r *NoteRepo) GetNoteById(id int) (model.Note, error) {
	var note model.Note

	stmt := `
	select
		id, author_id, body, is_public, created_at, updated_at 
	from notes
	where id = $1`

	if err := r.sqlstore.db.QueryRow(stmt, id).Scan(
		&note.ID, &note.AuthorID, &note.Body,
		&note.IsPublic, &note.CreatedAt, &note.UpdatedAt,
	); err != nil {
		return note, err
	}

	return note, nil
}

func (r *NoteRepo) GetNotes(userId int, filter_body string, filter_author int, limit int, offset int) ([]model.Note, error) {
	var stmt string
	var rows *sql.Rows
	var err error

	if filter_author != 0 {
		stmt = `select
			id, author_id, body, is_public, created_at, updated_at 
		from notes
		where 
			(is_public = true or author_id = $1)
			and (lower(body) like concat('%', $2::text, '%'))
			and (author_id = $3)
		order by is_public, id asc
		limit $4
		offset $5			
		`
		rows, err = r.sqlstore.db.Query(stmt, userId, filter_body, filter_author, limit, offset)
	} else {
		stmt = `
		select
			id, author_id, body, is_public, created_at, updated_at 
		from notes
		where 
			(is_public = true or author_id = $1)
			and (lower(body) like concat('%', $2::text, '%'))
		order by is_public, id asc
		limit $3
		offset $4			
		`
		rows, err = r.sqlstore.db.Query(stmt, userId, filter_body, limit, offset)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return parseRowsToSliceNotes(rows)
}

func (r *NoteRepo) HeatCache() ([]model.Note, error) {
	stmt := `
		select id, author_id, body, is_public, created_at, updated_at 
		from notes
		where is_public = true
		order by id asc	
		`

	rows, err := r.sqlstore.db.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return parseRowsToSliceNotes(rows)
}

func (r *NoteRepo) Set(n *model.Note) (*model.Note, error) {
	stmt := `insert	into notes 
			(author_id, body, is_public, created_at, updated_at) 
			values ($1, $2, $3, $4, $5)
			returning id`
	now := time.Now().UTC()
	n.CreatedAt = now
	n.UpdatedAt = now
	err := r.sqlstore.db.QueryRow(stmt,
		n.AuthorID, n.Body,	n.IsPublic, n.CreatedAt, n.UpdatedAt, 
	).Scan(&n.ID)
	return n, err
}

func (r *NoteRepo) Patch(n *model.Note) error {
	stmt := `update notes
			set body = $1, is_public=$2, updated_at=$3 
			where id = $4`
	now := time.Now().UTC()
	res, err := r.sqlstore.db.Exec(stmt,
		n.Body,	n.IsPublic, now, n.ID, 
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if count == 0 {
		return errors.New("no notes updated")
	}
	return err
}

func (r *NoteRepo) Delete(id int) error {
	stmt := `delete from notes where id = $1`
	res, err := r.sqlstore.db.Exec(stmt, id) 
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if count == 0 {
		return errors.New("no notes deleted")
	}
	return err
}








// func (r *NoteRepo) SetNotes() error {
// 	for i:=0; i<20; i++ {
// 		note := &model.Note{
// 			AuthorID:  rand.Intn(10),
// 			Body:      strings.Repeat("a", i),
// 			IsPublic:  []bool{true, false}[rand.Intn(2)],
// 		}
// 		err := r.Set(note)
// 		if err != nil {
// 			return err
// 		}	
// 	}
// 	return nil
// }