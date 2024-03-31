package mockstore

import "github.com/rautaruukkipalich/go_notes/internal/model"

var Test_notes = []model.Note{
	{
		AuthorID: 1,
		Body:     "test note 1",
		IsPublic: true,
	},
	{
		AuthorID: 1,
		Body:     "test note 2",
		IsPublic: true,
	},
	{
		AuthorID: 1,
		Body:     "test note 4",
		IsPublic: true,
	},
	{
		AuthorID: 2,
		Body:     "test note 2",
		IsPublic: true,
	},
	{
		AuthorID: 2,
		Body:     "test note 4",
		IsPublic: true,
	},
}
