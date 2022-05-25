package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bank/domain"
	"github.com/bank/server"
	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

var Dbpool *pgxpool.Pool

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

// getAlbums responds with the list of all albums as JSON.
func GetAccounts(c *gin.Context) {

	account := domain.Account{}
	var accounts []domain.Account

	accounts, err := account.Read(Dbpool, 0)
	if err == nil {
		//log.WithFields(log.Fields{"accounts": accounts}).Debug("Returned accounts")
		c.IndentedJSON(http.StatusOK, accounts)
	} else {
		c.IndentedJSON(http.StatusNotFound, err.Error())
	}
}

// getAlbums responds with the list of all albums as JSON.
func GetAccountById(c *gin.Context) {
	id := c.Param("id")

	account := domain.Account{}

	account, err := account.ReadById(Dbpool, id)
	if err == nil {
		//log.WithFields(log.Fields{"accounts": accounts}).Debug("Returned accounts")
		c.IndentedJSON(http.StatusOK, account)
	} else {
		c.IndentedJSON(http.StatusNotFound, err.Error())
	}
}

// postAlbums adds an album from JSON received in the request body.
func PostAccount(c *gin.Context) {
	var newAccount domain.Account

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAccount); err != nil {
		return
	}

	if newAccount.Id == 0 {
		log.Println("new account id is empty")
		/*
			var id_s = albums[len(albums)-1].ID
			var id int
			id, err := strconv.Atoi(id_s)
			if err == nil {
				newAlbum.ID = strconv.Itoa(id + 1)
			} else {
				c.IndentedJSON(http.StatusBadRequest, newAlbum)
				return
			}
		*/
	}
	// Add the new account to the database.
	newAccount.Write(Dbpool)
	//	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAccount)
}

// get pool information
func GetPool(c *gin.Context) {
	stat := Dbpool.Stat()

	var status domain.DbpoolStat
	status.AcquireConns = stat.AcquiredConns()
	status.AcquireCount = stat.AcquireCount()
	status.AcquireDuration = stat.AcquireDuration()
	status.ConstructingConns = stat.ConstructingConns()
	status.EmptyAcquireCount = stat.EmptyAcquireCount()
	status.IdleConns = stat.IdleConns()
	status.MaxConns = stat.MaxConns()
	status.TotalConns = stat.TotalConns()
	c.IndentedJSON(http.StatusOK, status)

}

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}
func serve() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	router.GET("/albums", server.GetAlbums)
	router.GET("/albums/:id", server.GetAlbumByID)
	router.POST("/albums", server.PostAlbums)
	router.GET("/accounts", GetAccounts)
	router.GET("/accounts/:id", GetAccountById)
	router.POST("/accounts", PostAccount)
	router.GET("/pool", GetPool)
	router.Use(JSONMiddleware())

	//router.Run("localhost:8080")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

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

	Dbpool, err = pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	log.WithFields(log.Fields{"pool": Dbpool}).
		Debug("Open database")
	CheckError(err)

	err = Dbpool.Ping(context.Background())
	CheckError(err)

	// always close database at program exit
	defer func() {
		log.Debug("Close database")
		Dbpool.Close()
	}()

	log.Debug("Connected to database!")

	//test.DoTransactionTest(dbpool)
	//test.Server()

	//server.Serve()
	serve()

	log.Debug("Closed  Bank")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
