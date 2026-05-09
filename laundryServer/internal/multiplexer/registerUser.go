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

// добавил обертку для установки заголовка, resp.json ожидает
// json, а по дефолту ставился plane/text
// поэтому в ручную ставил "application/json"

func writeJSONResponse(writer http.ResponseWriter, data []byte, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write(data)
}

func registerUserHandler(userss *users.Users) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writeJSONResponse(writer, incorrectRequestRespond, http.StatusMethodNotAllowed)
			return
		}

		registerUserRequest := &registerUserRequest{}
		if err := json.NewDecoder(request.Body).Decode(registerUserRequest); err != nil {
			writeJSONResponse(writer, incorrectRequestRespond, http.StatusBadRequest)
			return
		}

		user := &users.User{
			TelegramId: registerUserRequest.TelegramId,
			Name:       registerUserRequest.Name,
			Surname:    registerUserRequest.Surname,
			Email:      registerUserRequest.Email,
		}

		status := good
		if err := userss.AddUser(user); err != nil {
			status = bad
		}

		response, _ := json.Marshal(&registerUserRespond{
			respond{Status: status},
		})

		writeJSONResponse(writer, response, http.StatusOK)
	}
}
