package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"record-services/internal/models"
	"record-services/internal/repositories/user_repository"
	"record-services/pkg/consts"
	"record-services/pkg/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

// Константы для кодов ошибок и сообщений
const (
	tokenExpiration = 24 * time.Hour
	cookieMaxAge    = 3600
)

type AuthHandlers struct {
	mux        *http.ServeMux
	logger     *zerolog.Logger
	repository user_repository.UserRepository
	validator  *validator.Validate
	hashSecret string
	JwtSecret  string
}

func NewAuthHandlers(mux *http.ServeMux, logger *zerolog.Logger, repository user_repository.UserRepository, validator *validator.Validate, hashSecret string, jwtSecret string) *AuthHandlers {
	authHandlers := &AuthHandlers{
		mux:        mux,
		logger:     logger,
		repository: repository,
		validator:  validator,
		hashSecret: hashSecret,
		JwtSecret:  jwtSecret,
	}

	authHandlers.mux.HandleFunc("POST /api/auth/register", authHandlers.register)
	authHandlers.mux.HandleFunc("POST /api/auth/login", authHandlers.login)
	authHandlers.mux.HandleFunc("POST /api/auth/logout", authHandlers.logout)

	return authHandlers
}

func (h *AuthHandlers) register(w http.ResponseWriter, r *http.Request) {
	var registerData struct {
		Name           string `json:"name" validate:"required,min=2,max=100"`
		Email          string `json:"email" validate:"required,email,max=255"`
		Password       string `json:"password" validate:"required,min=6,max=15"`
		PasswordRepeat string `json:"passwordRepeat" validate:"required,min=6,max=15"`
	}

	if !h.decodeAndValidate(w, r, &registerData) {
		return
	}

	if registerData.Password != registerData.PasswordRepeat {
		h.sendError(w, "Пароли не совпадают", http.StatusBadRequest)
		return
	}

	// Проверка наличия пользователя
	existUser, err := h.repository.GetByEmail(registerData.Email)
	if err != nil {
		h.logger.Error().Err(err).Msgf("Ошибка при проверке пользователя: %s", registerData.Email)
		h.sendError(w, "Ошибка при регистрации", http.StatusInternalServerError)
		return
	}

	if existUser != nil {
		h.logger.Info().Msgf("Пользователь с email: %s уже существует", registerData.Email)
		h.sendError(w, "Пользователь с таким email уже существует", http.StatusConflict)
		return
	}

	newUser := &models.User{
		Name:         registerData.Name,
		Email:        registerData.Email,
		PasswordHash: utils.CreateHash(registerData.Password, h.hashSecret),
		IsActive:     false,
		IsAdmin:      false,
	}

	if _, err := h.repository.Create(newUser); err != nil {
		h.logger.Error().Err(err).Msgf("Ошибка при создании пользователя: %s", registerData.Email)
		h.sendError(w, "Ошибка при регистрации", http.StatusInternalServerError)
		return
	}

	h.sendJSONResponse(w, map[string]string{"status": "ok"})
}

func (h *AuthHandlers) login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email" validate:"required,email,max=255"`
		Password string `json:"password" validate:"required,min=6,max=100"`
	}

	if !h.decodeAndValidate(w, r, &loginData) {
		return
	}

	user, err := h.repository.GetByEmail(loginData.Email)
	if err != nil {
		h.logger.Error().Err(err).Msgf("Ошибка при получении пользователя: %s", loginData.Email)
		h.sendError(w, "Ошибка при авторизации", http.StatusInternalServerError)
		return
	}

	if user == nil || !utils.VerifyHash(loginData.Password, user.PasswordHash, h.hashSecret) {
		h.sendError(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	if !user.IsActive {
		h.sendError(w, "Пользователь не активный, обратитесь к администратору", http.StatusUnauthorized)
		return
	}

	// Создаем токен
	token, err := utils.CreateToken(utils.UserClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, []byte(h.JwtSecret))

	if err != nil {
		h.logger.Error().Err(err).Msgf("Ошибка при создании токена: %s", user.Email)
		h.sendError(w, "Ошибка при авторизации", http.StatusInternalServerError)
		return
	}

	// Устанавливаем cookies
	h.setAuthCookies(w, token, user)

	h.sendJSONResponse(w, map[string]string{
		"status": "ok",
		"token":  token,
		"user":   user.Name, // Добавляем информацию о пользователе
	})
}

func (h *AuthHandlers) logout(w http.ResponseWriter, r *http.Request) {
	h.clearAuthCookies(w)
	h.sendJSONResponse(w, map[string]string{"status": "ok"})
}

// Вспомогательные методы

func (h *AuthHandlers) decodeAndValidate(w http.ResponseWriter, r *http.Request, data interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		h.sendError(w, "Невалидный запрос", http.StatusBadRequest)
		return false
	}

	if err := h.validator.Struct(data); err != nil {
		h.sendError(w, "Невалидные данные", http.StatusBadRequest)
		return false
	}

	return true
}

func (h *AuthHandlers) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   http.StatusText(statusCode),
		"message": message,
	})
}

func (h *AuthHandlers) sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *AuthHandlers) setAuthCookies(w http.ResponseWriter, token string, user *models.User) {
    // Токен cookie (без изменений)
    http.SetCookie(w, &http.Cookie{
        Name:     string(consts.CookieTokenKey),
        Value:    token,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
        MaxAge:   cookieMaxAge,
    })

    // User data cookie - используем base64
    userData := map[string]string{
        "name":  user.Name,
        "email": user.Email,
    }
    
    userDataJSON, err := json.Marshal(userData)
    if err != nil {
        h.logger.Error().Err(err).Msg("Ошибка при сериализации данных пользователя")
        return
    }

    // Кодируем в base64 чтобы избежать проблем с кавычками
    encodedUserData := base64.URLEncoding.EncodeToString(userDataJSON)

    http.SetCookie(w, &http.Cookie{
        Name:     string(consts.CookieUserDataKey),
        Value:    encodedUserData, // Безопасное значение
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
        MaxAge:   cookieMaxAge,
    })
}

func (h *AuthHandlers) clearAuthCookies(w http.ResponseWriter) {
	cookies := []string{
		string(consts.CookieTokenKey),
		string(consts.CookieUserDataKey),
	}

	for _, name := range cookies {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
		})
	}
}
