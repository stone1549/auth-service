package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/stone1549/auth-service/repository"
	"net/http"
)

type newSessionRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type newSessionResponse struct {
	Token string `json:"token"`
}

func (nsr newSessionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// NewSessionMiddleware middleware to authenticate a user from the request parameters
func NewSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var reqUser newSessionRequest
		err := decoder.Decode(&reqUser)
		if err != nil {
			render.Render(w, r, errInvalidRequest(err))
			return
		}

		if reqUser.Email == "" {
			render.Render(w, r, errInvalidRequest(errors.New("email is required")))
			return
		}

		if reqUser.Password == "" {
			render.Render(w, r, errInvalidRequest(errors.New("password is required")))
			return
		}

		userRepo, ok := r.Context().Value("repo").(repository.UserRepository)

		if !ok {
			render.Render(w, r, errRepository(errors.New("UserRepository not found in context")))
			return
		}

		id, err := userRepo.Authenticate(r.Context(), reqUser.Email, reqUser.Password)

		if err != nil {
			render.Render(w, r, errRepository(err))
			return
		}

		tokenFactory, ok := r.Context().Value("tokenFactory").(TokenFactory)

		if !ok {
			render.Render(w, r, errUnknown(errors.New("token factory not found in context")))
			return
		}

		token, err := tokenFactory.NewToken(NewClaims(id, reqUser.Email))

		if err != nil {
			render.Render(w, r, errUnknown(errors.New("unable to create token")))
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewSession responds to authentication request with jwt token or appropriate error
func NewSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token, ok := ctx.Value("token").(string)

	if !ok {
		render.Render(w, r, errUnknown(errors.New("unable to authenticate user")))
		return
	}

	if err := render.Render(w, r, newSessionResponse{token}); err != nil {
		render.Render(w, r, errUnknown(err))
		return
	}
}
