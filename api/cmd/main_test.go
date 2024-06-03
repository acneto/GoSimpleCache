package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var fakeTCPGet = func(command string) (string, error) {
	if command == "GET my_key" {
		return "my_value", nil
	}
	return "", fmt.Errorf("error from fake tcp server")
}

var fakeTCPSet = func(command string) (string, error) {
	if command == "SET my_key my_value" {
		return "OK\n", nil
	}
	return "", fmt.Errorf("error from fake tcp server")
}

func TestHandleSet(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		value          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "No key provided",
			key:            "",
			value:          "my_value",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Key and value are required\n",
		},
		{
			name:           "No value provided",
			key:            "my_key",
			value:          "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Key and value are required\n",
		},
		{
			name:           "Valid key and value",
			key:            "my_key",
			value:          "my_value",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK\n",
		},
		{
			name:           "Error from TCP server",
			key:            "invalidkey",
			value:          "invalidvalue",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "error from fake tcp server\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/set?key="+tt.key+"&value="+tt.value, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := handleSet(fakeTCPSet)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestHandleGet(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "No key",
			key:            "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Key is required\n",
		},
		{
			name:           "Valid key",
			key:            "my_key",
			expectedStatus: http.StatusOK,
			expectedBody:   "my_value",
		},
		{
			name:           "Error from TCP server",
			key:            "invalidkey",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "error from fake tcp server\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/get?key="+tt.key, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := handleGet(fakeTCPGet)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
