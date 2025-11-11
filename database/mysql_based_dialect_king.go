package database

import "database/sql"

type MysqlBasedKingbaseDialect struct {
	*BaseMysqlBasedDialect
}

func NewMysqlBasedKingbaseDialect(db *sql.DB) *MysqlBasedKingbaseDialect {
	return &MysqlBasedKingbaseDialect{BaseMysqlBasedDialect: NewBaseMysqlBasedDialect(db)}
}

func (m *MysqlBasedKingbaseDialect) GetTableColumns(tableName string) ([]ColumnInfo, error) {
	schema, err := m.BaseMysqlBasedDialect.GetTableSchema(tableName)
	if err != nil {
		return nil, err
	}
	return getColumnsFroPGLikeSchema(schema)
}
