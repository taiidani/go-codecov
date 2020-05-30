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

func TestClient_ListRepositories(t *testing.T) {
	type test struct {
		name    string
		ctx     context.Context
		handler func(w http.ResponseWriter, r *http.Request)
		want    ListRepositoriesResponse
		wantErr bool
	}

	tests := []test{
		{
			name: "Successful",
			ctx:  context.Background(),
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, `
{
	"repos": [
		{
			"fork": null,
			"name": "guess-my-word",
			"language": "go",
			"activated": false,
			"private": false,
			"updatestamp": "2020-05-10 17:55:29.623184+00:00",
			"latest_commit": "2020-05-10 17:54:16",
			"branch": "master",
			"coverage": 82.07547,
			"repoid": "9384825"
		},
		{
			"fork": null,
			"name": "terraform-provider-jenkins",
			"language": "go",
			"activated": false,
			"private": false,
			"updatestamp": null,
			"latest_commit": null,
			"branch": "master",
			"coverage": null,
			"repoid": "9384835"
		}
	],
	"meta": {
		"status": 200,
		"limit": 20,
		"page": 1
	}
}`)
			},
			want: ListRepositoriesResponse{
				Repos: []Repository{
					{
						Name:        "guess-my-word",
						Branch:      "master",
						Coverage:    82.07547,
						RepoID:      "9384825",
						Updatestamp: &Time{Time: time.Date(2020, time.May, 10, 17, 55, 29, 623184000, time.UTC)},
						Language:    "go",
					},
					{
						Name:        "terraform-provider-jenkins",
						Branch:      "master",
						Coverage:    0,
						RepoID:      "9384835",
						Updatestamp: nil,
						Language:    "go",
					},
				},
				Meta: Meta{
					Status: 200,
				},
			},
		},
		{
			name: "Error",
			ctx:  context.Background(),
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(404)
				fmt.Fprintln(w, `
{
	"repos": [],
	"meta": {
		"status": 404
	},
	"error": {
		"reason": "Team not found.",
		"context": null
	  }
	}
}`)
			},
			wantErr: true,
		},
		{
			name: "Invalid request",
			ctx:  nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, `{}`)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/gh/", tt.handler)
			srv := httptest.NewServer(mux)
			defer srv.Close()

			c := NewClient("testing")
			c.SetEndpoint(url.URL{Scheme: "http", Host: strings.ReplaceAll(srv.URL, "http://", "")})

			response, err := c.ListRepositories(tt.ctx, "test-account")
			if (err != nil) != tt.wantErr {
				t.Error(err)
			} else if !reflect.DeepEqual(response, tt.want) {
				t.Errorf("Received %#v, want %#v", response, tt.want)
			}
		})
	}
}

func TestClient_GetRepository(t *testing.T) {
	type test struct {
		name    string
		ctx     context.Context
		handler func(w http.ResponseWriter, r *http.Request)
		want    GetRepositoryResponse
		wantErr bool
	}

	tests := []test{
		{
			name: "Successful",
			ctx:  context.Background(),
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, `
{
  "repo": {
    "using_integration": true,
    "name": "guess-my-word",
    "language": "go",
    "deleted": false,
    "bot_username": null,
    "activated": true,
    "private": true,
    "updatestamp": "2020-05-10T17:55:29.623184+00:00",
    "branch": "master",
    "active": true
  },
  "meta": {
    "status": 200
  }
}`)
			},
			want: GetRepositoryResponse{
				Repo: Repository{
					Name:             "guess-my-word",
					Branch:           "master",
					RepoID:           "",
					Updatestamp:      &Time{Time: time.Date(2020, time.May, 10, 17, 55, 29, 623184000, time.UTC)},
					UsingIntegration: true,
					Activated:        true,
					Private:          true,
					Language:         "go",
				},
				Meta: Meta{
					Status: 200,
				},
			},
		},
		{
			name: "Error",
			ctx:  context.Background(),
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(404)
				fmt.Fprintln(w, `
{
	"repos": [],
	"meta": {
		"status": 404
	},
	"error": {
		"reason": "GitHub API: Not Found",
		"context": null
	  }
	}
}`)
			},
			wantErr: true,
		},
		{
			name: "Invalid request",
			ctx:  nil,
			handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, `{}`)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/gh/", tt.handler)
			srv := httptest.NewServer(mux)
			defer srv.Close()

			c := NewClient("testing")
			c.SetEndpoint(url.URL{Scheme: "http", Host: strings.ReplaceAll(srv.URL, "http://", "")})

			response, err := c.GetRepository(tt.ctx, "test-account", "test-repo")
			if (err == nil) == tt.wantErr {
				t.Error(err)
			} else if !reflect.DeepEqual(response, tt.want) {
				t.Errorf("Received %#v, want %#v", response, tt.want)
			}
		})
	}
}
