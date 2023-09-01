package mocks

import (
	"context"
	"sort"
	"strconv"
	"sync"

	"github.com/nshimiyimanaamani/paypack-backend/core/notifs"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

type mockSmsRepository struct {
	mu            sync.Mutex
	counter       uint64
	notifications map[string]notifs.Notification
}

// NewRepository ...
func NewRepository() notifs.Repository {
	return &mockSmsRepository{
		counter:       0,
		notifications: make(map[string]notifs.Notification),
	}
}

func (repo *mockSmsRepository) Save(ctx context.Context, sms notifs.Notification) (notifs.Notification, error) {
	const op errors.Op = "core/notifs/mocks/mocksSmsRepository.Save"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	var empty notifs.Notification

	for _, n := range repo.notifications {
		if n.ID == sms.ID {
			return empty, errors.E(op, "conflicting message ids", errors.KindAlreadyExists)
		}
	}

	repo.counter++
	sms.ID = strconv.FormatUint(repo.counter, 10)
	repo.notifications[sms.ID] = sms

	return sms, nil
}

func (repo *mockSmsRepository) Find(ctx context.Context, id string) (notifs.Notification, error) {
	const op errors.Op = "core/notifs/mocks/mocksSmsRepository.Find"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	var empty notifs.Notification

	val, ok := repo.notifications[id]
	if !ok {
		return empty, errors.E(op, "sms not found", errors.KindNotFound)
	}
	return val, nil
}

func (repo *mockSmsRepository) List(ctx context.Context, nspace string, offset, limit uint64) (notifs.NoticationPage, error) {
	const op errors.Op = "core/notifs/mocks/mocksSmsRepository.List"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	items := make([]notifs.Notification, 0)

	if offset < 0 || limit <= 0 {
		return notifs.NoticationPage{}, nil
	}

	first := uint64(offset) + 1
	last := first + uint64(limit)

	for _, v := range repo.notifications {
		id, _ := strconv.ParseUint(v.ID, 10, 64)
		if v.Sender == nspace && id >= first && id < last {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := notifs.NoticationPage{
		Notifications: items,
		PageMetadata: notifs.PageMetadata{
			Total:  repo.counter,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}

func (repo *mockSmsRepository) Count(ctx context.Context, nspace string) (uint64, error) {
	const op errors.Op = "core/notifs/mocks/mocksSmsRepository.Count"

	repo.mu.Lock()
	defer repo.mu.Unlock()

	var count uint64

	for _, v := range repo.notifications {
		if v.Sender == nspace {
			count++
		}
	}
	return count, nil
}
