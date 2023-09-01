package mocks

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"github.com/nshimiyimanaamani/paypack-backend/core/feedback"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

var _ (feedback.Repository) = (*repositoryMock)(nil)

type repositoryMock struct {
	mu       sync.Mutex
	counter  uint64
	messages map[string]feedback.Message
}

// NewRepository ...
func NewRepository() feedback.Repository {
	return &repositoryMock{
		messages: make(map[string]feedback.Message),
	}
}

func (repo *repositoryMock) Save(ctx context.Context, msg *feedback.Message) (*feedback.Message, error) {
	const op errors.Op = "app/feedback/mocks/repositoryMock.Save"

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

func (repo *repositoryMock) Retrieve(ctx context.Context, id string) (feedback.Message, error) {
	const op errors.Op = "app/feedback/mocks/repositoryMock.Retrieve"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	val, ok := repo.messages[id]
	if !ok {
		return feedback.Message{}, errors.E(op, "message not found", errors.KindNotFound)
	}

	return val, nil
}

func (repo *repositoryMock) Update(ctx context.Context, msg feedback.Message) error {
	const op errors.Op = "app/feedback/mocks/repositoryMock.Update"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.messages[msg.ID]; !ok {
		return errors.E(op, "message not found", errors.KindNotFound)
	}

	repo.messages[msg.ID] = msg

	return nil
}

func (repo *repositoryMock) Delete(ctx context.Context, id string) error {
	const op errors.Op = "app/feedback/mocks/repositoryMock.Delete"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, ok := repo.messages[id]
	if !ok {
		return errors.E(op, "message not found", errors.KindNotFound)
	}

	delete(repo.messages, id)
	return nil
}

func (repo *repositoryMock) RetrieveAll(ctx context.Context, offset, limit uint64) (feedback.MessagePage, error) {
	const op errors.Op = "app/feedback/mocks/repository.List"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	items := make([]feedback.Message, 0)

	if offset < 0 || limit <= 0 {
		return feedback.MessagePage{}, nil
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	for _, v := range repo.messages {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if id >= first && id < last {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := feedback.MessagePage{
		Messages: items,
		PageMetadata: feedback.PageMetadata{
			Total:  repo.counter,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}
