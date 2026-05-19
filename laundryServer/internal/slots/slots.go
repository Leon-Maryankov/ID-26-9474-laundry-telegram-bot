package slots

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"time"
)

type Slots struct {
	database *sql.DB
}

type MachineType int

const (
	Washing MachineType = iota
	Dryer
)

func FromString(string string) MachineType {
	switch string {
	case "washing":
		return Washing

	case "dryer":
		return Dryer
	}

	panic(fmt.Errorf("MachineType.FromString: unknown type"))
}

func (machineType MachineType) String() string {
	switch machineType {
	case Washing:
		return "washing"

	case Dryer:
		return "dryer"
	}

	panic(fmt.Errorf("MachineType.String: unknown type"))
}

func (machineType MachineType) MarshalJSON() ([]byte, error) {
	return json.Marshal(machineType.String())
}

const (
	amountOfWashingMachines = 3
	amountOfDryerMachines   = 1
)

const (
	beginHour = 6
	endHour   = 23
)

var excludedHours = []int{16}

type Slot struct {
	Id          int
	Date        time.Time
	MachineId   int
	MachineType MachineType
	UserId      int `json:"-"`
}

func NewSlots(database *sql.DB) (*Slots, error) {
	_, err := database.Exec(
		`CREATE TABLE IF NOT EXISTS slots (
			id INTEGER UNIQUE NOT NULL PRIMARY KEY,
			date INTEGER NOT NULL,
			machineId INTEGER NOT NULL,
			machineType TEXT NOT NULL,
			userId INTEGER,
			UNIQUE(date, machineId, machineType),
			FOREIGN KEY (userId) REFERENCES users (id)
		);`,
	)
	if err != nil {
		return nil, fmt.Errorf("slots.NewSlots: failed to initialize database: %w", err)
	}

	result := &Slots{
		database,
	}

	go func() {
		for {
			now := time.Now()
			monday := now
			for monday.Weekday() != time.Monday {
				monday = monday.AddDate(0, 0, 1)
			}
			monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 22, 15, 0, 0, monday.Location())
			if !monday.After(now) {
				monday = monday.AddDate(0, 0, 7)
			}

			time.Sleep(time.Until(monday))

			err := result.removeAllSlots()
			if err != nil {
				panic(fmt.Errorf("SLOTS WEREN'T REMOVED: %w", err)) // TODO: Handle this properly.
			}

			err = result.generateSlots(monday)
			if err != nil {
				panic(fmt.Errorf("SLOTS WEREN'T ADDED: %w", err)) // TODO: Handle this properly.
			}
		}
	}()

	return result, nil
}

func (slots *Slots) removeAllSlots() error {
	_, err := slots.database.Exec("DELETE FROM slots;")
	if err != nil {
		return fmt.Errorf("slots.removeAllSlots: failed to remove all slots: %w", err)
	}

	return nil
}

func (slots *Slots) addEmptySlot(slot *Slot) error {
	_, err := slots.database.Exec(
		`INSERT INTO slots (date, machineId, machineType, userId)
		VALUES (?, ?, ?, NULL);
	`, slot.Date.Unix(), slot.MachineId, slot.MachineType.String())
	if err != nil {
		return fmt.Errorf("slots.addEmptySlot: failed to add empty slot: %w", err)
	}

	return nil
}

func (slots *Slots) generateSlots(start time.Time) error {
	for day := range 7 {
		currentDay := start.AddDate(0, 0, day)

		for hour := beginHour; hour < endHour; hour++ {
			if _, ok := slices.BinarySearch(excludedHours, hour); ok {
				continue
			}

			date := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(), hour, 0, 0, 0, currentDay.Location())

			for washingMachineId := range amountOfWashingMachines {
				err := slots.addEmptySlot(&Slot{
					Date:        date,
					MachineId:   washingMachineId,
					MachineType: Washing,
				})
				if err != nil {
					return fmt.Errorf("slots.generateSlots: failed to generate slot: %w", err)
				}
			}

			for dryerMachineId := range amountOfDryerMachines {
				err := slots.addEmptySlot(&Slot{
					Date:        date,
					MachineId:   dryerMachineId,
					MachineType: Dryer,
				})
				if err != nil {
					return fmt.Errorf("slots.generateSlots: failed to generate slot: %w", err)
				}
			}
		}
	}

	return nil
}

