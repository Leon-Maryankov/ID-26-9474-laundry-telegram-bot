package multiplexer

import (
	"encoding/json"
	"net/http"

	"laundryServer/internal/users"
)

type doesUserExistRequest struct {
	TelegramId int `json:"telegramId"`
}

type doesUserExistRespond struct {
	respond

	Answer bool `json:"answer"`
	User *users.User `json:"user,omitempty"`
}

func doesUserExistHandler(users *users.Users) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writeJSONResponse(writer, incorrectRequestRespond, http.StatusMethodNotAllowed)
			return
		}

		doesUserExistRequest := &doesUserExistRequest{}
		if err := json.NewDecoder(request.Body).Decode(doesUserExistRequest); err != nil {
			writeJSONResponse(writer, incorrectRequestRespond, http.StatusBadRequest)
			return
		}

		user, err := users.GetUserByTelegramId(doesUserExistRequest.TelegramId)
		if err != nil {
			writeJSONResponse(writer, incorrectRequestRespond, http.StatusInternalServerError)
			return
		}

		doesUserExistRespond, _ := json.Marshal(&doesUserExistRespond{
			respond: respond{
				Status: good,
			},
			Answer: user != nil,
			User: user,
		})

		writeJSONResponse(writer, doesUserExistRespond, http.StatusOK)
	}
}
