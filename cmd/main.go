package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/internal/booking"
	"example.com/internal/utils"
	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("Welcome to the Cinema Booking System!")
	fmt.Println("Rest API")

	mux := http.NewServeMux() // to send requests to the right handler
	fmt.Println("Server is running on http://localhost:8080")
	mux.HandleFunc("GET /movies", listMoviesHandler)
	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	store := booking.NewRedisStore(redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	}))
	svc := booking.NewService(store)
	bookingHandler := booking.NewHandler(svc)

	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", bookingHandler.HoldSeat)

	mux.HandleFunc("PUT /sessions/{sessionID}/confirm", bookingHandler.ConfirmSession)
	mux.HandleFunc("DELETE /sessions/{sessionID}", bookingHandler.ReleaseSession)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

// Maybe move to a separate file?
type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}

var movies = []movieResponse{
	{ID: "1", Title: "The Shawshank Redemption", Rows: 10, SeatsPerRow: 15},
	{ID: "2", Title: "The Godfather", Rows: 12, SeatsPerRow: 10},
}

func listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, movies)
}
