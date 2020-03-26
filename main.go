package main

import (
	"bandersnatch/api/handlers"
	"bandersnatch/pkg/entities"
	"bandersnatch/pkg/game"
	"bandersnatch/pkg/player"
	"bandersnatch/utils"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})
	log.SetOutput(os.Stdout)
	log.Printf("Running on %s", os.Getenv("ENV"))
	if os.Getenv("ENV") != "PROD" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func connectToDb() *gorm.DB {
	conn, err := pq.ParseURL(os.Getenv("DB_URI"))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	if os.Getenv("DEBUG") == "true" {
		db = db.Debug()
	}

	db.AutoMigrate(&entities.Player{})
	return db
}

func initNegroni() *negroni.Negroni {
	n := negroni.New()
	n.Use(negronilogrus.NewCustomMiddleware(log.DebugLevel, &log.JSONFormatter{PrettyPrint: true}, "API requests"))
	n.Use(negroni.NewRecovery())
	return n
}

func main() {
	utils.PrintAsciiArt()

	nexus := &game.Nexus{}
	if err := nexus.LoadFromFile(os.Getenv("NEXUS_FILE")); err != nil {
		log.Fatal(err)
		return
	}

	n := initNegroni()
	r := httprouter.New()
	n.UseHandler(r)
	db := connectToDb()

	playerRepo := player.NewPostgresRepo(db)

	playerSvc := player.NewService(playerRepo)
	gameSvc := game.NewService(nexus)

	handlers.MakePlayerHandlers(r, playerSvc)
	handlers.MakeGameHandlers(r, playerSvc, gameSvc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1729"
	}

	r.HandlerFunc("POST", "/api/bandersnatch/health", func(w http.ResponseWriter, r *http.Request) {
		utils.RespWrap(w, http.StatusOK, "Good")
	})

	log.WithField("event", "START").Info("Listening on port " + port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatal(err)
		return
	}

}
