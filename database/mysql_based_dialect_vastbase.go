package database

import "database/sql"

type MysqlBasedVastbaseDialect struct {
	*BaseMysqlBasedDialect
}

func NewMysqlBasedVastbaseDialect(db *sql.DB) *MysqlBasedVastbaseDialect {
	return &MysqlBasedVastbaseDialect{BaseMysqlBasedDialect: NewBaseMysqlBasedDialect(db)}
}

func (m *MysqlBasedVastbaseDialect) GetTableColumns(tableName string) ([]ColumnInfo, error) {
	schema, err := m.BaseMysqlBasedDialect.GetTableSchema(tableName)
	if err != nil {
		return nil, err
	}
	return getColumnsFroPGLikeSchema(schema)
}
