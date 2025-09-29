package user_repository

import (
	"errors"
	"record-services/internal/models"
	"record-services/pkg/consts"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetById(id uint) (*models.User, error)
	GetByIdWithOutPassword(id uint) (*models.UserResponse, error)
	GetByEmail(email string) (*models.User, error)
	GetByEmailWithOutPassword(email string) (*models.UserResponse, error)
	Create(user *models.User) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id uint) error
	GetAll(limit, offset int, name string) ([]models.User, error)
	GetAllWithPagination(limit, page int, name string) (*models.PaginatedUsers, error)
}

type userRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
	cache  *cachedData
}

func NewUserRepository(db *gorm.DB, logger *zerolog.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
		cache:  newCachedData(5 * time.Minute),
	}
}

func (r *userRepository) GetById(id uint) (*models.User, error) {
	// Пытаемся получить из кеша
	user, exists := r.cache.getById(id)

	if exists {
		return user, nil
	}
	// Если нет в кеше, получаем из БД
	user = &models.User{}
	result := r.db.First(user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error().Err(result.Error).Msgf("ошибка при получении пользователя по id: %d", id)
		return nil, result.Error
	}
	// Сохраняем в кеш
	r.cache.add(user)

	return user, nil
}

func (r *userRepository) GetByIdWithOutPassword(id uint) (*models.UserResponse, error) {
	user, err := r.GetById(id)
	if err != nil {
		return nil, err
	}
	return &models.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
		IsAdmin:  user.IsAdmin,
	}, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	// Пытаемся получить из кеша
	user, exists := r.cache.getByEmail(email)

	if exists {
		return user, nil
	}
	
	user = &models.User{}
	result := r.db.First(user, "email = ?", email)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error().Err(result.Error).Msgf("ошибка при получении пользователя по email: %s", email)
		return nil, result.Error
	}

	// Добавляем в кеш
	r.cache.add(user)

	return user, nil
}

func (r *userRepository) GetByEmailWithOutPassword(email string) (*models.UserResponse, error) {
	user, err := r.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return &models.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
		IsAdmin:  user.IsAdmin,
	}, nil
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	result := r.db.Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, consts.ErrAlreadyExists
		}
		r.logger.Error().Err(result.Error).Msgf("ошибка при создании пользователя: %v", user)
		return nil, result.Error
	}
	
	// Добавляем в кеш
	r.cache.add(user)
	
	return user, nil
}

func (r *userRepository) Update(user *models.User) (*models.User, error) {
	result := r.db.Save(user)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msgf("ошибка при обновлении пользователя: %v", user)
		return nil, result.Error
	}
	r.cache.delete(user.ID)
	
	return user, nil
}

func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msgf("ошибка при удалении пользователя по id: %d", id)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return consts.ErrNotFound
	}

	r.cache.delete(id)

	return nil
}

func (r *userRepository) GetAll(limit, offset int, name string) ([]models.User, error) {
	var users []models.User

	query := r.db.Preload("User")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	
	result := query.Find(&users)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("ошибка при получении пользователей")
		return nil, result.Error
	}
	return users, nil
}

func (r *userRepository) GetAllWithPagination(limit, page int, name string) (*models.PaginatedUsers, error) {
	var users []models.User
	var totalCount int64

	countQuery := r.db.Model(&models.User{})

	dataQuery := r.db.Preload("User")

	// Применяем фильтры к обоим запросам
	if name != "" {
		countQuery = countQuery.Where("name LIKE ?", "%"+name+"%")
		dataQuery = dataQuery.Where("name LIKE ?", "%"+name+"%")
	}	

	// Получаем общее количество записей
	if err := countQuery.Count(&totalCount).Error; err != nil {
		r.logger.Error().Err(err).Msg("ошибка при подсчете пользователей")
		return nil, err
	}

	// Рассчитываем пагинацию
	if limit <= 0 {
		limit = 30
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	// Получаем данные с пагинацией
	result := dataQuery.
		Order("name").
		Limit(limit).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		r.logger.Error().Err(result.Error).Msg("ошибка при получении пользователей")
		return nil, result.Error
	}

	// Рассчитываем общее количество страниц
	totalPages := int(totalCount) / limit
	if int(totalCount)%limit > 0 {
		totalPages++
	}

	// Проверяем, есть ли следующая страница
	hasMore := page < totalPages

	usersResponse := make([]models.UserResponse, len(users))
	for i, user := range users {
		usersResponse[i] = models.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			IsActive: user.IsActive,
			IsAdmin:  user.IsAdmin,
		}
	}

	return &models.PaginatedUsers{
		Users:      usersResponse,
		TotalCount: totalCount,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
		HasMore:    hasMore,
	}, nil
}
