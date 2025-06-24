package models

import (
	"encoding/json"
	"time"

	"github.com/beego/beego/v2/client/orm"
	_ "modernc.org/sqlite"
)

func init() {
	orm.RegisterModel(new(HumanLoop))
	orm.RegisterDriver("sqlite3", orm.DRSqlite)
}

// HumanLoop 人机协作请求主模型
type HumanLoop struct {
	ID             int64     `orm:"auto;pk;column(id)" json:"-"`
	TaskId         string    `orm:"size(128);index" json:"task_id"`
	ConversationId string    `orm:"size(128);index" json:"conversation_id"`
	RequestId      string    `orm:"size(128);unique" json:"request_id"`
	SpNo           string    `orm:"size(128);index" json:"sp_no"`
	LoopType       string    `orm:"size(64)" json:"loop_type"`
	Context        string    `orm:"type(text)" json:"-"` // 实际存储
	Platform       string    `orm:"size(64);index" json:"platform"`
	Metadata       string    `orm:"type(text);null" json:"-"` // 实际存储
	Status         string    `orm:"size(32);default(pending)" json:"status"`
	Response       string    `orm:"type(text);null" json:"response,omitempty"` // 实际存储
	Feedback       string    `orm:"type(text);null" json:"feedback,omitempty"` // 实际存储
	RespondedBy    string    `orm:"size(128);null" json:"responded_by,omitempty"`
	RespondedAt    time.Time `orm:"null;type(datetime)" json:"responded_at,omitempty"`
	Created        time.Time `orm:"auto_now_add;type(datetime)" json:"-"`
	Updated        time.Time `orm:"auto_now;type(datetime)" json:"-"`
	DeletedAt      time.Time `orm:"null;type(datetime)" json:"-"` // 软删除字段
	IsDeleted      bool      `orm:"default(false)" json:"-"`      // 是否已删除标记

	// 非数据库字段，用于JSON序列化
	ContextMap  map[string]any `orm:"-" json:"context"`
	MetadataMap map[string]any `orm:"-" json:"metadata,omitempty"`
}

// 字符串转map
func (h *HumanLoop) AfterLoad() {
	// 反序列化JSON字段
	json.Unmarshal([]byte(h.Context), &h.ContextMap)
	json.Unmarshal([]byte(h.Metadata), &h.MetadataMap)
}

// Map转字符串存数据库
func (h *HumanLoop) BeforeSave() {
	// 序列化JSON字段
	if ctx, err := json.Marshal(h.ContextMap); err == nil {
		h.Context = string(ctx)
	}
	if meta, err := json.Marshal(h.MetadataMap); err == nil {
		h.Metadata = string(meta)
	}
}

// TableName 设置表名
func (h *HumanLoop) TableName() string {
	return "human_loops"
}

type APIKey struct {
	ID        int64     `orm:"auto;pk;column(id)" json:"id"`
	Name      string    `orm:"size(128);unique" json:"name"`
	Key       string    `orm:"size(128);unique" json:"key"`
	Status    bool      `orm:"default(true)" json:"status"`
	Created   time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
	Updated   time.Time `orm:"auto_now;type(datetime)" json:"updated_at"`
	DeletedAt time.Time `orm:"null;type(datetime)" json:"-"` // 软删除字段
	IsDeleted bool      `orm:"default(false)" json:"-"`      // 是否已删除标记
}

// TableName 设置表名
func (k *APIKey) TableName() string {
	return "api_keys"
}
