package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"some-http-server/internal/transport"
	"testing"
)

func Test_Demonstration(t *testing.T) {

	tests := []struct {
		name     string
		req      *http.Request
		wantCode int
	}{
		// TODO fix test cases
		{
			". Success on POST to /quotes",
			httptest.NewRequest("POST", "http://localhost:8888/quotes", ),
			http.StatusOK,
		},
		{
			". Success on POST to /quotes",
			httptest.NewRequest("POST", "http://localhost:8888/quotes", ),
			http.StatusOK,
		},
		{
			". Success on POST to /quotes",
			httptest.NewRequest("POST", "http://localhost:8888/quotes", ),
			http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := "add service mock here" //TODO mock svc here

			r := mux.NewRouter()
			handler := transport.NewHandler(svc)
			handler.Register(r)

			respRec := httptest.NewRecorder()
			r.ServeHTTP(respRec, tt.req)

			resp := respRec.Result()
			body, _ := io.ReadAll(resp.Body)

			if respRec.Code != tt.wantCode {
				t.Errorf("code = %d, wantCode %d", respRec.Code, tt.wantCode)
				return
			}
			fmt.Printf("Status Code: '%d'\nContent-Type: '%s'\nBody: '%s'\n",
				resp.StatusCode, resp.Header.Get("Content-Type"), string(body))
		})
	}
}
