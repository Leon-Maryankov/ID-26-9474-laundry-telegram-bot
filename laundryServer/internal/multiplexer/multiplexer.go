package multiplexer

import (
	"encoding/json"
	"net/http"

	"laundryServer/internal/slots"
	"laundryServer/internal/users"
)

type status string
const (
	good status = "good"
	bad status = "bad"
)

type respond struct {
	Status status `json:"status"`
}

var incorrectRequestRespond = func() []byte {
	type incorrectRequestRespond struct {
		respond

		Reason string `json:"reason"`
	}

	result, _ := json.Marshal(incorrectRequestRespond{
		respond: respond{
			Status: bad,
		},
		Reason: "Incorrent request.",
	})

	return result
}()

// CORS middleware wrapper
func corsMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}

type Multiplexer = http.ServeMux

func NewMultiplexer(token string, users *users.Users, slots *slots.Slots) *Multiplexer {
	result := http.NewServeMux()
	result.HandleFunc("/" + token + "/doesUserExist", corsMiddleware(doesUserExistHandler(users)))
	result.HandleFunc("/" + token + "/registerUser", corsMiddleware(registerUserHandler(users)))
	result.HandleFunc("/" + token + "/getAvailableSlots", corsMiddleware(getAvailableSlotsHandler(slots)))
	result.HandleFunc("/" + token + "/occupySlot", corsMiddleware(occupySlotHandler(slots)))
	result.HandleFunc("/" + token + "/cancelSlot", corsMiddleware(cancelSlotHandler(slots)))
	result.HandleFunc("/" + token + "/getSlotsAttachedToUser", corsMiddleware(getSlotsAttachedToUserHandler(slots)))
	result.HandleFunc("/" + token + "/getNotifications", corsMiddleware(getNotificationsHandler(users, slots)))

	return result
}
