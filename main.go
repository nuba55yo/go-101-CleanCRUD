package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/nuba55yo/go-101-CleanCRUD/docs/v1"
	_ "github.com/nuba55yo/go-101-CleanCRUD/docs/v2"

	"github.com/nuba55yo/go-101-CleanCRUD/application/usecase"
	"github.com/nuba55yo/go-101-CleanCRUD/infrastructure/logging"
	gormp "github.com/nuba55yo/go-101-CleanCRUD/infrastructure/persistence/gorm"
	httpx "github.com/nuba55yo/go-101-CleanCRUD/presentation/http/router"
)

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

func main() {
	_ = godotenv.Load()

	// DB
	db, err := gormp.Open()
	if err != nil {
		log.Fatal(err)
	}
	if err := gormp.AutoMigrateTables(db); err != nil {
		log.Fatal(err)
	}
	if err := gormp.EnsureIndexes(db); err != nil {
		log.Fatal(err)
	}

	// Logger (use case)
	appLogger, flush, err := logging.NewZapLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = flush() }()

	// DI: Repository -> UseCase -> Router
	bookRepository := gormp.NewBookRepositoryGorm(db)
	bookUseCase := usecase.NewBookUseCase(bookRepository, systemClock{}, appLogger)
	router := httpx.NewRouter(bookUseCase) // ??? /api/v1, /api/v2, /docs, /swagger

	// Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	_ = router.Run(":" + port)
}
