package main

import (
	"context"
	"errors"
	ihttp "github.com/rainu/r-ray/internal/http"
	"github.com/rainu/r-ray/internal/http/controller"
	"github.com/rainu/r-ray/internal/processor"
	"github.com/rainu/r-ray/internal/store"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := readConfig()
	if err != nil {
		logrus.WithError(err).Error("Error while reading config.")
		os.Exit(1)
		return
	}

	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	userStore := store.NewUser()
	for _, credential := range cfg.RequestCredentials {
		userStore.Add(credential.UsernameAndPassword())
	}

	p := processor.New(userStore)
	cors := ihttp.CorsMiddleware{
		Delegate: controller.NewProxy(cfg.RequestHeaderPrefix, cfg.ForwardRequestHeader, cfg.ForwardResponseHeader, p),

		Origins: cfg.CorsAllowOrigin,
		Methods: cfg.CorsAllowMethods,
		Headers: cfg.CorsAllowHeaders,
		MaxAge:  cfg.CorsAllowMaxAge,
	}
	server := ihttp.NewServer(cfg.BindingAddr, cors)

	errChan := make(chan error)
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT)

	select {
	case err := <-errChan:
		logrus.WithError(err).Error("Error while listen and serve.")
		os.Exit(2)
		return
	case <-signalChan:
		logrus.Info("Shutting down...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logrus.WithError(err).Error("Error while shutdown server.")
	}
}
