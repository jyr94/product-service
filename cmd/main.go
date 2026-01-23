package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/jyr94/product-service/internal/adapters/cache"
	httpadapter "github.com/jyr94/product-service/internal/adapters/httphandler"
	"github.com/jyr94/product-service/internal/adapters/persistence"
	"github.com/jyr94/product-service/internal/application"
	"github.com/jyr94/product-service/internal/config"
)

func main() {

	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	authMiddleware := httpadapter.BasicAuthMiddleware(
		cfg.HTTP.BasicAuthUser,
		cfg.HTTP.BasicAuthPass,
	)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	var cacheAdapter *cache.RedisCache

	if cfg.Redis.Enabled {
		rdb := redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})

		if err := rdb.Ping(context.Background()).Err(); err != nil {
			log.Fatalf("failed to connect redis: %v", err)
		}

		cacheAdapter = cache.NewRedisCache(rdb)
		log.Println("redis connected")
	}

	productRepo := persistence.NewProductRepository(db)

	productService := application.NewProductService(
		productRepo,
		cacheAdapter,
	)

	productHandler := httpadapter.NewProductHandler(productService)

	router := mux.NewRouter()
	router.Use(authMiddleware)
	productHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("HTTP server running on :%s", cfg.HTTP.Port)
	log.Fatal(server.ListenAndServe())
}
