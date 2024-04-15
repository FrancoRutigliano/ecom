package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	typePayload "github.com/FrancoRutigliano/ecom/types"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserstore{}
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := typePayload.RegisterUserPayload{
			FirstName: "user",
			LastName:  "123",
			Email:     "invalid",
			Password:  "asd",
		}

		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/register", handler.handleRegister)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should correctly the user", func(t *testing.T) {
		payload := typePayload.RegisterUserPayload{
			FirstName: "user",
			LastName:  "123",
			Email:     "valid@gmail.com",
			Password:  "asd",
		}

		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/register", handler.handleRegister)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}

type mockUserstore struct{}

func (m *mockUserstore) GetUserByEmail(email string) (*typePayload.User, error) {
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserstore) GetUserByID(id int) (*typePayload.User, error) {
	return nil, nil
}

func (m *mockUserstore) CreateUser(user typePayload.User) error {
	return nil
}
