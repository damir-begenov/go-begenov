package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// By default, the keys in the JSON object are equal to the field names in the struct ( ID,
// CreatedAt, Title and so on).
type Movie struct {
	ID        int64     `json:"id"`                       // Unique integer ID for the movie
	CreatedAt time.Time `json:"-"`                        // Timestamp for when the movie is added to our database, "-" directive, hidden in response
	Title     string    `json:"title"`                    // Movie title
	Year      int32     `json:"year,omitempty"`           // Movie release year, "omitempty" - hide from response if empty
	Runtime   int32     `json:"runtime,omitempty,string"` // Movie runtime (in minutes), "string" - convert int to string
	Genres    []string  `json:"genres,omitempty"`         // Slice of genres for the movie (romance, comedy, etc.)
	Version   int32     `json:"version"`                  // The version number starts at 1 and will be incremented each
	// time the movie information is updated
}
type Actor struct {
	ID         int64     `json:"id"`             // Unique integer ID for the movie
	CreatedAt  time.Time `json:"-"`              // Timestamp for when the movie is added to our database, "-" directive, hidden in response
	Fullname   string    `json:"fullname"`       // Movie title
	Year       int32     `json:"year,omitempty"` // Movie release year, "omitempty" - hide from response if empty
	Films      []string  `json:"films,omitempty"`
	Girlfriend string    `json:"girlfriend"` // Movie title

}

