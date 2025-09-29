package auth

import (
	"encoding/json"
	"net/http"
	"record-services/internal/models"
	"record-services/internal/repositories/user_repository"
	"record-services/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type registerRequest struct {
	Name           string `json:"name" validate:"required,min=2,max=100"`
	Email          string `json:"email" validate:"required,email,max=255"`
	Password       string `json:"password" validate:"required,min=6,max=15"`
	PasswordRepeat string `json:"passwordRepeat" validate:"required,min=6,max=15"`
}

type loginRequest struct {
	Email    string `json:"email"  validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

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

	authHandlers.mux.HandleFunc("POST /auth/register", authHandlers.register)
	authHandlers.mux.HandleFunc("POST /auth/login", authHandlers.login)

	return authHandlers
}

func (h *AuthHandlers) register(w http.ResponseWriter, r *http.Request) {
	var registerData registerRequest
	err := json.NewDecoder(r.Body).Decode(&registerData)
	if err != nil {
		http.Error(w, "Невалидный запрос", http.StatusBadRequest)
		return
	}
	// проверим валидные ли данные
	err = h.validator.Struct(&registerData)
	if err != nil {
		http.Error(w, "Невалидный запрос", http.StatusBadRequest)
		return
	}

	if registerData.Password != registerData.PasswordRepeat {
		http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
		return
	}

	// Проверка наличие пользователя с таким email
	existUser, err := h.repository.GetByEmail(registerData.Email)
	if existUser != nil && err == nil {
		h.logger.Info().Msgf("Пользователь с email: %s уже существует", registerData.Email)
		http.Error(w, "Пользователь с таким email уже существует", http.StatusConflict)
		return
	}
	if err != nil {
		h.logger.Error().Err(err).Msgf("ошибка при получении пользователя по email: %s", registerData.Email)
		http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
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
		h.logger.Error().Err(err).Msgf("ошибка при регистрации пользователя: %s", registerData.Email)
		http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{
		"status": "ok"
	}`))

}

func (h *AuthHandlers) login(w http.ResponseWriter, r *http.Request) {

	var loginRequest loginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Не валидное тело запроса", http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(&loginRequest)
	if err != nil {
		http.Error(w, "Не валидное тело запроса", http.StatusBadRequest)
		return
	}
	user, err := h.repository.GetByEmail(loginRequest.Email)
	if err != nil {
		http.Error(w, "Ошибка при получении пользователя", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	if !utils.VerifyHash(loginRequest.Password, user.PasswordHash, h.hashSecret) {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	if !user.IsActive {
		http.Error(w, "Пользователь не активный, обратитесь к администратору", http.StatusUnauthorized)
		return
	}

	// Создаем токен для передачи в cookie
	token, err := utils.CreateToken(utils.UserClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, []byte(h.hashSecret))

	if err != nil {
		h.logger.Error().Err(err).Msgf("ошибка при создании токена для пользователя: %s", user.Email)
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	// Установим cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
