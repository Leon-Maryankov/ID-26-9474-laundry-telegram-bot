package multiplexer

import (
	"encoding/json"
	"net/http"

	"laundryServer/internal/slots"
	"laundryServer/internal/users"
)

type getNotificationsRequest struct {}

type getNotificationsRespond struct {
	respond

	Ids []int `json:"ids"`
}

func getNotificationsHandler(users *users.Users, slots*slots.Slots) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userIds, err := slots.GetUsersForNotification()
		if err != nil {
			writer.Write(incorrectRequestRespond)

			return
		}

		result := []int{}
		for _, id := range userIds {
			user, err := users.GetUserById(id)
			if err != nil {
				writer.Write(incorrectRequestRespond)
				return
			}

			if user != nil {
				result = append(result, user.TelegramId)
			}
		}

		responseBytes, _ := json.Marshal(&getNotificationsRespond{
			respond: respond{
				Status: good,
			},
			Ids: result,
		})

		writer.Write(responseBytes)
	}
}
