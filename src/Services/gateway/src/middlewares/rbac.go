package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/bibimoni/Online-judge/gateway/src/common"
	"github.com/bibimoni/Online-judge/gateway/src/infrastructure/config"
)

type VerifyRequestBody struct {
	Token      string `json:"token"`
	Permission string `json:"permission"`
}

type VerifyResponseBody struct {
	Allowed bool        `json:"allowed"`
	User    UserPayload `json:"user"`
}

type UserPayload struct {
	Id          int      `json:"id"`
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

func WithPermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authURL := config.Load().Endpoints.Auth
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "No token provided", http.StatusUnauthorized)
				return
			}

			token := extractBearer(authHeader)
			if token == "" {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			reqBody := VerifyRequestBody{
				Token:      token,
				Permission: permission,
			}

			jsonBody, _ := json.Marshal(reqBody)

			res, err := common.SendRequest[VerifyResponseBody](r.Context(), common.APIRequest{
				Method:  "POST",
				URL:     fmt.Sprintf("%s/auth/verify", authURL),
				Timeout: 10 * time.Second,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: bytes.NewBuffer(jsonBody),
			})

			if err != nil {
				if res != nil {
					http.Error(w, "Auth Service Error", res.StatusCode)
					return
				}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !res.PayLoad.Allowed {
				http.Error(w, "Forbidden: Insufficient Permissions", http.StatusForbidden)
				return
			}

			r.Header.Set("X-User-Id", strconv.Itoa(res.PayLoad.User.Id))
			r.Header.Set("X-User-Role", res.PayLoad.User.Role)

			permsJson, _ := json.Marshal(res.PayLoad.User.Permissions)
			r.Header.Set("X-User-Permissions", string(permsJson))

			org, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Can't read body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			var obj map[string]any
			if len(org) > 0 {
				if err := json.Unmarshal(org, &obj); err != nil {
					http.Error(w, "Invalid Json", http.StatusBadRequest)
				}
			} else {
				obj = make(map[string]any)
			}

			obj["username"] = res.PayLoad.User.Username
			newBodyBytes, _ := json.Marshal(obj)

			r.Body = io.NopCloser(bytes.NewBuffer(newBodyBytes))
			r.ContentLength = int64(len(newBodyBytes))
			r.Header.Set("Content-Length", strconv.Itoa(len(newBodyBytes)))

			next.ServeHTTP(w, r)
		})
	}
}
