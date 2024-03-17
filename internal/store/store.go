package store

type Store interface {
	Note() NoteRepo
}

type Cache interface {
	Note() NoteCache
	User() UserCache
}
