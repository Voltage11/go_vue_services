package middleware

import (
	"context"
	"net/http"
	"record-services/internal/auth"
	"record-services/pkg/consts"
	"record-services/pkg/utils"
)

var publicPath = map[string]bool{
	"/api/auth/login": true,
	"/api/auth/register": true,
}

func isPublicPath(path string) bool {
	if _, ok := publicPath[path]; ok {
		return true
	}
	return false
}

func AuthMiddleware(authService *auth.AuthHandlers) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем JWT токен из cookie
			cookie, _ := r.Cookie(string(consts.TokenCookieKey))
			// Если нет cookie, то проверяем на публичные пути
			if cookie.Value == "" {
				if isPublicPath(r.URL.Path) {
					next.ServeHTTP(w, r)
					return
				} else {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}

			tokenString := cookie.Value

			// Валидируем токен
			user, err := utils.ParseToken(tokenString, []byte(authService.JwtSecret))
			if err != nil {
				// Если токен невалиден, удаляем испорченную cookie
				http.SetCookie(w, &http.Cookie{
					Name:   string(consts.TokenCookieKey),
					Value:  "",
					Path:   "/",
					MaxAge: -1, // удаляем cookie
				})
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Добавляем пользователя в контекст
			ctx := context.WithValue(r.Context(), string(consts.TokenCookieKey), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Вспомогательная функция для получения пользователя из контекста
// func GetUserFromContext(ctx context.Context) *auth.User {
// 	if user, ok := ctx.Value(consts.TokenCookieKey).(*auth.User); ok {
// 		return user
// 	}
// 	return nil
// }