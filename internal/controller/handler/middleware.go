package handler

import (
	"context"
	"distributed-configuration/pkg/utils"
	"net/http"
	"strings"
)

func (h handler) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		switch token {
		case h.cfg.AdminSecret:
			ctx := context.WithValue(r.Context(), "role", utils.RoleAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		case h.cfg.ControllerSecret:
			ctx := context.WithValue(r.Context(), "role", utils.RoleAgent)

			if r.URL.Path == "/config" {
				agentID := r.Header.Get("X-Agent-ID")
				if agentID == "" {
					http.Error(w, "missing agent id", http.StatusUnauthorized)
					return
				}

				err := h.agent.Verify(ctx, agentID)
				if err != nil {
					status, msg := utils.MapError(err)
					http.Error(w, msg, status)
					return
				}
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		default:
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		}
	})
}

func (h handler) RoleBase(allowed ...utils.Role) func(http.Handler) http.Handler {
	allowedMap := make(map[utils.Role]struct{})
	for _, r := range allowed {
		allowedMap[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value("role").(utils.Role)
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			if _, exists := allowedMap[role]; !exists {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
