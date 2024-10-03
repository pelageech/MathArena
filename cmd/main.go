package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/pelageech/matharena/internal/data"
	"github.com/pelageech/matharena/internal/handlers"
	"github.com/pelageech/matharena/internal/pkg/ioutil"
	"github.com/pelageech/matharena/internal/postgres"
)

const (
	bindAddress     = ":8080"
	shutdownTimeout = 30 * time.Second
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 120 * time.Second
)

func main() {
	l := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
	})

	// Set up a routerDeck
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(
		cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		}))

	// Set up a context
	ctx := context.Background()

	// Get the connection string from the environment
	connStr := os.Getenv("DB_CONN_STR")
	if connStr == "" {
		l.Fatal("DB_CONN_STR env var is not set")
	}

	// wait until database is up
	time.Sleep(5 * time.Second)

	// Set up a database connection

	psqlDB, err := postgres.NewPSQLDatabase(ctx, connStr, l)
	if err != nil {
		log.Fatal(err)
	}

	// set up context so that we can ping the database and don't wait forever
	cancelCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// check whether connection is established
	err = psqlDB.Ping(cancelCtx)
	if err != nil {
		log.Fatal("Unable to ping database in main", "error", err)
	}

	// get salt length from env
	saltLengthString := os.Getenv("SALT_LENGTH")
	if saltLengthString == "" {
		l.Fatal("SALT_LENGTH env var is not set")
	}

	// convert to int
	saltLength, err := strconv.Atoi(saltLengthString)
	if err != nil {
		l.Fatal("Unable to convert salt length to int", "error", err)
	}

	// get token expiration time from env
	tokenExpirationTime := os.Getenv("TOKEN_EXPIRATION_TIME")
	if tokenExpirationTime == "" {
		l.Fatal("TOKEN_EXPIRATION_TIME env var is not set")
	}

	// convert to int
	tokenExpirationTimeInt, err := strconv.Atoi(tokenExpirationTime)
	if err != nil {
		l.Fatal("Unable to convert token expiration time to int", "error", err)
	}

	// convert to time.Duration
	duration := time.Duration(tokenExpirationTimeInt) * time.Hour

	// get sign key from env
	tokenSignKey := os.Getenv("TOKEN_SECRET")
	if tokenSignKey == "" {
		l.Fatal("TOKEN_SECRET env var is not set")
	}

	// Set up a datalayer
	dl := data.New(psqlDB, saltLength, duration, []byte(tokenSignKey))

	// Set up error writer
	ew := ioutil.JSONErrorWriter{Logger: l}

	authHandlers := handlers.NewAuthorization(dl, ew, l)

	// Set up routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/signup", authHandlers.SignUp)
		r.Options("/signup", authHandlers.SignUp)
		r.Post("/signin", authHandlers.SignIn)
		r.Get("/user/{id}", authHandlers.GetUserInfo)
	})

	// create a new server
	s := http.Server{
		Addr:         bindAddress, // configure the bind address
		Handler:      r,           // set the default handler
		ErrorLog:     l.StandardLog(),
		ReadTimeout:  readTimeout,  // max time to read request from the client
		WriteTimeout: writeTimeout, // max time to write response to the client
		IdleTimeout:  idleTimeout,  // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Info("Starting server", "port", bindAddress)

		l.Fatal("Error form server", "error", s.ListenAndServe())
	}()

	// trap interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Got signal: %v", sig)
	l.Infof("Shutting down...")

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	cancelCtx, cancel = context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = s.Shutdown(cancelCtx)
	if err != nil {
		l.Fatal("Error shutting down server", "error", err)
	}
}
