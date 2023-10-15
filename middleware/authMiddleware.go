package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Maxcarrassco/blog_aggregator/handler"
)


func AuthMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		if authz != "" {	
			auth := strings.Split(authz, " ")
			if len(auth) == 2 || auth[0] == "ApiKey" || auth[1] != "" {
				userParam, err := handler.ApiCfg.DB.GetUserByApiKey(handler.ApiCfg.Ctx, auth[1])
				if err == nil {
					ctx := context.WithValue(r.Context(), "user", userParam)
					r = r.WithContext(ctx)
				}
			}
		}

		next.ServeHTTP(w, r);
	})
}
