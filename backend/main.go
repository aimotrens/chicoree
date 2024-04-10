package main

import (
	noroute "chicoree/controllers/no_route"
	"chicoree/ent"
	"chicoree/ent/migrate"
	"context"
	"embed"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	//go:embed static/*
	staticFiles embed.FS
)

func main() {
	if stat, err := os.Stat(".env"); err == nil && !stat.IsDir() {
		err := godotenv.Load()
		if err != nil {
			panic("Error loading .env file")
		}
	}

	dbType := os.Getenv("DATABASE_TYPE")
	dbUri := os.Getenv("DATABASE_URI")
	{
		client, err := ent.Open(dbType, dbUri)
		if err != nil {
			log.Fatalf("failed opening connection to sqlite: %v", err)
		}

		// Run the auto migration tool.
		if err := client.Schema.Create(context.Background(),
			migrate.WithDropColumn(true),
			migrate.WithDropIndex(true),
			migrate.WithGlobalUniqueID(true),
		); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}

		client.Close()
	}

	g := gin.Default()

	noroute.
		NewController(staticFiles).
		RegisterRoutes(g)

	g.Run(":8080")
}
