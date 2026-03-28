package booking

import (
	"context"
	"errors"
)

type MemoryStore struct {
	// map seatsId -> booking
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: make(map[string]Booking),
	}
}

func (s *MemoryStore) Book(b Booking) error {
	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyTaken
	}
	s.bookings[b.SeatID] = b
	return nil
}

func (s *MemoryStore) Hold(b Booking) (Booking, error) {
	if _, exists := s.bookings[b.SeatID]; exists {
		return Booking{}, ErrSeatAlreadyTaken
	}
	s.bookings[b.SeatID] = b
	return b, nil
}

func (s *MemoryStore) ListBookings(movieID string) ([]Booking, error) {
	var result []Booking
	for _, booking := range s.bookings {
		if booking.MovieID == movieID {
			result = append(result, booking)
		}
	}
	return result, nil
}

func (s *MemoryStore) Confirm(_ context.Context, _ string, _ string) (Booking, error) {
	return Booking{}, errors.New("not implemented")
}

func (s *MemoryStore) Release(_ context.Context, _ string, _ string) error {
	return errors.New("not implemented")
}
