package database

import (
	"ksogit.kingsoft.net/kgo/mysql"
)

type KingbaseDialect struct {
	*BaseDialect
}

func NewKingbaseDialect(db mysql.DBAdapter) *KingbaseDialect {
	return &KingbaseDialect{BaseDialect: NewBaseDialect(db)}
}

func (m *KingbaseDialect) GetTableColumns(tableName string) ([]ColumnInfo, error) {
	schema, err := m.BaseDialect.GetTableSchema(tableName)
	if err != nil {
		return nil, err
	}
	return getColumnsFroPGLikeSchema(schema)
}
