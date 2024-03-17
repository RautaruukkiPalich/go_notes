package store

type Store interface {
	Note() NoteRepo
}