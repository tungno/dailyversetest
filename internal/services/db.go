// internal/services/db.go
package services

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

type DatabaseInterface interface {
	Collection(path string) *firestore.CollectionRef
	Close() error
}

type FirestoreDB struct {
	Client *firestore.Client
}

func NewFirestoreClient(ctx context.Context) (DatabaseInterface, error) {
	//sa := option.WithCredentialsFile("./serviceAccountKey.json")
	client, err := firestore.NewClient(ctx, "prog2052-project")
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Firestore successfully.")
	return &FirestoreDB{Client: client}, nil
}

func (db *FirestoreDB) Collection(path string) *firestore.CollectionRef {
	return db.Client.Collection(path)
}

func (db *FirestoreDB) Close() error {
	return db.Client.Close()
}
