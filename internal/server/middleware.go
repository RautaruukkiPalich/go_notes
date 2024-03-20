package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rautaruukkipalich/go_notes/internal/model"
)

var mockUserInfo = &model.User{
	ID: 1,
	TokenTTL: time.Time(time.Now().UTC().Add(time.Minute*10)),
	Valid: true,
}

func getUserInfoByToken(token string) (*model.User, error) {
	user := mockUserInfo
	// if token expires -> request to sso -> get new token
	// user.TokenTTL = time.Now().UTC().Add(time.Second*15)
	return user, nil
}

func validateUserData(user *model.User) error {
	if user.ID == 0 {return ErrInvalidToken}
	if !user.Valid {return ErrInvalidToken}
	// if user.TokenTTL.Unix() < time.Now().UTC().Unix() {return ErrInvalidToken}
	return nil
}

func (s *Server) RedisGetUser(token string) (*model.User, error) {
	s.logger.Info(fmt.Sprintf("user get token from cache: %s", token))
	data, err := s.cache.User().Get(token)

	if err != nil {
		fmt.Printf("%T", err)
		s.logger.Errorf("redis err: %v", err)
		return nil, err
	}

	var user model.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, fmt.Errorf("nil data")
	}

	return &user, nil
}

func (s *Server) RedisSetUser(token string, user *model.User) error{
	s.logger.Info(fmt.Sprintf("redis set %s, %v", token, user.TokenTTL))
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.cache.User().Set(token, data, user.TokenTTL)
}

func (s *Server) AuthMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s.logger.Info(r.URL)

			authHeader := r.Header.Get("Authorization")
			authHeaderChunks := strings.Split(authHeader, " ")

			if len(authHeaderChunks) != 2 {
				s.logger.Info("user is not authorized")
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserKey, 0)))
				return
			}

			token := authHeaderChunks[1]

			user, err := s.RedisGetUser(token)
			if err != nil {
				user, err = getUserInfoByToken(token)
				// user.TokenTTL = time.Now().UTC().Add(time.Duration(time.Second*15))
				if err != nil {
					s.logger.Errorf("get user info b token: %v", err)
					s.error(w, r, errorResponse{Error: ErrInternalServerError.Error(), Code: http.StatusInternalServerError})
					return
				}

				//validate user
				if err := validateUserData(user); err != nil {
					s.logger.Errorf("get user info b token: %v", err)
					s.error(w, r, errorResponse{Error: ErrInvalidToken.Error(), Code: http.StatusForbidden})
					return
				}
				
				//set redis
				if err := s.RedisSetUser(token, user); err != nil{
					s.logger.Errorf("redis: set user err: %v", err)
				}
			} 
			
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserKey, user.ID)))
		},
	)
}