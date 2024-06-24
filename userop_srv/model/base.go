package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type GormList []string

// Value implements the driver.Valuer interface for GormList
func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// Scan implements the sql.Scanner interface for GormList
func (g *GormList) Scan(value interface{}) error {
	if value == nil {
		*g = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, g)
	case string:
		return json.Unmarshal([]byte(v), g)
	default:
		return fmt.Errorf("cannot sql.Scan() unsupported type %T into GormList", value)
	}
}

type BaseModel struct {
	ID        int32          `gorm:"primarykey;type:int" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool           `json:"-"`
}
