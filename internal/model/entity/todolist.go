package entity

type Todos struct {
	Id          int64  `gorm:"primaryKey;column:todos_id" json:"id"`
	Title       string `gorm:"type:varchar(99)" json:"title"`
	Description string `gorm:"type:varchar(999)" json:"description"`
	Status      bool   `gorm:"default:false" json:"status"`
}
