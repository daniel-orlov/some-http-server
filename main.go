package main

import (
	"github.com/myntra/golimit/store"
	"net/http"
	"some-http-server/internal/cache"
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
	pflag.String("svc_name", "quotes", "Service name")

	pflag.Bool("auto_migrate_enable", false, "Enable automatic db migration")
	pflag.Int("db_migrate_to", 1, "Version of migration to be used for auto-migration")

	pflag.Bool("cache_enable", false, "Enable cache")

	pflag.Int32("limit_threshold", 10, "Max allowed create quote requests per user per window")
	pflag.Int32("limit_window", 60, "Rate limiter window in seconds")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
}

func main() {
	// Establishing connection to all databases.
	db := database.NewDBFromEnv()
	defer db.Close()

	if viper.GetBool("auto_migrate_enable") {
		if err := database.Migrate(db.DB, viper.GetInt("db_migrate_to")); err != nil {
			logrus.Fatal(err)
		}
	}

	// Initiating all repositories.
	quoteRepo := repository.NewQuoteRepo(db)

	// Wrapping them into cash.
	cachedQR := getCachedQuoteRepo(quoteRepo)

	// Connecting clients.
	xSvcClientStub := repository.NewExternalSvcClientStub(repository.XPConfig{})

	// Instantiating main service.
	svc := service.NewQuoteService(xSvcClientStub, cachedQR)

	// Creating a router, a limit store and registering handlers on it.
	r := mux.NewRouter()
	s := store.NewStore()
	defer s.Close()
	transport.NewHandler(svc, s).Register(r)

	srv := &http.Server{
		Handler:           r,
		Addr:              viper.GetString("address"),
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	logrus.Infof("Starting server %s", viper.GetString("address"))
	logrus.Fatal(srv.ListenAndServe())
}

func getCachedQuoteRepo(repo *repository.QuoteRepository) service.QuoteRepo {
	cacheStorage := cache.NewCacheStore(cache.Ristretto, cache.GetOptsFromEnv())
	return repository.NewProxyQuoteRepo(repo, cacheStorage)
}
