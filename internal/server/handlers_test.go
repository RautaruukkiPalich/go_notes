package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rautaruukkipalich/go_notes/internal/model"
	mockcachestore "github.com/rautaruukkipalich/go_notes/internal/store/mockcachestore"
	mockstore "github.com/rautaruukkipalich/go_notes/internal/store/mocksqlstore"
)

const authToken1 = "testcase1"
const authToken2 = "testcase2"
const invalidAuthToken = "invalid_auth_token"

var testCases = []model.Note{
	{
		Body: "testcase1",
		AuthorID: 1,
		IsPublic: true,
	},
	{
		Body: "testcase2",
		AuthorID: 2,
		IsPublic: true,
	},
	{
		Body: "testcase3",
		AuthorID: 1,
		IsPublic: false,
	},
}


func setTokens(s *Server, t *testing.T) {
	// valid
	s.RedisSetUser(authToken1, model.TestUser(t))

	// 
	u := model.TestUser(t)
	u.ID = 2
	s.RedisSetUser(authToken2, u)

	//invalid
	u = model.TestUser(t)
	u.TokenTTL = time.Now().UTC().Add(time.Minute * -1)
	s.RedisSetUser(invalidAuthToken, u)
}

func setNotes(s *Server, t *testing.T) []int {
	ids := []int{}
	for _, tc := range testCases {
		note := tc 
		s.PostNote()
		n, err := s.store.Note().Set(&note)
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, n.ID)
	}
	// n, _ := s.store.Note().GetNotes(0, "", 0, 5, 0)
	// fmt.Println(n)
	// fmt.Println(ids)
	return ids
}


func TestServer_HandlePostNote(t *testing.T) {

	store, err := mockstore.New()
	if err != nil {
		t.Fatal(err)
	}
	cache, err := mockcachestore.New()
	if err != nil {
		t.Fatal(err)
	}

	s := NewServer(store, cache)
	setTokens(s, t)

	testCases := []struct {
		name string
		payload any
		authHeader string
		expectedCode int
	}{
		{
			name: "valid post form",
			payload: map[string]any{
				"body": "test body 1",
				"is_public": true,
			},
			authHeader: authToken1,
			expectedCode: http.StatusOK,
		},
		{
			name: "valid post form 2",
			payload: map[string]any{
				"body": "test body 2",
				"is_public": false,
			},
			authHeader: authToken1,
			expectedCode: http.StatusOK,
		},
		{
			name: "empty auth token",
			payload: map[string]any{
				"body": "test body 3",
				"is_public": true,
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "invalid auth token",
			payload: map[string]any{
				"body": "test body 4",
				"is_public": true,
			},
			authHeader: invalidAuthToken,
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			var json_data bytes.Buffer
			_ = json.NewEncoder(&json_data).Encode(tc.payload)

			req, _ := http.NewRequest(http.MethodPost, "/notes", &json_data)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.authHeader))

			s.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)
		})
	}
}


func TestServer_HandleGetNote(t *testing.T) {
	store, err := mockstore.New()
	if err != nil {
		t.Fatal(err)
	}
	cache, err := mockcachestore.New()
	if err != nil {
		t.Fatal(err)
	}

	s := NewServer(store, cache)
	setTokens(s, t)
	ids := setNotes(s, t)

	testCases := []struct {
		name string
		payload int
		authHeader string
		expectedCode int
	}{
		{
			name: "valid id",
			payload: ids[0],
			authHeader: authToken2,
			expectedCode: http.StatusOK,
		},
		{
			name: "non public note",
			payload: ids[2],
			authHeader: authToken2,
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			var json_data bytes.Buffer
			_ = json.NewEncoder(&json_data).Encode(tc.payload)

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/notes/%d", tc.payload), &json_data)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.authHeader))

			s.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedCode, rr.Code)
		})
	}

}


func TestServer_HandleGetNotes(t *testing.T) {
	store, err := mockstore.New()
	if err != nil {
		t.Fatal(err)
	}
	cache, err := mockcachestore.New()
	if err != nil {
		t.Fatal(err)
	}

	s := NewServer(store, cache)
	setTokens(s, t)
	
}


func TestServer_HandlePatchNote(t *testing.T) {
	store, err := mockstore.New()
	if err != nil {
		t.Fatal(err)
	}
	cache, err := mockcachestore.New()
	if err != nil {
		t.Fatal(err)
	}

	s := NewServer(store, cache)
	setTokens(s, t)
	
}

func TestServer_HandleDeleteNote(t *testing.T) {
	store, err := mockstore.New()
	if err != nil {
		t.Fatal(err)
	}
	cache, err := mockcachestore.New()
	if err != nil {
		t.Fatal(err)
	}

	s := NewServer(store, cache)
	setTokens(s, t)
	
}
