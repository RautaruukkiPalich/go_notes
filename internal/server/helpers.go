package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
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

func getClaimsFromJWT(token string)(map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(
		token,
		claims, 
		func(token *jwt.Token) (interface{}, error) {return []byte{}, nil},
	)
	return claims, err
}

func getUserInfoFromAuth(token string, user *model.User) (*model.User, error) {
	client := &http.Client{}

	auth_url := "http://localhost:8080/me"
	auth_request, err := http.NewRequest(http.MethodGet, auth_url, nil)
	
	if err != nil {
		// TODO: handle error here
		return nil, fmt.Errorf("404")
	}

	auth_request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(auth_request)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("check token: status code = %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
    if err != nil {
		return nil, err
    }

	defer res.Body.Close()

	if err := json.Unmarshal(body, user); err != nil {
		return nil, err
	}
	return user, nil
}

func getUserByToken(token string) (*model.User, error) {

	var user model.User
	
	claims, _ := getClaimsFromJWT(token)
	exp := claims["exp"]
	if exp == nil {
		// TODO: handle error token expired
		return &user, errors.New("token expired")
	}
	user.TokenTTL = time.Unix(int64(exp.(float64)), 0)

	return getUserInfoFromAuth(token, &user)
}

func validateUserData(user *model.User) error {
	if user.ID == 0 {return ErrInvalidToken}
	if user.TokenTTL.Unix() < time.Now().UTC().Unix() {return ErrInvalidToken}
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
	s.JSONrespond(w, r, err.Code, err)
}

func (s *Server) JSONrespond(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.WriteHeader(code)
	if data != nil {
		w.Header().Add("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.logger.Printf("error: %v", err)
		}
	}
}