package response

import (
	"net/http"

	"github.com/go-chi/render"
)

type JWT struct {
	Status         string `json:"status,omitempty"`
	Message        string `json:"message"`
	Type           string `json:"type"`
	AccessToken    string `json:"access_token"`
	ExpirationTime string `json:"expiration_time"`
	RefreshToken   string `json:"refresh_token"`
}

func JsonJWT(w http.ResponseWriter, r *http.Request, code int, response *JWT) {
	render.Status(r, code)
	render.JSON(w, r, response)
}
