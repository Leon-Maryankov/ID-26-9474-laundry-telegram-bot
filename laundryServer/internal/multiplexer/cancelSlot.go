package multiplexer

import (
	"encoding/json"
	"net/http"

	"laundryServer/internal/slots"
)

type cancelSlotRequest struct {
	SlotId int `json:"slotId"`
	UserId int `json:"userId"`
}

type cancelSlotRespond struct {
	respond
}

func cancelSlotHandler(slots *slots.Slots) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.Write(incorrectRequestRespond)

			return
		}

		cancelSlotRequest := &cancelSlotRequest{}
		if err := json.NewDecoder(request.Body).Decode(cancelSlotRequest); err != nil {
			writer.Write(incorrectRequestRespond)

			return
		}

		err := slots.ReleaseSlot(cancelSlotRequest.SlotId, cancelSlotRequest.UserId)
		if err != nil {
			writer.Write(incorrectRequestRespond)

			return
		}

		cancelSlotRespond, _ := json.Marshal(&cancelSlotRespond{
			respond: respond{
				Status: good,
			},
		})

		writer.Write(cancelSlotRespond)
	}
}
