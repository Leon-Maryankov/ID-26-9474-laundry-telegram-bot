package multiplexer

import (
	"encoding/json"
	"net/http"

	"laundryServer/internal/slots"
)

type getAvailableSlotsRequest struct{}

type getAvailableSlotsRespond struct {
	respond

	Slots []slots.Slot
}

func getAvailableSlotsHandler(slots *slots.Slots) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		slots, err := slots.GetAvailableSlots()
		if err != nil {
			writeJSONResponse(writer, incorrectRequestRespond, http.StatusInternalServerError)
			return
		}

		getAvailableSlotsRespond, _ := json.Marshal(&getAvailableSlotsRespond{
			respond: respond{
				Status: good,
			},
			Slots: slots,
		})

		writeJSONResponse(writer, getAvailableSlotsRespond, http.StatusOK)
	}
}
