/**
 *  Provides a utility function to initialize a Firestore client for database operations.
 *
 *  @file       db.go
 *  @package    services
 *
 *  @functions
 *  - NewFirestoreClient(ctx) - Creates and returns a new Firestore client for the specified context.
 *
 *  @dependencies
 *  - "cloud.google.com/go/firestore": Provides Firestore client capabilities.
 *  - Google Cloud Project: The project must be configured and accessible for Firestore operations.
 *
 *  @behaviors
 *  - Establishes a connection to the Firestore database using the provided context.
 *  - Logs a success message upon successful connection.
 *  - Returns an error if the client initialization fails.
 *
 *  @example
 *  ```
 *  ctx := context.Background()
 *  client, err := NewFirestoreClient(ctx)
 *  if err != nil {
 *      log.Fatalf("Failed to connect to Firestore: %v", err)
 *  }
 *  defer client.Close()
 *  ```
 *
 *  @errors
 *  - Returns an error if the Firestore client cannot be created.
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package services

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

// NewFirestoreClient creates and returns a new Firestore client.
// It takes a context as an argument, which is used to manage the lifecycle of the client connection.
func NewFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, "prog2052-project") // Replace "prog2052-project" with your actual Google Cloud Project ID.
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Firestore successfully.") // Log successful connection.
	return client, nil
}
