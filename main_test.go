package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuscaCepHandler(t *testing.T) {
	// Note que não é necessário iniciar um servidor real aqui.

	testCases := []struct {
		name           string
		cep            string
		expectedStatus int
		expectedBody   string
	}{
		{"Valid CEP", "79052564", http.StatusOK, ""},
		{"Invalid CEP", "invalid", http.StatusUnprocessableEntity, "invalid zipcode"},
		{"CEP Not found", "11111111", http.StatusNotFound, "can not find zipcode"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "/cep/"+tc.cep, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(BuscaCepHandler)

			handler.ServeHTTP(rr, request)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			if tc.expectedBody != "" && !strings.Contains(rr.Body.String(), tc.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.expectedBody)
			}
		})
	}
}
