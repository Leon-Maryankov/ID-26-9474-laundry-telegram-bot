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

type Multiplexer = http.ServeMux

func NewMultiplexer(token string, users *users.Users, slots *slots.Slots) *Multiplexer {
	result := http.NewServeMux()
	result.HandleFunc("/" + token + "/doesUserExist", doesUserExistHandler(users))
	result.HandleFunc("/" + token + "/registerUser", registerUserHandler(users))
	result.HandleFunc("/" + token + "/getAvailableSlots", getAvailableSlotsHandler(slots))
	result.HandleFunc("/" + token + "/occupySlot", occupySlotHandler(slots))
	result.HandleFunc("/" + token + "/cancelSlot", cancelSlotHandler(slots))
	result.HandleFunc("/" + token + "/getSlotsAttachedToUser", getSlotsAttachedToUserHandler(slots))
	result.HandleFunc("/" + token + "/getNotifications", getNotificationsHandler(users, slots))

	return result
}
