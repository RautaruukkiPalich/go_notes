package cachestore

type Cache interface {
	Note() NoteCache
	User() UserCache
}

type NoteCache interface {
	Getter
	Setter
	Deletter
}

type UserCache interface {
	Getter
	Setter
	Deletter
}

type Getter interface {
	Get()
}

type Setter interface {
	Set()
}

type Deletter interface {
	Delete()
}