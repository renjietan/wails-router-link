package entity

// uniqueIndex;
type ConfigEntity struct {
	BaseEntity
	// 配置名称
	Name string `gorm:"unique;column:name;type:varchar(20);not null;comment:配置名称"`
	// 配置的值
	Value   string               `gorm:"column:value;type:text;not null"`
	Details []ConfigDetailEntity `gorm:"foreignKey:Cid;references:ID" json:"details,omitempty"`
}

func (m *ConfigEntity) TableName() string {
	return "config"
}
