package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zakkbob/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id != 1 {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:        1,
		CreatedAt: time.Now(),
		Year:      2006,
		Title:     "Borat! Cultural Learnings of America for Make Benefit Glorious Nation of Kazakhstan",
		Runtime:   84,
		Genres: []string{
			"comedy",
			"mockumentary",
			"satire",
		},
		Version: 1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

