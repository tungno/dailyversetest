// internal/repositories/firestore_friend_repository.go
package repositories

import (
	"context"
	"proh2052-group6/pkg/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FirestoreFriendRepository struct {
	Client *firestore.Client
}

func NewFirestoreFriendRepository(client *firestore.Client) FriendRepository {
	return &FirestoreFriendRepository{Client: client}
}

func (fr *FirestoreFriendRepository) CreateFriendRequest(ctx context.Context, friend *models.Friend) error {
	docID := friend.Email + "_" + friend.FriendEmail
	_, err := fr.Client.Collection("friends").Doc(docID).Set(ctx, friend)
	return err
}

func (fr *FirestoreFriendRepository) GetFriendRequest(ctx context.Context, senderEmail, recipientEmail string) (*models.Friend, error) {
	docID := senderEmail + "_" + recipientEmail
	doc, err := fr.Client.Collection("friends").Doc(docID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var friend models.Friend
	if err := doc.DataTo(&friend); err != nil {
		return nil, err
	}
	return &friend, nil
}

func (fr *FirestoreFriendRepository) UpdateFriendRequest(ctx context.Context, senderEmail, recipientEmail string, updates map[string]interface{}) error {
	docID := senderEmail + "_" + recipientEmail
	_, err := fr.Client.Collection("friends").Doc(docID).Set(ctx, updates, firestore.MergeAll)
	return err
}

func (fr *FirestoreFriendRepository) DeleteFriendRequest(ctx context.Context, senderEmail, recipientEmail string) error {
	docID := senderEmail + "_" + recipientEmail
	_, err := fr.Client.Collection("friends").Doc(docID).Delete(ctx)
	return err
}

func (fr *FirestoreFriendRepository) GetFriends(ctx context.Context, userEmail string) ([]models.Friend, error) {
	var friends []models.Friend

	// Query for friends where the user is the sender
	iter := fr.Client.Collection("friends").Where("Email", "==", userEmail).Where("Status", "==", "accepted").Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var friend models.Friend
		if err := doc.DataTo(&friend); err != nil {
			continue
		}
		friends = append(friends, friend)
	}

	// Query for friends where the user is the recipient
	iter = fr.Client.Collection("friends").Where("FriendEmail", "==", userEmail).Where("Status", "==", "accepted").Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var friend models.Friend
		if err := doc.DataTo(&friend); err != nil {
			continue
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

func (fr *FirestoreFriendRepository) GetPendingFriendRequests(ctx context.Context, userEmail string) ([]models.Friend, error) {
	var pendingRequests []models.Friend
	iter := fr.Client.Collection("friends").Where("FriendEmail", "==", userEmail).Where("Status", "==", "pending").Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var friend models.Friend
		if err := doc.DataTo(&friend); err != nil {
			continue
		}
		pendingRequests = append(pendingRequests, friend)
	}
	return pendingRequests, nil
}
