package booking

import (
	"context"
	"errors"
	"sync"
)

type ConcurrentSafeStore struct {
	bookings map[string]Booking
	sync.RWMutex
}

func NewConcurrentSafeStore() *ConcurrentSafeStore {
	return &ConcurrentSafeStore{
		bookings: make(map[string]Booking),
	}
}

func (s *ConcurrentSafeStore) Book(b Booking) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyTaken
	}
	s.bookings[b.SeatID] = b
	return nil
}

func (s *ConcurrentSafeStore) Hold(b Booking) (Booking, error) {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.bookings[b.SeatID]; exists {
		return Booking{}, ErrSeatAlreadyTaken
	}
	s.bookings[b.SeatID] = b
	return b, nil
}

func (s *ConcurrentSafeStore) ListBookings(movieID string) ([]Booking, error) {
	s.RLock()
	defer s.RUnlock()
	var result []Booking
	for _, booking := range s.bookings {
		if booking.MovieID == movieID {
			result = append(result, booking)
		}
	}
	return result, nil
}

func (s *ConcurrentSafeStore) Confirm(_ context.Context, _ string, _ string) (Booking, error) {
	return Booking{}, errors.New("not implemented")
}

func (s *ConcurrentSafeStore) Release(_ context.Context, _ string, _ string) error {
	return errors.New("not implemented")
}