func (slots *Slots) GetAvailableSlots() ([]Slot, error) {
	rows, err := slots.database.Query(
		`SELECT id, date, machineId, machineType
		FROM slots
		WHERE userId IS NULL;
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("slots.GetAvailableSlots: failed to query rows: %w", err)
	}
	defer rows.Close()

	result := []Slot{}
	for rows.Next() {
		var (
			id          int
			date        int64
			machineId   int
			machineType string
		)

		err := rows.Scan(
			&id,
			&date,
			&machineId,
			&machineType,
		)
		if err != nil {
			return nil, fmt.Errorf("slots.GetAvailableSlots: failed to scan rows: %w", err)
		}

		slot := Slot{
			Id:          id,
			Date:        time.Unix(date, 0),
			MachineId:   machineId,
			MachineType: FromString(machineType),
		}

		result = append(result, slot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("slots.GetAvailableSlots: rows' error: %w", err)
	}

	return result, nil
}

var (
	SlotNotFound        = fmt.Errorf("slots.OccupySlots: slot not found")
	SlotAlreadyOccupied = fmt.Errorf("slots.OccupySlots: slot already occupied")
	WashingLimit        = fmt.Errorf("slots.OccupySlots: washing machine limit exceeded")
	DryerLimit          = fmt.Errorf("slots.OccupySlots: dryer limit exceeded")
)

func (slots *Slots) OccupySlot(slotId int, userId int) error {
	result, err := slots.database.Exec(
		`UPDATE slots
		SET userId = ?
		WHERE id = ?
			AND userId IS NULL
			AND (
				(machineType = 'washing' AND (
					SELECT COUNT(*) FROM slots
					WHERE userId = ? AND machineType = 'washing'
				) < 2)
			OR
				(machineType = 'dryer' AND (
					SELECT COUNT(*) FROM slots
					WHERE userId = ? AND machineType = 'dryer'
				) < 1)
			);`,
		userId, slotId, userId, userId,
	)
	if err != nil {
		return fmt.Errorf("slots.OccupySlot: failed to update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("slots.OccupySlot: failed to affect row: %w", err)
	}

	if rowsAffected > 0 {
		return nil
	}

	var machineType string
	var existingUserId sql.NullInt64

	err = slots.database.QueryRow(`
		SELECT machineType, userId
		FROM slots
		WHERE id = ?;
	`, slotId).Scan(&machineType, &existingUserId)

	if err == sql.ErrNoRows {
		return SlotNotFound
	}
	if err != nil {
		return fmt.Errorf("slots.OccupySlot: failed to select slot: %w", err)
	}

	if existingUserId.Valid {
		return SlotAlreadyOccupied
	}

	var count int

	switch FromString(machineType) {
	case Washing:
		err = slots.database.QueryRow(
			`SELECT COUNT(*) FROM slots
			WHERE userId = ? AND machineType = 'washing';`, userId).
			Scan(&count)
		if err != nil {
			return fmt.Errorf("slots.OccupySlot: failed to count washings: %w", err)
		}
		if count >= 2 {
			return WashingLimit
		}

	case Dryer:
		err = slots.database.QueryRow(
			`SELECT COUNT(*) FROM slots
			WHERE userId = ? AND machineType = 'dryer';`, userId).
			Scan(&count)
		if err != nil {
			return fmt.Errorf("slots.OccupySlot: failed to count dryers: %w", err)
		}
		if count >= 1 {
			return DryerLimit
		}
	}

	panic("SOMETHING WENT COMPLETELY WRONG!")
}

func (slots *Slots) GetUsersForNotification() ([]int, error) {
	now := time.Now()
	from := now.Add(10 * time.Minute).Unix()
	to := now.Add(11 * time.Minute).Unix()

	rows, err := slots.database.Query(
		`SELECT DISTINCT userId
		FROM slots
		WHERE userId IS NOT NULL
			AND date >= ?
			AND date < ?;`,
		from, to,
	)
	if err != nil {
		return nil, fmt.Errorf("slots.GetUsersForNotification: failed to query: %w", err)
	}
	defer rows.Close()

	result := []int{}
	for rows.Next() {
		userId := 0
		if err := rows.Scan(&userId); err != nil {
			return nil, fmt.Errorf("slots.GetUsersForNotification: failed to scan rows: %w", err)
		}

		result = append(result, userId)
	}

	return result, nil
}

func (slots *Slots) ReleaseSlot(slotId int, userId int) error {
	result, err := slots.database.Exec(
		`UPDATE slots
		SET userId = NULL
		WHERE id = ? AND userId = ?;`,
		slotId, userId,
	)
	if err != nil {
		return fmt.Errorf("slots.ReleaseSlot: failed to update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("slots.ReleaseSlot: failed to affect rows: %w", err)
	}

	if rowsAffected > 0 {
		return nil
	}

	var existingUserId sql.NullInt64

	err = slots.database.QueryRow(
		`SELECT userId FROM slots WHERE id = ?;`,
		slotId,
	).Scan(&existingUserId)

	if err == sql.ErrNoRows {
		return fmt.Errorf("slot.ReleaseSlot: slot not found.")
	}
	if err != nil {
		return fmt.Errorf("slots.ReleaseSlot: failed to select: %w", err)
	}

	if !existingUserId.Valid {
		return fmt.Errorf("slots.ReleaseSlot: slot is already free.")
	}

	if int(existingUserId.Int64) != userId {
		return fmt.Errorf("slots.ReleaseSlot: slot belongs to another user.")
	}

	return fmt.Errorf("slots.ReleaseSlot: unknown error.")
}

func (slots *Slots) GetSlotsAttachedToUser(id int) ([]Slot, error) {
	rows, err := slots.database.Query(
		`SELECT id, date, machineId, machineType, userId
		FROM slots
		WHERE userId = ?`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("slots.GetSlotsAttachedToUser: failed to query: %w", err)
	}
	defer rows.Close()

	result := []Slot{}

	for rows.Next() {
		var (
			slot        Slot
			timestamp   int64
			machineType string
			userId      sql.NullInt64
		)

		err := rows.Scan(
			&slot.Id,
			&timestamp,
			&slot.MachineId,
			&machineType,
			&userId,
		)
		if err != nil {
			return nil, fmt.Errorf("slots.GetSlotsAttachedToUser: scan failed: %w", err)
		}

		slot.Date = time.Unix(timestamp, 0)
		slot.MachineType = FromString(machineType)

		if userId.Valid {
			slot.UserId = int(userId.Int64)
		}

		result = append(result, slot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("slots.GetSlotsAttachedToUser: rows error: %w", err)
	}

	return result, nil
}
