package main

import (
	"log"
	"net/http"
	"os"

	internalhttp "github.com/BimaPDev/SignalStack/api/internal/http"
	"github.com/BimaPDev/SignalStack/api/internal/repo"
	"github.com/joho/godotenv"
)

// func main
// - load config from environment via config.Load()
// - open database connection via repo.Open()
// - create structured logger via observability.NewLogger()
// - create HTTP router via http.NewRouter(db, logger)
// - start http.Server on configured port
// - listen for SIGINT/SIGTERM signals
// - on signal: graceful shutdown with timeout context
func main() {
	godotenv.Load()
	db, err := repo.Open(os.Getenv("POSTGRES_ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	router := internalhttp.HandlerRouter()
	http.ListenAndServe(":3000", router)

}