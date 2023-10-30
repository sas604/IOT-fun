package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) listnAndServe() error {

	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", app.config.port),
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	shutDownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutDownError <- err
		}

		app.logger.Info("completing bg tasks", "add", srv.Addr)
		// app.wg.Wait()

		shutDownError <- nil
	}()

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutDownError
	if err != nil {
		return err
	}

	app.logger.Info("stoped server", "addr", srv.Addr)
	return nil

}
