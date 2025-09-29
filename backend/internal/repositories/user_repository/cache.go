package user_repository

import (
    "record-services/internal/models"
    "sync"
    "time"
)

type cachedUser struct {
    user      *models.User
    expiresAt time.Time
}

type cachedData struct {
    usersById    map[uint]*cachedUser
    usersByEmail map[string]*cachedUser
    mu           sync.RWMutex
    ttl          time.Duration
}

func newCachedData(ttl time.Duration) *cachedData {
    c := &cachedData{
        usersById:    make(map[uint]*cachedUser),
        usersByEmail: make(map[string]*cachedUser),
        ttl:          ttl,
    }
    
    // Запускаем фоновую очистку
    go c.startCleanup()
    
    return c
}

// Запускает фоновую очистку устаревших записей
func (c *cachedData) startCleanup() {
    ticker := time.NewTicker(time.Minute) // Очищаем каждую минуту
    defer ticker.Stop()
    
    for range ticker.C {
        c.clearExpired()
    }
}

// Очищает только устаревшие записи без блокировки чтения
func (c *cachedData) clearExpired() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    now := time.Now()
    for id, item := range c.usersById {
        if now.After(item.expiresAt) {
            delete(c.usersById, id)
            delete(c.usersByEmail, item.user.Email)
        }
    }
}

func (c *cachedData) getById(id uint) (*models.User, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if item, ok := c.usersById[id]; ok {
        if time.Now().Before(item.expiresAt) {
            return item.user, true
        }
    }
    return nil, false
}

func (c *cachedData) getByEmail(email string) (*models.User, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if item, ok := c.usersByEmail[email]; ok {
        if time.Now().Before(item.expiresAt) {
            return item.user, true
        }
    }
    return nil, false
}

// Добавляет пользователя в кеш
func (c *cachedData) add(user *models.User) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    cached := &cachedUser{
        user:      user,
        expiresAt: time.Now().Add(c.ttl),
    }
    
    // Удаляем старые записи если они есть
    if old, ok := c.usersById[user.ID]; ok {
        delete(c.usersByEmail, old.user.Email)
    }
    
    c.usersById[user.ID] = cached
    c.usersByEmail[user.Email] = cached
}

// Удаляет пользователя из кеша
func (c *cachedData) delete(id uint) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if item, ok := c.usersById[id]; ok {
        delete(c.usersById, id)
        delete(c.usersByEmail, item.user.Email)
    }
}
