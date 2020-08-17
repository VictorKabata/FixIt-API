package models

type Work struct {
	ID       uint32 `gorm:"primary_key;auto_increment" json:"id"`
	UserID   uint32 `gorm:"not null" json:"user_id"`
	WorkerID uint32 `gorm:"not null" json:"user_id"`
	PostID   uint32 `gorm:"not null" json:"post_id"`
	Status   string `gorm:"not null" json:"status"`
}
