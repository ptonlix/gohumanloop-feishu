package services

import (
	"time"

	"github.com/beego/beego/orm"
	"github.com/ptonlix/gohumanloop-wework/init/sqlite"
	"github.com/ptonlix/gohumanloop-wework/models"
)

type GoHumanloopService struct {
	ormer orm.Ormer
}

func NewGoHumanloopService() *GoHumanloopService {
	return &GoHumanloopService{
		ormer: sqlite.Mydb,
	}
}

var GHLDataService *GoHumanloopService = NewGoHumanloopService()

// CreateHumanLoop 创建人机协作记录
func (s *GoHumanloopService) CreateHumanLoop(hl *models.HumanLoop) (int64, error) {
	return s.ormer.Insert(hl)
}

// GetHumanLoopByID 通过ID获取人机协作记录
func (s *GoHumanloopService) GetHumanLoopByID(id int64) (*models.HumanLoop, error) {
	hl := &models.HumanLoop{ID: id}
	err := s.ormer.Read(hl)
	return hl, err
}

// GetHumanLoopByRequestID 通过ConversationId, RequestID获取人机协作记录
func (s *GoHumanloopService) GetHumanLoopByRequestId(conversationId, requestId, platform string) (*models.HumanLoop, error) {
	hl := &models.HumanLoop{RequestId: requestId, ConversationId: conversationId, Platform: platform}
	err := s.ormer.Read(hl, "ConversationId", "RequestId", "Platform")
	return hl, err
}

// UpdateHumanLoop 更新人机协作记录
func (s *GoHumanloopService) UpdateHumanLoop(hl *models.HumanLoop) (int64, error) {
	return s.ormer.Update(hl)
}

// DeleteHumanLoop 删除人机协作记录(软删除)
func (s *GoHumanloopService) DeleteHumanLoop(id int64) error {
	hl := &models.HumanLoop{
		ID:        id,
		IsDeleted: true,
		DeletedAt: time.Now(), // 设置当前时间为删除时间
	}
	_, err := s.ormer.Update(hl, "IsDeleted", "DeletedAt")
	return err
}

// ListHumanLoops 获取人机协作记录列表
func (s *GoHumanloopService) ListHumanLoops(filter map[string]interface{}, page, pageSize int) ([]*models.HumanLoop, int64, error) {
	qs := s.ormer.QueryTable(&models.HumanLoop{}).Filter("IsDeleted", false)

	for key, value := range filter {
		qs = qs.Filter(key, value)
	}

	total, _ := qs.Count()
	var hls []*models.HumanLoop
	_, err := qs.Limit(pageSize, (page-1)*pageSize).All(&hls)
	return hls, total, err
}

// GetHumanLoopsByConversationId 通过ConversationId获取所有人机协作记录
func (s *GoHumanloopService) GetHumanLoopsByConversationId(conversationId string, platform string) ([]*models.HumanLoop, error) {
	var hls []*models.HumanLoop
	_, err := s.ormer.QueryTable(&models.HumanLoop{}).
		Filter("ConversationId", conversationId).
		Filter("Platform", platform).
		Filter("IsDeleted", false).
		OrderBy("-Created"). // 通常按创建时间降序排列
		All(&hls)
	return hls, err
}

// BatchUpdateHumanLoops 批量更新人机协作记录（支持事务）
func (s *GoHumanloopService) BatchUpdateHumanLoops(hls []*models.HumanLoop) error {
	// 开始事务
	err := s.ormer.Begin()
	if err != nil {
		return err
	}

	// 确保在函数返回时处理事务
	defer func() {
		if p := recover(); p != nil {
			// 发生panic时回滚事务
			s.ormer.Rollback()
			panic(p) // 重新抛出panic
		}
	}()

	// 遍历更新每条记录
	for _, hl := range hls {
		rows, err := s.ormer.Update(hl)
		if err != nil {
			// 更新失败时回滚事务
			s.ormer.Rollback()
			return err
		}

		// 新增检查：如果更新行数为0，认为更新失败
		if rows == 0 {
			s.ormer.Rollback()
			return orm.ErrNoRows
		}
	}

	// 提交事务
	return s.ormer.Commit()
}
