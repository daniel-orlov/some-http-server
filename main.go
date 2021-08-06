package main

import (
	"net/http"
	"some-http-server/internal/database"
	"some-http-server/internal/repository"
	"some-http-server/internal/service"
	"some-http-server/internal/transport"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.String("address", "0.0.0.0:8888", "Address from which to serve")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
}

func main() {
	// Establishing connection to database.
	db := database.NewDBFromEnv()
	defer db.Close()

	// Initiating all repositories.
	quoteRepo := repository.NewQuoteRepo(db)

	// Connecting clients
	xSvcClientStub := repository.NewExClientStub(repository.XPConfig{})

	// Instantiating main service.
	svc := service.NewQuoteService(xSvcClientStub, quoteRepo)

	r := mux.NewRouter()

	transport.NewTaskHandler(svc).Register(r)

	address := viper.GetString("address")

	srv := &http.Server{
		Handler:           r,
		Addr:              address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	logrus.Infof("Starting server %s", address)
	logrus.Fatal(srv.ListenAndServe())
}
