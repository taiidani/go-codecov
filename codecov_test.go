// Package codecov defines an API client for the CodeCov REST API.
//
// https://docs.codecov.io/reference
package codecov

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestClient_doRequest(t *testing.T) {
	baseRequest, _ := http.NewRequest("GET", "/gh/", nil)

	type args struct {
		request  *http.Request
		response interface{}
	}
	tests := []struct {
		name         string
		args         args
		handler      func(w http.ResponseWriter, r *http.Request)
		wantResponse Response
		wantErr      bool
	}{
		{
			name: "Successful",
			args: args{
				request:  baseRequest,
				response: &Response{},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, `{ "meta": { "status": 200 } }`)
			},
			wantResponse: Response{
				Meta: Meta{
					Status: 200,
				},
			},
		},
		{
			name: "Error",
			args: args{
				request:  baseRequest,
				response: &Response{},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(404)
				fmt.Fprintln(w, `{ "meta": { "status": 404 }, "error": { "reason": "Item not found." } } }`)
			},
			wantErr: true,
		},
		{
			name: "Bad Response",
			args: args{
				request:  baseRequest,
				response: &Response{},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, `uhhh`)
			},
			wantErr: true,
		},
		{
			name: "Timeout",
			args: args{
				request: baseRequest.Clone(func() context.Context {
					ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
					return ctx
				}()),
				response: &Response{},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(time.Second * 20)
				fmt.Fprintln(w, `{ "meta": { "status": 200 } }`)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Running server in separate goroutine so the Timeout test does not block
			var srv *httptest.Server
			ready := make(chan bool)
			done := make(chan bool)
			go func() {
				mux := http.NewServeMux()
				mux.HandleFunc("/gh/", tt.handler)
				srv = httptest.NewServer(mux)
				defer srv.Close()
				ready <- true // Ready to receive requests
				<-done        // Test complete, shut down server
			}()

			<-ready
			c := NewClient("testing")
			c.SetEndpoint(url.URL{Scheme: "http", Host: strings.ReplaceAll(srv.URL, "http://", "")})

			if err := c.doRequest(tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("Client.doRequest() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(tt.args.response, &tt.wantResponse) {
				t.Errorf("Received %#v, want %#v", tt.args.response, tt.wantResponse)
			}
			done <- true
		})
	}
}
