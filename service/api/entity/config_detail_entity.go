package entity

type ConfigDetailEntity struct {
	BaseEntity
	Content string       `gorm:"column:content;type:varchar(20);not null;comment: 配置详情"`
	Cid     string       `gorm:"column:cId;type:varchar(20);not null;comment: 配置ID"`
	Config  ConfigEntity `gorm:"foreignKey:Cid;references:ID" json:"-"`
}

func (m *ConfigDetailEntity) TableName() string {
	return "config_detail"
}
