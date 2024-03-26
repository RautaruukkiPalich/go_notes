package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rautaruukkipalich/go_notes/internal/model"
)

type (
	userKey string
	
	notePostForm struct {
		Body string `json:"body"`
		IsPublic bool `json:"is_public"`
	}
)

const (
	UserKey userKey = "userID" 
)


func getUserInfoByToken(token string) (*model.User, error) {
	// user := mockUserInfo

	var user model.User

	res, err := checkAuth(token)

	if err != nil {
		return &user, err
	}

	if res.StatusCode != 200 {
		return &user, fmt.Errorf("check token: status code = %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
    if err != nil {
		return &user, err
    }

	fmt.Println(string(body))

	defer res.Body.Close()

	if err := json.Unmarshal(body, &user); err != nil {
		return &user, err
	}

	// if token expires -> request to sso -> get new token

	return &user, nil
}


func checkAuth(token string) (*http.Response, error) {
	client := &http.Client{}

	auth_url := "http://localhost:8080/me"
	auth_request, err := http.NewRequest(http.MethodGet, auth_url, nil)
	
	if err != nil {
		return nil, fmt.Errorf("404")
	}

	auth_request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return client.Do(auth_request)
}

func validateUserData(user *model.User) error {
	if user.ID == 0 {return ErrInvalidToken}
	// if user.TokenTTL.Unix() < time.Now().UTC().Unix() {return ErrInvalidToken}
	return nil
}

func getIdVarFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	strID := vars["id"]
	return strconv.Atoi(strID)
}


//limit, offset and body/author filters
func getFiltersFromRequest(r *http.Request) (int, int, int, string) {
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		offset = 0
	}

	filter_author, err := strconv.Atoi(r.FormValue("filter_author"))
	if err != nil {
		filter_author = 0
	}

	filterBody := r.FormValue("filter_body")

	return limit, offset, filter_author, filterBody
}

func (s *Server) GetNoteById(id int) (model.Note, error) {
	var note model.Note

	s.logger.Infof("get note %d from cache...", id)
	note, err := s.cache.Note().GetNoteById(id)
	if err == nil {
		return note, nil
	}

	s.logger.Infof("get note %d from db...", id)
	note, err = s.store.Note().GetNoteById(id)
	if err != nil {
		return note, err
	}

	s.logger.Infof("set note %d to cache...", id)
	if err := s.cache.Note().Set(&note); err != nil {
		return note, err
	}
	return note, err
}

func (s *Server) getUserIdFromContext(r *http.Request) int {
	ctx := r.Context()
	id, ok := ctx.Value(UserKey).(int)

	if !ok {
		return 0
	}
	return id
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, err errorResponse) {
	s.logger.Printf("error: %v", err)
	s.respond(w, r, err.Code, err)
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.WriteHeader(code)
	if data != nil {
		w.Header().Add("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.logger.Printf("error: %v", err)
		}
	}
}