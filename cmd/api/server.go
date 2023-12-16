package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		quit := make(chan os.Signal, 1)

		// Listen for incoming SIGINT and SIGTERM signals and relay them to
		// the quit channel.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Read the signal from the quit channel. This code will block until
		// a signal is received.
		s := <-quit

		app.logger.PrintInfo("caugth signal", map[string]string{
			"signal": s.String(),
		})
		os.Exit(0)
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	return srv.ListenAndServe()
}
