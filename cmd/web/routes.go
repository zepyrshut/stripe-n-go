package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// routers from url to methods
func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/virtual-terminal", app.VirtualTerminal)
	mux.Post("/payment-succeeded", app.PaymentSucceeded)

	mux.Get("/charge-once", app.ChargeOnce)

	// load all static elements
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
