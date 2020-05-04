package persistent

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/hxhieu/oceanbridge-wrapper-go-common/models"
	"github.com/hxhieu/oceanbridge-wrapper-go-common/persistent"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type firebaseUser struct {
	email    string
	password string
}

// FirestoreUserStore that uses Firestore as the backend
type FirestoreUserStore struct {
	store   *firestore.Client
	context *context.Context
}

func (s FirestoreUserStore) getSessions(email string) (*[]models.AuthSession, error) {
	// Login OK, check existing sessions
	query := s.store.Collection("sessions").Where("Email", "==", email).Documents(*s.context)
	var sessions []models.AuthSession
	for {
		doc, err := query.Next()
		if err == iterator.Done {
			break
		}
		var session models.AuthSession
		if err := doc.DataTo(&session); err == nil {
			sessions = append(sessions, session)
		}
	}

	return &sessions, nil
}

// NewUserStoreFirestore creates an instance of user store backed by Firestore
func NewUserStoreFirestore() (persistent.UserStore, error) {
	// Create a Firebase instance
	serviceAccount := os.Getenv("FIREBASE_SERVICE_ACCOUNT")
	if len(serviceAccount) <= 0 {
		return nil, status.Error(codes.Aborted, "Firestore configure is missing")
	}
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON([]byte(serviceAccount)))
	if err != nil {
		return nil, err
	}
	// Get the Firestore
	store, err := app.Firestore(context.Background())
	if err != nil {
		return nil, err
	}
	return &FirestoreUserStore{store: store, context: &ctx}, nil
}

// Login and gets sessions linked to the user
func (s FirestoreUserStore) Login(email string, password string) (*[]models.AuthSession, error) {
	doc, err := s.store.Collection("users").Doc(email).Get(*s.context)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Hide the detail error from Firebase
			return nil, status.Error(codes.NotFound, "User not found.")
		}
		return nil, err
	}

	var user persistent.User
	if err = doc.DataTo(&user); err != nil || (err == nil && user.Password != password) {
		return nil, status.Error(codes.Unauthenticated, "Invalid email and password combination.")
	}

	// Login OK, check existing sessions
	sessions, err := s.getSessions(email)
	if err != nil {
		return nil, status.Error(codes.Aborted, "An error has occured while querying the data.")
	}

	return sessions, nil
}

// Logout and clear the sessions
func (s FirestoreUserStore) Logout(sessionIDs ...string) error {
	query := s.store.Collection("sessions").Where("ID", "in", sessionIDs).Documents(*s.context)
	for {
		doc, err := query.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		doc.Ref.Delete(*s.context)
	}
	return nil
}

// SetToken update a sesion with a new token
func (s FirestoreUserStore) SetToken(sessionID string, accessToken string) error {
	if _, err := s.store.Collection("sessions").Doc(sessionID).Update(*s.context, []firestore.Update{{
		Path:  "AccessToken",
		Value: accessToken,
	}}); err != nil {
		return status.Error(codes.Aborted, "Invalid session.")
	}
	return nil
}

// GetSessions of a user from the store
func (s FirestoreUserStore) GetSessions(email string) (*[]models.AuthSession, *persistent.User, error) {
	// Get session
	sessions, err := s.getSessions(email)
	if err != nil {
		return nil, nil, err
	}

	// Get associated user
	userDoc, err := s.store.Collection("users").Doc(email).Get(*s.context)
	if err != nil {
		return nil, nil, status.Error(codes.Aborted, "Invalid session.")
	}
	var user persistent.User
	if err := userDoc.DataTo(&user); err != nil {
		return nil, nil, err
	}

	return sessions, &user, nil
}

// CreateSession create a new session for the user
func (s FirestoreUserStore) CreateSession(email string, accessToken string) error {
	ref := s.store.Collection("sessions").NewDoc()
	session := models.AuthSession{
		Email:       email,
		AccessToken: accessToken,
		ID:          ref.ID,
	}
	_, err := ref.Set(*s.context, session)
	if err != nil {
		return status.Error(codes.Aborted, "Fail to create login session.")
	}
	return nil
}
