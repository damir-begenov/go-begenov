package main

import (
	"errors"
	"fmt"
	"github.com/shynggys9219/greenlight/internal/data"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

func (app *application) createActorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Fullname   string   `json:"fullname"`
		Year       int32    `json:"year"`
		Films      []string `json:"films"`
		Girlfriend string   `json:"girlfriend"`
	}
	err := app.readJSON(w, r, &input) //non-nil pointer as the target decode destination
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}
	actor := &data.Actor{
		Fullname:   input.Fullname,
		Year:       input.Year,
		Films:      input.Films,
		Girlfriend: input.Girlfriend,
	}
	err = app.models.Actor.INSERTACTOR(actor)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", actor.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"actor": actor}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) createDirectorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string   `json:"name"`
		Surname string   `json:"surname"`
		Awords  []string `json:"awords"`
	}
	err := app.readJSON(w, r, &input) //non-nil pointer as the target decode destination
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}
	directors := &data.Directors{
		Name:    input.Name,
		Surname: input.Surname,
		Awords:  input.Awords,
	}
	err = app.models.Directors.InsertDirector(directors)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/directors/%d", directors.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"director": directors}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) showActorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	actor, err := app.models.Actor.GetActors(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"actor": actor}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}
	err := app.readJSON(w, r, &input) //non-nil pointer as the target decode destination
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Call the Get() method to fetch the data for a specific movie. We also need to // use the errors.Is() function to check if it returns a data.ErrRecordNotFound // error, in which case we send a 404 Not Found response to the client.
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//}
//func (app *application) showMovieHandlerTitle(w http.ResponseWriter, r *http.Request) {
//	title, err := app.readIDParam(r)
//	if err != nil {
//		app.notFoundResponse(w, r)
//		return
//	}
//	// Call the Get() method to fetch the data for a specific movie. We also need to // use the errors.Is() function to check if it returns a data.ErrRecordNotFound // error, in which case we send a 404 Not Found response to the client.
//	movie, err := app.models.Movies.GetByTitle(string(title))
//	if err != nil {
//		switch {
//		case errors.Is(err, data.ErrRecordNotFound):
//			app.notFoundResponse(w, r)
//		default:
//			app.serverErrorResponse(w, r, err)
//		}
//		return
//	}
//	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//	}
//}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) deleteActorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Actor.DeleteActor(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "actor successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) updateActorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	actor, err := app.models.Actor.GetActors(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Fullname   string   `json:"fullname"`
		Year       int32    `json:"year"`
		Films      []string `json:"films"`
		Girlfriend string   `json:"girlfriend"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	actor.Fullname = input.Fullname
	actor.Year = input.Year
	actor.Films = input.Films
	actor.Girlfriend = input.Girlfriend

	err = app.models.Actor.UpdateActor(actor)
	fmt.Println("INPUT", actor)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"actor": actor}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	err = app.models.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
func (app *application) updateMovieeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
		return
	}
	// Use pointers for the Title, Year and Runtime fields.
	var input struct {
		Title   *string  `json:"title"`   // This will be nil if there is no corresponding key in the JSON.
		Year    *int32   `json:"year"`    // Likewise...
		Runtime *int32   `json:"runtime"` // Likewise...
		Genres  []string `json:"genres"`  // We don't need to change this because slices already have the zero-value nil.
	}
	// Decode the JSON as normal.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres // Note that we don't need to dereference a slice.
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	// Call the GetAll() method to retrieve the movies, passing in the various filter // parameters.
	movies, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSON response containing the movie data.
	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listDirectorsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string
		Surname string
		Awords  []string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Name = app.readString(qs, "name", "")
	input.Surname = app.readString(qs, "surname", "")
	input.Awords = app.readCSV(qs, "awords", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "surname", "-id", "-name", "-awords", "awords", "-surname", "-runtime"}

	// Call the GetAll() method to retrieve the movies, passing in the various filter // parameters.
	directors, err := app.models.Directors.GetAllDirectors(input.Name, input.Surname, input.Awords, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSON response containing the movie data.
	err = app.writeJSON(w, http.StatusOK, envelope{"director": directors}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
