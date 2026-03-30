package multiplexer

import (
	"encoding/json"
	"net/http"

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

func NewMultiplexer(token string, users *users.Users) *Multiplexer {
	result := http.NewServeMux()
	result.HandleFunc("/" + token + "/doesUserExist", doesUserExistHandler(users))
	result.HandleFunc("/" + token + "/registerUser", registerUserHandler(users))

	return result
}