type Directors struct {
	ID      int64    `json:"id"`      // Unique integer ID for the movie
	Name    string   `json:"name"`    // Movie title
	Surname string   `json:"surname"` // Movie title
	Awords  []string `json:"awords,omitempty"`
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type MovieModel struct {
	DB *sql.DB
}
type ActorModel struct {
	DB *sql.DB
}
type DirectorModel struct {
	DB *sql.DB
}

func (m DirectorModel) InsertDirector(directors *Directors) error {
	query := `
		INSERT INTO directors(name, surname,awords)
		VALUES ($1, $2, $3)
		RETURNING id`
	return m.DB.QueryRow(query, &directors.Name, &directors.Surname, pq.Array(&directors.Awords)).Scan(&directors.ID)
}

func (m ActorModel) INSERTACTOR(actor *Actor) error {
	query := `
		INSERT INTO actor(fullname, year,girlfriend, films)
		VALUES ($1, $2, $3,$4)
		RETURNING id, created_at`

	return m.DB.QueryRow(query, &actor.Fullname, &actor.Year, &actor.Girlfriend, pq.Array(&actor.Films)).Scan(&actor.ID, &actor.CreatedAt)
}

// method for inserting a new record in the movies table.
func (m MovieModel) Insert(movie *Movie) error {
	query := `
		INSERT INTO movies(title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`
	return m.DB.QueryRow(query, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres)).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT id, created_at, title, year, runtime, genres, version FROM movies
WHERE id = $1`
	var movie Movie
	err := m.DB.QueryRow(query, id).Scan(&movie.ID,
		&movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &movie, nil
}
func (m ActorModel) GetActors(id int64) (*Actor, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT id, created_at, fullname, year, films, girlfriend FROM actor
WHERE id = $1`
	var actor Actor
	err := m.DB.QueryRow(query, id).Scan(&actor.ID,
		&actor.CreatedAt, &actor.Fullname, &actor.Year, pq.Array(&actor.Films), &actor.Girlfriend,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &actor, nil
}
func (m MovieModel) GetByTitle(title string) (*Movie, error) {
	if title < "null" {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT id, created_at, title, year, runtime, genres, version FROM movies
WHERE title = $1`
	var movie Movie
	err := m.DB.QueryRow(query, title).Scan(&movie.ID,
		&movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &movie, nil
}

func (m ActorModel) UpdateActor(actor *Actor) error {
	// Declare the SQL query for updating the record and returning the new version     // number.
	query := `
	UPDATE actor
	SET fullname = $1, year = $2, films = $3, girlfriend= $4
	WHERE id = $5
	RETURNING girlfriend;`
	args := []any{
		actor.Fullname,
		actor.Year,
		pq.Array(actor.Films),
		actor.Girlfriend,
		actor.ID,
	}
	// Use the QueryRow() method to execute the query, passing in the args slice as a     // variadic parameter and scanning the new version value into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&actor.Girlfriend)
}
func (m MovieModel) Update(movie *Movie) error {
	query := `
UPDATE movies
SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1 WHERE id = $5 AND version = $6
RETURNING version`
	args := []any{movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version, // Add the expected movie version.
	}
	// Execute the SQL query. If no matching row could be found, we know the movie // version has changed (or the record has been deleted) and we return our custom // ErrEditConflict error.
	err := m.DB.QueryRow(query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
func (m ActorModel) DeleteActor(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM actor WHERE id = $1`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows     // affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we     // return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// method for deleting a specific record from the movies table.
func (m MovieModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM movies WHERE id = $1`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows     // affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we     // return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

//	func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) {
//		// Construct the SQL query to retrieve all movie records.
//		query := `
//
// SELECT id, created_at, title, year, runtime, genres, version FROM movies
// ORDER BY id`
//
//		// Create a context with a 3-second timeout.
//		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//		defer cancel()
//		// Use QueryContext() to execute the query. This returns a sql.Rows resultset // containing the result.
//		rows, err := m.DB.QueryContext(ctx, query)
//		if err != nil {
//			return nil, err
//		}
//		// Importantly, defer a call to rows.Close() to ensure that the resultset is closed // before GetAll() returns.
//		defer rows.Close()
//		// Initialize an empty slice to hold the movie data.
//		movies := []*Movie{}
//		// Use rows.Next to iterate through the rows in the resultset.
//		for rows.Next() {
//			// Initialize an empty Movie struct to hold the data for an individual movie.
//			var movie Movie
//			// Scan the values from the row into the Movie struct. Again, note that we're // using the pq.Array() adapter on the genres field here.
//			err := rows.Scan(
//				&movie.ID,
//				&movie.CreatedAt,
//				&movie.Title,
//				&movie.Year,
//				&movie.Runtime,
//				pq.Array(&movie.Genres),
//				&movie.Version,
//			)
//			if err != nil {
//				return nil, err
//			}
//
//			// Add the Movie struct to the slice.
//			movies = append(movies, &movie)
//		}
//		// When the rows.Next() loop has finished, call rows.Err() to retrieve any error // that was encountered during the iteration.
//		if err = rows.Err(); err != nil {
//			return nil, err
//		}
//		// If everything went OK, then return the slice of movies.
//		return movies, nil
//	}
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) { // Update the SQL query to include the filter conditions.
	query := fmt.Sprintf(`
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') AND (genres @> $2 OR $2 = '{}')
ORDER BY %s %s, id ASC
LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// As our SQL query now has quite a few placeholder parameters, let's collect the // values for the placeholders in a slice. Notice here how we call the limit() and // offset() methods on the Filters struct to get the appropriate values for the
	// LIMIT and OFFSET clauses.
	args := []any{title, pq.Array(genres), filters.limit(), filters.offset()}
	// And then pass the args slice to QueryContext() as a variadic parameter.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	movies := []*Movie{}
	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.ID,
			&movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}

func (m DirectorModel) GetAllDirectors(name string, surname string, awords []string, filters Filters) ([]*Directors, error) { // Update the SQL query to include the filter conditions.
	query := fmt.Sprintf(`
SELECT id, name, surname, awords
FROM directors
WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '') AND (awords @> $2 OR $2 = '{}')
ORDER BY %s %s, id ASC
LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// As our SQL query now has quite a few placeholder parameters, let's collect the // values for the placeholders in a slice. Notice here how we call the limit() and // offset() methods on the Filters struct to get the appropriate values for the
	// LIMIT and OFFSET clauses.
	args := []any{name, pq.Array(awords), filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	director := []*Directors{}
	for rows.Next() {
		var directors Directors
		err := rows.Scan(&directors.ID,
			&directors.Name, &directors.Surname, pq.Array(&directors.Awords),
		)
		if err != nil {
			return nil, err
		}
		director = append(director, &directors)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return director, nil
}
