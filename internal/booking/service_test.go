package booking

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestConcurrentBooking_ExactkyOneWins(t *testing.T) {
	store := NewRedisStore(redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	}))

	fmt.Println("Connected to Redis, starting test...")
	fmt.Printf("Connected to Redis at %s\n", store.rdb)

	svc := NewService(store)

	const numGoroutines = 100_000 // 100k users trying to book the same seat

	var (
		successes atomic.Int64
		failures  atomic.Int64
		wg        sync.WaitGroup
	)

	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func(userNum int) {
			defer wg.Done()
			err := svc.Book(Booking{
				MovieID: "movie-123",
				UserID:  uuid.New().String(),
				SeatID:  "seat-1",
			})
			if err == nil {
				successes.Add(1)
			} else if err == ErrSeatAlreadyTaken {
				failures.Add(1)
			} else {
				t.Errorf("unexpected error: %v", err)
			}
		}(i)
	}

	wg.Wait()

	if successes.Load() != 1 {
		t.Errorf("expected exactly 1 success, got %d", successes.Load())
	}
	if failures.Load() != numGoroutines-1 {
		t.Errorf("expected %d failures, got %d", numGoroutines-1, failures.Load())
	}
}
