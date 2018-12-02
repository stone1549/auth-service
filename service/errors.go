package service

import (
	"github.com/go-chi/render"
	"net/http"
)

type errResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func errInvalidRequest(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func errRepository(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Unable to handle request at this time.",
		ErrorText:      err.Error(),
	}
}

func errUnknown(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Unknown server error.",
		ErrorText:      err.Error(),
	}
}

var errNotFound = &errResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
