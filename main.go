package main

import (
	"bandersnatch/api/handlers"
	"bandersnatch/pkg/entities"
	"bandersnatch/pkg/game"
	"bandersnatch/pkg/player"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Bandersnatch: A Dynamically Randomized State Automaton (a.k.a DYRASTAT)")
	nexus := &game.Nexus{}
	if err := nexus.LoadFromFile("sample.json"); err != nil {
		fmt.Println(err)
	}

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	conn, err := pq.ParseURL(os.Getenv("DB_URI"))
	if err != nil {
		panic(err)
		return
	}

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		panic(err)
		return
	}

	db = db.Debug()
	db.AutoMigrate(&entities.Player{})

	playerRepo := player.NewPostgresRepo(db)

	playerSvc := player.NewService(playerRepo)
	gameSvc := game.NewService(nexus)

	r := httprouter.New()
	handlers.MakePlayerHandlers(r, playerSvc)
	handlers.MakeGameHandlers(r, playerSvc, gameSvc)


	port := os.Getenv("PORT")
	if port == "" {
		port = "1729"
	}

	fmt.Println("Bandersnatch server up and running ...")
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		panic(err)
		return
	}


}
