package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func SetUpUserHandlerDep() *UserHandler {
	h := &UserHandler{
		UserDB:   &MockUserDB{},
		TokenMng: &MockTokenManager{},
	}
	return h
}

// ---------------- POST TESTS --------------------
func TestUserHandler_PostUserLogin(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"username": "John Doe"}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "Correct")

	handler := SetUpUserHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.LoginHandle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusOK)
	}
}

func TestUserHandler_PostUserLogin_IncorrectBodyRequest(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"username": -1}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpUserHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.LoginHandle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusBadRequest)
	}
}

func TestUserHandler_PostUserLogin_IncorrectUser(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"username": "ErrNoRows"}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpUserHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.LoginHandle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusUnauthorized)
	}
}

func TestUserHandler_PostUserLogin_IncorrectUserServerErr(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"username": "ServerError"}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpUserHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.LoginHandle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusInternalServerError)
	}
}

func TestUserHandler_PostUserLogin_IncorrectTokenInRequest(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"username": "Incorrect"}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpUserHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.LoginHandle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusInternalServerError)
	}
}

func TestUserHandler_PostUserReg(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"username": "Reg"}`)
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpUserHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.RegHandle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusCreated)
	}
}

func TestUserHandler_PostUserReg_IncorrectBodyRequest(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"username": -1}`)
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpUserHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.RegHandle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusBadRequest)
	}
}

func TestUserHandler_PostUserReg_IncorrectName(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"username": "John Doe"}`)
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpUserHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.RegHandle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusConflict)
	}

	expectedErrorMessage := "User with this name is already registered"
	actualErrorMessage := rr.Header().Get("X-Error-Message")
	if actualErrorMessage != expectedErrorMessage {
		t.Errorf("Received incorrect error message: received %v, expected %v",
			actualErrorMessage, expectedErrorMessage)
	}
}

func TestUserHandler_PostUserReg_ServerError(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"username": "ServerError"}`)
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpUserHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.RegHandle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusInternalServerError)
	}
}

func TestUserHandler_PostUserReg_ErrNoRows(t *testing.T) {
	// Arrange
	expenseJSON := []byte(`{"username": "ErrNoRows"}`)
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(expenseJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := SetUpUserHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.RegHandle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusConflict)
	}
}

// ---------------- END USER POST TESTS --------------------

// -------------- NOTALLOWEDMETHOD TESTS --------------
func TestUserHandler_UnknownMethodLogin_NotAllowed(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("NOTALLOWED", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "Correct")

	handler := SetUpUserHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.LoginHandle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestUserHandler_UnknownMethodReg_NotAllowed(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("NOTALLOWED", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", "Correct")

	handler := SetUpUserHandlerDep()

	rr := httptest.NewRecorder()

	// Act
	handler.RegHandle(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Отримано некоректний статус-код: отримано %v, очікувалося %v",
			status, http.StatusMethodNotAllowed)
	}
}

// -------------- END NOTALLOWEDMETHOD TESTS --------------
