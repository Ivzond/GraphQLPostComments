package main

import (
	"GraphQLPostComments/api"
	"GraphQLPostComments/api/generated"
	"GraphQLPostComments/internal/storage"
	"GraphQLPostComments/internal/storage/memory"
	"GraphQLPostComments/internal/storage/postgres"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

type Config struct {
	StorageType         string `yaml:"storage_type" default:"memory"`
	PostgresDatabaseURL string `yaml:"postgres_database_url"`
}

func loadConfig(configPath string) (*Config, error) {
	config := &Config{}
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yml"
	}

	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	var store storage.Storage
	switch config.StorageType {
	case "postgres":
		pgStore, err := postgres.NewStorage(config.PostgresDatabaseURL)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		store = pgStore
	case "memory":
		store = memory.NewStorage()
	default:
		log.Fatalf("Unsupported storage type: %v", config.StorageType)
	}

	resolver := api.NewResolver(store)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.Handle("/subscriptions", handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: resolver}),
	))
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
