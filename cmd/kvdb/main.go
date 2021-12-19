package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bilalislam/kvdb/pkg/kvdb"
	"github.com/bilalislam/kvdb/pkg/kvdb/aol"
)

const (
	defaultPort          = 8080
	defaultBasePath      = "/tmp"
	defaultMaxRecordSize = 64 * 1024
)

var (
	httpPort      int
	basePath      string
	maxRecordSize int
	async         bool
)

func init() {
	flag.IntVar(&httpPort, "port", defaultPort, "http server listening port")
	flag.StringVar(&basePath, "path", defaultBasePath, "storage path")
	flag.IntVar(&maxRecordSize, "maxRecordSize", defaultMaxRecordSize, "max size of a database record")
	flag.BoolVar(&async, "async", async, "file sync for durability")
}
func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)
	async := true // interval ile olmalı

	db, err := aol.NewStore(aol.Config{
		BasePath:      basePath,
		MaxRecordSize: &maxRecordSize,
		Async:         &async,
		Logger:        logger,
	})

	if err != nil {
		logger.Fatal(err)
	}

	server := startHTTPServer(httpPort, logger, db)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	stopHTTPServer(server, logger)
	if err := db.Close(); err != nil {
		logger.Printf("Could not close database, %s", err)
	}
}

func startHTTPServer(port int, logger *log.Logger, db kvdb.Store) *http.Server {
	srv := http.Server{Addr: fmt.Sprintf(":%d", port)}
	logger.Printf(fmt.Sprintf("Started server on port: %d", port))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[1:]
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Key cannot be empty"))
			return
		}

		switch r.Method {
		case http.MethodGet:
			val, err := db.Get(key)
			if err != nil {
				handleError(w, r, logger, err)
				return
			}

			w.Write(val)

		case http.MethodPut, http.MethodPost:
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				handleError(w, r, logger, err)
				return
			}

			err = db.Set(key, data)
			if err != nil {
				handleError(w, r, logger, err)
				return
			} else {
				db.Sync()
			}

			w.WriteHeader(http.StatusCreated)

		default:
			w.WriteHeader(http.StatusNotImplemented)
		}

	})

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("Could not start server: %s", err)
		}
	}()

	return &srv
}

func stopHTTPServer(server *http.Server, logger *log.Logger) {
	logger.Print("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not close server: %s", err)
	}
}

func handleError(w http.ResponseWriter, r *http.Request, logger *log.Logger, err error) {
	if kvdb.IsNotFoundError(err) {
		w.WriteHeader(http.StatusNotFound)
	} else if kvdb.IsBadRequestError(err) {
		logger.Printf("%s %s failed: %s", r.Method, r.URL.Path, err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		logger.Printf("%s %s failed: %s", r.Method, r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
