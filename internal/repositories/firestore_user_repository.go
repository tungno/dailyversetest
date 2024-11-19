// internal/repositories/firestore_user_repository.go
package repositories

import (
	"context"
	"fmt"
	"proh2052-group6/pkg/models"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FirestoreUserRepository struct {
	Client *firestore.Client
}

func NewFirestoreUserRepository(client *firestore.Client) UserRepository {
	return &FirestoreUserRepository{Client: client}
}

func (ur *FirestoreUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	doc, err := ur.Client.Collection("users").Doc(email).Get(ctx)
	if err != nil {
		return nil, err
	}
	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *FirestoreUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	iter := ur.Client.Collection("users").Where("UsernameLower", "==", strings.ToLower(username)).Limit(1).Documents(ctx)
	defer iter.Stop()
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *FirestoreUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := ur.Client.Collection("users").Doc(user.Email).Set(ctx, user)
	return err
}

func (ur *FirestoreUserRepository) UpdateUser(ctx context.Context, email string, updates map[string]interface{}) error {
	_, err := ur.Client.Collection("users").Doc(email).Set(ctx, updates, firestore.MergeAll)
	return err
}

func (ur *FirestoreUserRepository) SearchUsersByUsername(ctx context.Context, query string) ([]*models.User, error) {
	iter := ur.Client.Collection("users").
		Where("UsernameLower", ">=", strings.ToLower(query)).
		Where("UsernameLower", "<=", strings.ToLower(query)+"\uf8ff").
		Documents(ctx)
	defer iter.Stop()

	var users []*models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			continue
		}
		users = append(users, &user)
	}

	return users, nil
}
