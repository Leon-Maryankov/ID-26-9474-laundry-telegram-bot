package multiplexer

import (
	"encoding/json"
	"net/http"

	"laundryServer/internal/users"
)

type registerUserRequest struct {
	TelegramId int    `json:"telegramId"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Email      string `json:"email"`
}

type registerUserRespond struct {
	respond
}

func registerUserHandler(userss *users.Users) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.Write(incorrectRequestRespond)

			return
		}

		registerUserRequest := &registerUserRequest{}
		if err := json.NewDecoder(request.Body).Decode(registerUserRequest); err != nil {
			writer.Write(incorrectRequestRespond)

			return
		}

		user := &users.User{
			TelegramId: registerUserRequest.TelegramId,
			Name:       registerUserRequest.Name,
			Surname:    registerUserRequest.Surname,
			Email:      registerUserRequest.Email,
		}

		registerUserRespond, _ := json.Marshal(&registerUserRespond{
			respond{
				Status: func() status {
					err := userss.AddUser(user)
					if err == nil {
						return good
					} else {
						return bad
					}
				}(),
			},
		})

		writer.Write(registerUserRespond)
	}
}
