package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/lucasthedev/rssagg/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	feed, err := urlToFeed("http://wagslane.dev/index.xml")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(feed)

	fmt.Println("inciando sistema")

	godotenv.Load(".env")

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("Porta não foi encontrada no ambiente")
	}

	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		log.Fatal("DB URL não foi encontrada no ambiente")
	}

	conn, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("Não foi possível conectar no banco de dados", err)
	}

	dbConnection := database.New(conn)
	apiCfg := apiConfig{
		DB: dbConnection,
	}

	go startScraping(dbConnection, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	routerPaths := chi.NewRouter()
	routerPaths.Get("/healthz", handlerReadiness)
	routerPaths.Get("/error", handlerError)
	routerPaths.Post("/createUser", apiCfg.handlerCreateUser)
	routerPaths.Get("/getUser", apiCfg.middleware(apiCfg.handlerGetUserByApiKey))
	routerPaths.Post("/createFeed", apiCfg.middleware(apiCfg.handlerCreateFeed))
	routerPaths.Get("/getFeeds", apiCfg.handlerGetFeeds)
	routerPaths.Post("/feedFollows", apiCfg.middleware(apiCfg.handlerCreateFeedFollow))
	routerPaths.Get("/feedFollows", apiCfg.middleware(apiCfg.handlerGetFeedFollows))
	routerPaths.Delete("/feedFollows/{feedFollowID}", apiCfg.middleware(apiCfg.handlerUnfollowFeed))

	router.Mount("/v1", routerPaths)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Servidor iniciou na porta %v", portString)
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
