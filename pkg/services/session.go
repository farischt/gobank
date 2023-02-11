package services

import (
	"fmt"
	"time"

	"github.com/farischt/gobank/pkg/store"
	"github.com/farischt/gobank/pkg/types"
	"golang.org/x/crypto/bcrypt"
)

type SessionService interface {
	Get(tokenId string) (*types.SerializedSessionToken, error)
	comparePassword(hashedPassword string, password []byte) bool
	Create(accountId uint, password string) (*types.SerializedSessionToken, error)
	IsValidSessionToken(tokenId string) (*types.SerializedSessionToken, bool)
	Delete(tokenId string) error
}

type sessionService struct {
	store store.Store
}

func NewSessionService(store store.Store) SessionService {
	return &sessionService{
		store: store,
	}
}

func (s *sessionService) Get(tokenId string) (*types.SerializedSessionToken, error) {
	t, err := s.store.SessionToken.GetSessionToken(tokenId)

	if err != nil {
		return nil, err
	}

	return t.Serialize(), nil
}

func (s *sessionService) comparePassword(hashedPassword string, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), password)
	return err == nil
}

func (s *sessionService) Create(accountId uint, password string) (*types.SerializedSessionToken, error) {

	if accountId <= 0 {
		return nil, fmt.Errorf("missing_account_number")
	}

	// Check if the account exists
	a, err := s.store.Account.GetAccount(accountId)
	if err != nil {
		return nil, fmt.Errorf("invalid_id")
	}

	// Compare the password
	if !s.comparePassword(a.Password, []byte(password)) {
		return nil, fmt.Errorf("invalid_password")
	}

	// Create a new session token
	token, err := s.store.SessionToken.CreateSessionToken(a.ID)

	return token.Serialize(), err
}

func (s *sessionService) IsValidSessionToken(tokenId string) (*types.SerializedSessionToken, bool) {
	st, err := s.Get(tokenId)
	if err != nil {
		return nil, false
	}

	elapsed := time.Since(st.CreatedAt)
	return st, elapsed <= time.Second*1000
}

func (s *sessionService) Delete(tokenId string) error {
	return s.store.SessionToken.DeleteSessionToken(tokenId)
}
