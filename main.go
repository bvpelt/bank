package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bank/server"
	"github.com/bank/util"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

//var Dbpool *pgxpool.Pool

func init() {
	//
	// Logging levels are: Trace, Debug, Info, Warning, Error, Fatal and Panic see https://github.com/Sirupsen/logrus
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		PadLevelText:    true,
		TimestampFormat: time.RFC3339,
	})
	/*
		customFormatter := new(logrus.TextFormatter)
		log.SetFormatter(customFormatter)
		customFormatter.TimestampFormat = "2016-03-28 15:04:05.000"
		customFormatter.FullTimestamp = true
	*/
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)
	log.SetLevel(log.DebugLevel)
}

func serve() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := server.StartServer()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func main() {
	log.Debug("Started Bank")
	// export DATABASE_URL=postgres://testuser:12345@localhost:5432/bank

	var err error
	var Dbpool *pgxpool.Pool

	Dbpool, err = util.Dbaccess()

	if err != nil {
		return
	}

	// always close database at program exit
	defer func() {
		log.Debug("Setup close of already opened database")
		Dbpool.Close()
	}()

	//test.DoTransactionTest(dbpool)
	//test.Server()

	//server.Serve()
	serve()

	log.Debug("Closed  Bank")
}
