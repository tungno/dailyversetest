// internal/services/db.go
package services

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

func NewFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, "prog2052-project")
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Firestore successfully.")
	return client, nil
}
