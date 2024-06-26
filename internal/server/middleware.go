package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rautaruukkipalich/go_notes/internal/model"
)

func (s *Server) RedisGetUser(token string) (*model.User, error) {
	s.logger.Info(fmt.Sprintf("try user get token from cache: %s", token))
	data, err := s.cache.User().Get(token)

	if err != nil {
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
	s.logger.Info(fmt.Sprintf("try redis set %s, %v", token, user.TokenTTL))
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

			authHeaderChunks := strings.Split(r.Header.Get("Authorization"), " ")

			if len(authHeaderChunks) != 2 {
				s.logger.Info("user is not authorized")
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserKey, 0)))
				return
			}

			token := authHeaderChunks[1]

			user, err := s.RedisGetUser(token)
			if err != nil {
				user, err = getUserByToken(token)
				if err != nil {
					s.logger.Errorf("try get user info by token: %v", err)
					s.logger.Errorln("msg:", err, "user:", user.Str(), "token:", token)
					// next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserKey, user.ID)))
					// return
					// s.error(w, r, errorResponse{Error: ErrInternalServerError.Error(), Code: http.StatusInternalServerError})
				}

				//validate user
				if err := validateUserData(user); err != nil {
					s.logger.Errorf("try validate user: %v", err)
					// next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserKey, user.ID)))
					// s.error(w, r, errorResponse{Error: ErrInvalidToken.Error(), Code: http.StatusForbidden})
					// return
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