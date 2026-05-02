package slots

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("cannot open test db: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, telegramId INTEGER UNIQUE, name TEXT, surname TEXT, email TEXT UNIQUE);`)
	if err != nil {
		t.Fatalf("cannot create users table: %v", err)
	}
	return db
}

func TestNewSlotsCreatesTable(t *testing.T) {
	db := setupTestDB(t)
	_, err := NewSlots(db)
	if err != nil {
		t.Fatal(err)
	}
	var count int
	err = db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE type='table' AND name='slots'`).Scan(&count)
	if err != nil || count == 0 {
		t.Error("table slots not created")
	}
}

func TestGetAvailableSlotsEmpty(t *testing.T) {
	db := setupTestDB(t)
	slotsRepo, _ := NewSlots(db)
	_, err := db.Exec(`DELETE FROM slots`)
	if err != nil {
		t.Fatal(err)
	}
	free, err := slotsRepo.GetAvailableSlots()
	if err != nil {
		t.Fatal(err)
	}
	if len(free) != 0 {
		t.Errorf("expected 0 slots, got %d", len(free))
	}
}

func TestOccupyAndReleaseSlot(t *testing.T) {
	db := setupTestDB(t)
	slotsRepo, _ := NewSlots(db)
	_, err := db.Exec(`INSERT INTO users (id, telegramId, name, surname, email) VALUES (1, 123, 'Test', 'User', 't@t.com')`)
	if err != nil {
		t.Fatal(err)
	}
	date := time.Now().Truncate(time.Hour)
	_, err = db.Exec(`INSERT INTO slots (date, machineId, machineType, userId) VALUES (?, 0, 'washing', NULL)`, date.Unix())
	if err != nil {
		t.Fatal(err)
	}
	var slotId int
	err = db.QueryRow(`SELECT id FROM slots WHERE date = ?`, date.Unix()).Scan(&slotId)
	if err != nil {
		t.Fatal(err)
	}

	err = slotsRepo.OccupySlot(slotId, 1)
	if err != nil {
		t.Fatal(err)
	}
	var userId sql.NullInt64
	err = db.QueryRow(`SELECT userId FROM slots WHERE id = ?`, slotId).Scan(&userId)
	if err != nil {
		t.Fatal(err)
	}
	if !userId.Valid || userId.Int64 != 1 {
		t.Errorf("slot not occupied, userId=%v", userId)
	}

	err = slotsRepo.ReleaseSlot(slotId, 1)
	if err != nil {
		t.Fatal(err)
	}
	err = db.QueryRow(`SELECT userId FROM slots WHERE id = ?`, slotId).Scan(&userId)
	if err != nil {
		t.Fatal(err)
	}
	if userId.Valid {
		t.Error("slot still occupied after release")
	}
}

func TestWashingLimit(t *testing.T) {
	db := setupTestDB(t)
	slotsRepo, _ := NewSlots(db)
	_, err := db.Exec(`INSERT INTO users (id, telegramId, name, surname, email) VALUES (2, 456, 'Limit', 'Tester', 'limit@t.com')`)
	if err != nil {
		t.Fatal(err)
	}
	date1 := time.Now().Truncate(time.Hour)
	date2 := date1.Add(time.Hour)
	date3 := date2.Add(time.Hour)
	for i, d := range []time.Time{date1, date2, date3} {
		_, err = db.Exec(`INSERT INTO slots (date, machineId, machineType, userId) VALUES (?, ?, 'washing', NULL)`, d.Unix(), i)
		if err != nil {
			t.Fatal(err)
		}
	}
	var ids [3]int
	for i := 0; i < 3; i++ {
		err = db.QueryRow(`SELECT id FROM slots WHERE machineId = ? AND date = ?`, i, []time.Time{date1, date2, date3}[i].Unix()).Scan(&ids[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	err = slotsRepo.OccupySlot(ids[0], 2)
	if err != nil {
		t.Fatal(err)
	}
	err = slotsRepo.OccupySlot(ids[1], 2)
	if err != nil {
		t.Fatal(err)
	}
	err = slotsRepo.OccupySlot(ids[2], 2)
	if err != WashingLimit {
		t.Errorf("expected WashingLimit, got %v", err)
	}
}
