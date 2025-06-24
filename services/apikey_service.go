package services

import (
	"sync"
	"time"

	"github.com/beego/beego/orm"
	"github.com/google/uuid"
	"github.com/ptonlix/gohumanloop-wework/init/sqlite"
	"github.com/ptonlix/gohumanloop-wework/models"
)

type APIKeyService struct {
	ormer orm.Ormer
	cache sync.Map // 线程安全的缓存
}

func NewAPIKeyService() *APIKeyService {
	return &APIKeyService{
		ormer: sqlite.Mydb,
	}
}

var APIKeyDataService *APIKeyService = NewAPIKeyService()

// CreateAPIKey 创建新的APIKey
func (s *APIKeyService) CreateAPIKey(name string) (*models.APIKey, error) {
	apiKey := &models.APIKey{
		Name:    name,
		Key:     uuid.NewString(), // 使用UUID生成唯一的APIKey
		Status:  true,
		Created: time.Now(),
		Updated: time.Now(),
	}

	_, err := s.ormer.Insert(apiKey)
	return apiKey, err
}

// GetAPIKeyByID 根据ID获取APIKey
func (s *APIKeyService) GetAPIKeyByID(id int64) (*models.APIKey, error) {
	apiKey := &models.APIKey{ID: id}
	err := s.ormer.Read(apiKey)
	if err != nil {
		return nil, err
	}
	// 检查是否已删除
	if apiKey.IsDeleted {
		return nil, orm.ErrNoRows
	}
	return apiKey, nil
}

// GetAPIKeyByKey 根据Key获取APIKey
func (s *APIKeyService) GetAPIKeyByKey(key string) (*models.APIKey, error) {
	// 先从缓存读取
	if val, ok := s.cache.Load(key); ok {
		apiKey := val.(*models.APIKey)
		// 检查缓存中的key是否已删除
		if apiKey.IsDeleted {
			s.cache.Delete(key) // 删除已失效的缓存
			return nil, orm.ErrNoRows
		}
		return apiKey, nil
	}

	// 缓存未命中则查询数据库
	apiKey := &models.APIKey{Key: key}
	err := s.ormer.Read(apiKey, "Key")
	if err != nil {
		return nil, err
	}

	// 检查是否已删除
	if apiKey.IsDeleted {
		return nil, orm.ErrNoRows
	}

	// 存入缓存
	s.cache.Store(key, apiKey)
	return apiKey, nil
}

// UpdateAPIKey 更新API Key时同时更新缓存
func (s *APIKeyService) UpdateAPIKey(apiKey *models.APIKey) error {
	apiKey.Updated = time.Now()
	_, err := s.ormer.Update(apiKey)
	if err == nil {
		s.cache.Store(apiKey.Key, apiKey)
	}
	return err
}

// DeleteAPIKey 删除API Key时清除缓存
func (s *APIKeyService) DeleteAPIKey(id int64) error {
	apiKey := &models.APIKey{ID: id}
	err := s.ormer.Read(apiKey)
	if err != nil {
		return err
	}

	s.cache.Delete(apiKey.Key)
	apiKey.DeletedAt = time.Now()
	apiKey.IsDeleted = true
	_, err = s.ormer.Update(apiKey, "DeletedAt", "IsDeleted")
	return err
}

// ListAPIKeys 获取所有APIKey列表
func (s *APIKeyService) ListAPIKeys() ([]*models.APIKey, error) {
	var apiKeys []*models.APIKey
	_, err := s.ormer.QueryTable(new(models.APIKey)).
		Filter("is_deleted", false).
		All(&apiKeys)
	return apiKeys, err
}

// EnableAPIKey 启用APIKey
func (s *APIKeyService) EnableAPIKey(id int64) error {
	apiKey := &models.APIKey{ID: id, Status: true}
	_, err := s.ormer.Update(apiKey, "Status")
	return err
}

// DisableAPIKey 禁用APIKey
func (s *APIKeyService) DisableAPIKey(id int64) error {
	apiKey := &models.APIKey{ID: id, Status: false}
	_, err := s.ormer.Update(apiKey, "Status")
	return err
}
