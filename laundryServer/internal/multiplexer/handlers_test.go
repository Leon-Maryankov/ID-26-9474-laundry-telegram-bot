package multiplexer

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"laundryServer/internal/slots"
	"laundryServer/internal/users"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sql.DB, *users.Users, *slots.Slots) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	usersRepo, _ := users.NewUsers(db)
	slotsRepo, _ := slots.NewSlots(db)
	usersRepo.AddUser(&users.User{
		TelegramId: 777,
		Name:       "Test",
		Surname:    "Bot",
		Email:      "bot@test.com",
	})
	return db, usersRepo, slotsRepo
}

func TestRegisterUserHandler(t *testing.T) {
	db, usersRepo, _ := setupTestDB(t)
	defer db.Close()

	handler := registerUserHandler(usersRepo)

	reqBody := []byte(`{"telegramId": 888, "name": "New", "surname": "User", "email": "new@test.com"}`)
	req := httptest.NewRequest("POST", "/token/registerUser", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	handler(w, req)

	var resp registerUserRespond
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Status != good {
		t.Errorf("expected status 'good', got '%s'", resp.Status)
	}
}

func TestOccupySlotHandlerIntegration(t *testing.T) {
	db, usersRepo, slotsRepo := setupTestDB(t)
	defer db.Close()

	user, err := usersRepo.GetUserByTelegramId(777)
	if err != nil || user == nil {
		t.Fatal("test user not found")
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

	handler := occupySlotHandler(slotsRepo)
	reqBody := []byte(`{"slotId": ` + strconv.Itoa(slotId) + `, "userId": ` + strconv.Itoa(user.Id) + `}`)
	req := httptest.NewRequest("POST", "/token/occupySlot", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	handler(w, req)

	var resp occupySlotRespond
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Status != good {
		t.Errorf("expected good, got %v", resp)
	}

	var userIdDB sql.NullInt64
	err = db.QueryRow(`SELECT userId FROM slots WHERE id = ?`, slotId).Scan(&userIdDB)
	if err != nil {
		t.Fatal(err)
	}
	if !userIdDB.Valid || userIdDB.Int64 != int64(user.Id) {
		t.Errorf("slot not occupied correctly: userIdDB=%v, expected %d", userIdDB, user.Id)
	}
}
