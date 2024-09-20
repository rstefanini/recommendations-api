package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

var (
	host       = os.Getenv("DB_HOST")
	port       = os.Getenv("DB_PORT")
	username   = os.Getenv("DB_USERNAME")
	pass       = os.Getenv("DB_ROOT_PASSWORD")
	database   = os.Getenv("DB_NAME")
	collection = os.Getenv("DB_COLLECTION")
)

func New() *Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", username, pass, host, port)))

	if err != nil {
		log.Fatal(err)

	}

	//TODO refactor Builder
	return &Service{
		client:     client,
		db:         client.Database(database),
		collection: client.Database(database).Collection(collection),
	}

}

func (s *Service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *Service) Client() *mongo.Client {
	return s.client
}

func (s *Service) Database() *mongo.Database {
	return s.db
}

func (s *Service) Collection() *mongo.Collection {
	return s.collection
}
