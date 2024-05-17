package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type GormList []string

func (g *GormList) Value() (driver.Value, error) {
	if g == nil {
		return "[]", nil
	}
	return json.Marshal(*g)
}

func (g *GormList) Scan(value interface{}) error {
	if value == nil {
		*g = GormList{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, g)
	case string:
		return json.Unmarshal([]byte(v), g)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

type BaseModel struct {
	ID        int32          `gorm:"primarykey;type:int" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool           `json:"-"`
}
