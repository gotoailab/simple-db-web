package database

import (
	"ksogit.kingsoft.net/kgo/mysql"
)

type VastbaseDialect struct {
	*BaseDialect
}

func NewVastbaseDialect(db mysql.DBAdapter) *VastbaseDialect {
	return &VastbaseDialect{BaseDialect: NewBaseDialect(db)}
}

func (m *VastbaseDialect) GetTableColumns(tableName string) ([]ColumnInfo, error) {
	return m.BaseDialect.GetTableColumns(tableName)
}
