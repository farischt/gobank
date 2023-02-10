package services

import (
	"fmt"
	"time"

	"github.com/farischt/gobank/store"
	"github.com/farischt/gobank/types"
)

type SessionService interface {
	Get(tokenId string) (*types.SerializedSessionToken, error)
	Create(accountId uint) (*types.SerializedSessionToken, error)
	IsValidSessionToken(tokenId string) (*types.SerializedSessionToken, bool)
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

func (s *sessionService) Create(accountId uint) (*types.SerializedSessionToken, error) {

	if accountId <= 0 {
		return nil, fmt.Errorf("missing_account_number")
	}

	// Check if the account exists
	a, err := s.store.Account.GetAccount(accountId)
	if err != nil {
		return nil, err
	}

	// TODO Check if the password is correct
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
