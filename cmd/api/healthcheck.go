package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"movie": env}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered an error a problem and could not process your request", http.StatusInternalServerError)
	}
}
