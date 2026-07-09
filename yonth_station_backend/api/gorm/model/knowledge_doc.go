package model

import "time"

// KnowledgeDoc 知识库文档表，用于存储 RAG 系统的知识库文档
type KnowledgeDoc struct {
	// Id 文档唯一标识，自增主键
	Id int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// Title 文档标题，用于展示和检索
	Title string `gorm:"column:title;type:varchar(256);not null" json:"title"`

	// Content 文档正文内容，将被向量化后存入 Milvus
	Content string `gorm:"column:content;type:text;not null" json:"content"`

	// Category 文档分类，便于管理端筛选和分类检索
	Category string `gorm:"column:category;type:varchar(64);default:''" json:"category"`

	// Status 文档索引状态
	// 0 - 待索引（已保存但尚未向量化）
	// 1 - 已索引（已成功向量化并存入 Milvus）
	Status int8 `gorm:"column:status;type:tinyint;default:0" json:"status"`

	// CreatedAt 文档创建时间（自动填充）
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`

	// UpdatedAt 文档最后更新时间（自动更新）
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName 返回数据库表名
func (KnowledgeDoc) TableName() string {
	return "knowledge_doc"
}
