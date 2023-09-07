package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zura-t/go_delivery_system/internal"
	"github.com/zura-t/go_delivery_system/token"
)

type renewAccessTokenRequest struct {
	RefreshToken string
}

type renewAccessTokenResponse struct {
	AccessToken          string
	AccessTokenExpiresAt time.Time
}

func RenewAccessTokenHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			cookies := r.Cookies()[0]
			config, err := internal.LoadConfig("../..")
			tokenMaker, err := token.NewJwtMaker(config.TokenSymmetricKey)
			refreshPayload, err := tokenMaker.VerifyToken(cookies.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
			}

			accessToken, accessPayload, err := tokenMaker.CreateToken(refreshPayload.UserId, refreshPayload.Email, config.AccessTokenDuration)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
			}

			rsp := renewAccessTokenResponse{
				AccessToken:          accessToken,
				AccessTokenExpiresAt: accessPayload.ExpiredAt,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(rsp)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
