package mocks

import (
	"context"
	"strconv"
	"sync"

	"github.com/rugwirobaker/paypack-backend/app/feedback"
)

var _ (feedback.Repository) = (*messageRepoMock)(nil)

type messageRepoMock struct {
	mu       sync.Mutex
	counter  uint64
	messages map[string]feedback.Message
}

// NewREPO ...
func NewREPO() feedback.Repository {
	return &messageRepoMock{
		messages: make(map[string]feedback.Message),
	}
}

func (repo *messageRepoMock) Save(ctx context.Context, msg *feedback.Message) (*feedback.Message, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, en := range repo.messages {
		if en.ID == msg.ID {
			return nil, feedback.ErrConflict
		}
	}
	repo.counter++
	msg.ID = strconv.FormatUint(repo.counter, 10)
	repo.messages[msg.ID] = *msg

	return msg, nil
}

func (repo *messageRepoMock) Retrieve(ctx context.Context, id string) (feedback.Message, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	val, ok := repo.messages[id]
	if !ok {
		return feedback.Message{}, feedback.ErrNotFound
	}

	return val, nil
}

func (repo *messageRepoMock) Update(ctx context.Context, msg feedback.Message) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.messages[msg.ID]; !ok {
		return feedback.ErrNotFound
	}

	repo.messages[msg.ID] = msg

	return nil
}

func (repo *messageRepoMock) Delete(ctx context.Context, id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, ok := repo.messages[id]
	if !ok {
		return feedback.ErrNotFound
	}

	delete(repo.messages, id)
	return nil
}
