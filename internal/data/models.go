package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies    MovieModel
	Actor     ActorModel
	Directors DirectorModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:    MovieModel{DB: db},
		Actor:     ActorModel{DB: db},
		Directors: DirectorModel{DB: db},
	}

}
