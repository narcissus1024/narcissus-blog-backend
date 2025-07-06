package model

const TableNameRole = "roles"

// Role mapped from table <roles>
type Role struct {
	ID   int   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Role int32 `gorm:"column:role;not null;comment:0:普通用户;1:系统管理员" json:"role"` // 0:普通用户;1:系统管理员
}

// TableName Role's table name
func (*Role) TableName() string {
	return TableNameRole
}
