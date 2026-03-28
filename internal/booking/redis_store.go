package booking

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const defaultHoldTTL = 1 * time.Minute

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func (s *RedisStore) Book(b Booking) error {
	session, err := s.Hold(b)
	if err != nil {
		return err
	}
	log.Printf("Session held: %+v", session)

	return nil
}

func (s *RedisStore) ListBookings(movieID string) ([]Booking, error) {
	pattern := fmt.Sprintf("seat:%s:*", movieID)
	var sessions []Booking // Placeholder for actual booking retrieval logic
	ctx := context.Background()

	iter := s.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		val, err := s.rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			return nil, err
		}
		session, err := parseSession(val)
		if err != nil {
			continue
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func parseSession(val string) (Booking, error) {
	var session Booking
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return Booking{}, err
	}
	return Booking{
		ID:      session.ID,
		MovieID: session.MovieID,
		UserID:  session.UserID,
		SeatID:  session.SeatID,
		Status:  session.Status,
	}, nil
}

func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func rdbKey(movieID, seatID string) string {
	return fmt.Sprintf("seat:%s:%s", movieID, seatID)
}

func (s *RedisStore) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {
	seatKey, err := s.rdb.Get(ctx, sessionKey(sessionID)).Result()
	if err == redis.Nil {
		return Booking{}, errors.New("session not found")
	}
	if err != nil {
		return Booking{}, err
	}

	val, err := s.rdb.Get(ctx, seatKey).Result()
	if err != nil {
		return Booking{}, err
	}

	b, err := parseSession(val)
	if err != nil {
		return Booking{}, err
	}

	if b.UserID != userID {
		return Booking{}, errors.New("unauthorized")
	}

	b.Status = 1
	data, _ := json.Marshal(b)
	s.rdb.Set(ctx, seatKey, data, 0)

	return b, nil
}

func (s *RedisStore) Release(ctx context.Context, sessionID string, userID string) error {
	seatKey, err := s.rdb.Get(ctx, sessionKey(sessionID)).Result()
	if err == redis.Nil {
		return errors.New("session not found")
	}
	if err != nil {
		return err
	}

	s.rdb.Del(ctx, sessionKey(sessionID), seatKey)
	return nil
}

func (s *RedisStore) Hold(b Booking) (Booking, error) {
	id := uuid.New().String()
	now := time.Now()
	ctx := context.Background()
	key := rdbKey(b.MovieID, b.SeatID)
	data := Booking{
		ID:        id,
		MovieID:   b.MovieID,
		UserID:    b.UserID,
		SeatID:    b.SeatID,
		Status:    0,
		ExpiresAt: now.Add(defaultHoldTTL),
	}
	val, _ := json.Marshal(data)

	res := s.rdb.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX",
		TTL:  defaultHoldTTL,
	})

	ok := res.Val() == "OK"
	if !ok {
		return Booking{}, ErrSeatAlreadyTaken
	}

	s.rdb.Set(ctx, sessionKey(id), key, defaultHoldTTL)

	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		UserID:    b.UserID,
		SeatID:    b.SeatID,
		Status:    0,
		ExpiresAt: now.Add(defaultHoldTTL),
	}, nil
}
