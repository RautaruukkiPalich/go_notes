package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")

type UserInfo struct{
	ID int `json:"id"`
	TokenTTL time.Time `json:"token_ttl"`
	Valid bool `json:"is_valid"`
}


var mockUserInfo = &UserInfo{
	ID: 1,
	TokenTTL: time.Time(time.Now().UTC().Add(time.Second*15)),
	Valid: true,
}

func getUserInfoByToken(token string) (*UserInfo, error) {
	user := mockUserInfo
	user.TokenTTL = time.Now().UTC().Add(time.Second*15)
	return user, nil
}

func validateUserData(user *UserInfo) error {
	if user.ID == 0 {return ErrInvalidToken}
	if !user.Valid {return ErrInvalidToken}
	if user.TokenTTL.Unix() < time.Now().UTC().Unix() {return ErrInvalidToken}
	return nil
}

func (s *Server) RedisGetUser(token string) (*UserInfo, error) {
	s.logger.Info(fmt.Sprintf("redis get %s", token))
	data, err := s.cache.User().Get(token)

	if err != nil {
		return nil, err
	}

	var user UserInfo
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, fmt.Errorf("nil data")
	}

	return &user, nil
}

func (s *Server) RedisSetUser(token string, user *UserInfo) error{
	s.logger.Info(fmt.Sprintf("redis set %s", token))
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.cache.User().Set(token, data, user.TokenTTL)
}

func (s *Server) AuthMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			authHeaderChunks := strings.Split(authHeader, " ")

			if len(authHeaderChunks) != 2 {
				// s.error(w, r, errorResponse{Error: ErrInvalidToken.Error(), Code: http.StatusForbidden})
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
					s.logger.Error(fmt.Sprintf("get user info b token: %v", err.Error()))
					s.error(w, r, errorResponse{Error: ErrInternalServerError.Error(), Code: http.StatusInternalServerError})
					return
				}
				//validate user
				if err := validateUserData(user); err != nil {
					s.logger.Error(fmt.Sprintf("get user info b token: %v", err.Error()))
					s.error(w, r, errorResponse{Error: ErrInvalidToken.Error(), Code: http.StatusForbidden})
					return
				}
				//set redis
				if err := s.RedisSetUser(token, user); err != nil{
					s.logger.Error(fmt.Sprintf("get user info b token: %v", err.Error()))
					// s.error(w, r, errorResponse{Error: ErrInvalidToken.Error(), Code: http.StatusForbidden})
					// return
				}
			} 
			
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserKey, user.ID)))
		},
	)
}