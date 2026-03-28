package booking

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSeatAlreadyTaken = errors.New("seat already taken")
)

type Booking struct {
	ID        string    `json:"session_id"`
	MovieID   string    `json:"movie_id"`
	UserID    string    `json:"user_id"`
	SeatID    string    `json:"seat_id"`
	Status    int       `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}

type BookingStore interface {
	Book(b Booking) error
	Hold(b Booking) (Booking, error)
	ListBookings(movieID string) ([]Booking, error)
	Confirm(ctx context.Context, sessionID string, userID string) (Booking, error)
	Release(ctx context.Context, sessionID string, userID string) error
}
