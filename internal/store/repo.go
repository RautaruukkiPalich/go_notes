package store

type NoteRepo interface {
	NoteCreater
	NoteGetter
	NoteSetter
	NoteDeletter
}

type NoteGetter interface {
	Get()
}

type NoteCreater interface {
	Create()
}

type NoteSetter interface {
	Set()
}

type NoteDeletter interface {
	Delete()
}