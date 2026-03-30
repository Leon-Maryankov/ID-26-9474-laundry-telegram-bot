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
}

func doesUserExistHandler(users *users.Users) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.Write(incorrectRequestRespond)

			return
		}

		doesUserExistRequest := &doesUserExistRequest{}
		if err := json.NewDecoder(request.Body).Decode(doesUserExistRequest); err != nil {
			writer.Write(incorrectRequestRespond)

			return
		}

		user, err := users.GetUserByTelegramId(doesUserExistRequest.TelegramId)
		if err != nil {
			writer.Write(incorrectRequestRespond)

			return
		}

		doesUserExistRespond, _ := json.Marshal(&doesUserExistRespond{
			respond: respond{
				Status: good,
			},
			Answer: user != nil,
		})
		writer.Write(doesUserExistRespond)
	}
}
