package multiplexer

import (
	"encoding/json"
	"net/http"

	"laundryServer/internal/slots"
)

type getSlotsAttachedToUserRequest struct {
	UserId int `json:"userId"`
}

type getSlotsAttachedToUserRespond struct {
	respond

	Slots []slots.Slot `json:"slots"`
}

func getSlotsAttachedToUserHandler(slots *slots.Slots) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.Write(incorrectRequestRespond)
			return
		}

		getSlotsAttachedToUserRequest := &getSlotsAttachedToUserRequest{}
		if err := json.NewDecoder(request.Body).Decode(getSlotsAttachedToUserRequest); err != nil {
			writer.Write(incorrectRequestRespond)

			return
		}

		slots, err := slots.GetSlotsAttachedToUser(getSlotsAttachedToUserRequest.UserId)
		if err != nil {
			writer.Write(incorrectRequestRespond)

			return
		}

		getSlotsAttachedToUserRespond, _ := json.Marshal(&getSlotsAttachedToUserRespond{
			respond: respond{
				Status: good,
			},
			Slots: slots,
		})

		writer.Write(getSlotsAttachedToUserRespond)
	}
}
