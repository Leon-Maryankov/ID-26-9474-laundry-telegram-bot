package multiplexer

import (
	"encoding/json"
	"net/http"

	"laundryServer/internal/slots"
)

type occupySlotRequest struct {
	SlotId int `json:"slotId"`
	UserId int `json:"userId"`
}

type reason = string

const (
	tooManyWashingMachines reason = "too many washing machines"
	tooManyDryerMachines reason = "too many dryer machines"
)

type occupySlotRespond struct {
	respond

	Reason reason `json:"reason,omitempty"`
}

func occupySlotHandler(slotss *slots.Slots) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.Write(incorrectRequestRespond)

			return
		}

		occupySlotRequest := &occupySlotRequest{}
		if err := json.NewDecoder(request.Body).Decode(occupySlotRequest); err != nil {
			writer.Write(incorrectRequestRespond)

			return
		}

		err := slotss.OccupySlot(occupySlotRequest.SlotId, occupySlotRequest.UserId)
		if err != nil {
			switch err {
			case slots.WashingLimit:
				occupySlotRespond, _ := json.Marshal(&occupySlotRespond{
					respond: respond{
						Status: bad,
					},
					Reason: tooManyWashingMachines,
				})

				writer.Write(occupySlotRespond)

				return

			case slots.DryerLimit:
				occupySlotRespond, _ := json.Marshal(&occupySlotRespond{
					respond: respond{
						Status: bad,
					},
					Reason: tooManyDryerMachines,
				})

				writer.Write(occupySlotRespond)

				return

			default:
				writer.Write(incorrectRequestRespond)

				return
			}
		}

		occupySlotRespond, _ := json.Marshal(&occupySlotRespond{
			respond: respond{
				Status: good,
			},
		})

		writer.Write(occupySlotRespond)
	}
}
