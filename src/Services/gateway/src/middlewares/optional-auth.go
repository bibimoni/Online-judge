package middlewares

import (
	// "bytes"
	"encoding/json"
	"fmt"

	// "io"
	"net/http"
	"strconv"
	"time"

	"github.com/bibimoni/Online-judge/gateway/src/common"
	"github.com/bibimoni/Online-judge/gateway/src/infrastructure/config"
)

func OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authURL := config.Load().Endpoints.Auth
		config.GetLogger().Debug().Msgf("auth url: %s", authURL)

		if r.Header.Get("Authorization") == "" {
			next.ServeHTTP(w, r)
			return
		}

		res, err := common.SendRequest[AuthResponseBody](r.Context(), common.APIRequest{
			Method:  "POST",
			URL:     fmt.Sprintf("%s/auth/validate", authURL),
			Timeout: 10 * time.Second,
			Headers: map[string]string{
				"Authorization": r.Header.Get("Authorization"),
				"Content-Type":  "application/json",
			}, Body: nil,
		})

		if err != nil {
			if res != nil {
				http.Error(w, "Auth Service Error", res.StatusCode)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// org, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	http.Error(w, "Can't read body", http.StatusBadRequest)
		// 	return
		// }
		// defer r.Body.Close()
		//
		// var obj map[string]any
		// if len(org) > 0 {
		// 	if err := json.Unmarshal(org, &obj); err != nil {
		// 		http.Error(w, "Invalid Json", http.StatusBadRequest)
		// 	}
		// } else {
		// 	obj = make(map[string]any)
		// }

		config.GetLogger().Debug().Msgf("Auth Response: %+v", res.PayLoad)

		r.Header.Set("X-User-Id", strconv.Itoa(res.PayLoad.Id))
		r.Header.Set("X-User-Role", res.PayLoad.Role)
		permsJson, _ := json.Marshal(res.PayLoad.Permissions)
		r.Header.Set("X-User-Permissions", string(permsJson))
		r.Header.Set("X-Username", res.PayLoad.Username)

		config.GetLogger().Debug().Msgf("%+v", r.Header)

		// // obj["username"] = res.PayLoad.Username
		// newBytes, _ := json.Marshal(obj)
		//
		// r.Body = io.NopCloser(bytes.NewReader(newBytes))
		//
		// r.ContentLength = int64(len(newBytes))
		//
		// r.Header.Set("Content-Length", strconv.Itoa(len(newBytes)))

		next.ServeHTTP(w, r)
	})
}
