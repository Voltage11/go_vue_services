package middleware

import (
	"context"
	"net/http"
	"record-services/internal/auth"
	"record-services/pkg/consts"
	"record-services/pkg/utils"
	"strings"
)

// Тип запроса middleware браузер или api
type requestType string

const (
	requestTypeCookie requestType = "cookie"
	requestTypeBearer requestType = "bearer"
	requestTypeNone   requestType = "none"
)

var publicPath = map[string]bool{
	"/api/auth/login":    true,
	"/api/auth/register": true,
}

func isPublicPath(path string) bool {
	return publicPath[path]
}

func AuthMiddleware(authService *auth.AuthHandlers) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if isPublicPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			tokenString, reqType := extractToken(r)
			
			if tokenString == "" {
				sendUnauthorizedError(w, reqType)
				return
			}

			// Валидируем токен
			user, err := utils.ParseToken(tokenString, []byte(authService.JwtSecret))
			if err != nil {
				handleInvalidToken(w, reqType)
				return
			}

			// Добавляем пользователя в контекст
			ctx := context.WithValue(r.Context(), string(consts.ContextUserKey), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) (string, requestType) {
	// Проверяем cookie
	cookie, err := r.Cookie(string(consts.CookieTokenKey))
	if err == nil && cookie != nil && cookie.Value != "" {
		return cookie.Value, requestTypeCookie
	}

	// Проверяем Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token != "" {
				return token, requestTypeBearer
			}
		}
	}

	return "", requestTypeNone
}

func sendUnauthorizedError(w http.ResponseWriter, reqType requestType) {
	message := "Неавторизован"
	if reqType == requestTypeBearer {
		message = "Невалидный или отсутствует Bearer токен"
	}
	
	http.Error(w, message, http.StatusUnauthorized)
}

func handleInvalidToken(w http.ResponseWriter, reqType requestType) {
	// Если токен был из cookie - очищаем его
	if reqType == requestTypeCookie {
		clearAuthCookie(w)
	}
	
	http.Error(w, "Невалидный токен", http.StatusUnauthorized)
}

func clearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     string(consts.CookieTokenKey),
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode, // добавляем для consistency
	})

	http.SetCookie(w, &http.Cookie{
		Name:     string(consts.CookieUserDataKey),
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // добавляем для полного удаления
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func GetUserFromContext(ctx context.Context) *utils.UserClaims {
	if user, ok := ctx.Value(string(consts.ContextUserKey)).(*utils.UserClaims); ok {
		return user
	}
	return nil
}