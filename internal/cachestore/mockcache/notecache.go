package mockcache

type NoteCache struct{
	client map[string]any
}

func newNoteCache() *NoteCache {
	return &NoteCache{
		client: make(map[string]any),
	}
}

func (n *NoteCache) Get()    {}
func (n *NoteCache) Set()    {}
func (n *NoteCache) Delete() {}