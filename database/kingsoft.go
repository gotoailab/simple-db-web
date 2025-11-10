package database

import (
	"context"
	"fmt"

	"ksogit.kingsoft.net/chat/lib/xmysql"
	xmysqlv2 "ksogit.kingsoft.net/chat/lib/xmysql/v2"
	"ksogit.kingsoft.net/kgo/mysql"
)

type BaseKingsoftDB struct {
	db      mysql.DBAdapter
	dialect Dialect
}

func NewKingsoftDB() *BaseKingsoftDB {
	return &BaseKingsoftDB{}
}

// Connect 建立MySQL连接
func (m *BaseKingsoftDB) Connect(dsn string) error {
	dbConfig, err := xmysql.GetDatabaseFromDSN(dsn)
	if err != nil {
		return err
	}
	db, err := xmysqlv2.NewDBBuilder(dbConfig, &xmysqlv2.ServiceInfo{}).WithNameSuffix("master").Build(context.Background())
	if err != nil {
		return err
	}
	m.db = db
	m.dialect = NewMySQLDialect(db)
	return nil
}

// Close 关闭连接
func (m *BaseKingsoftDB) Close() error {
	if m.db != nil {
		m.db.Close()
	}
	return nil
}

// GetTables 获取所有表名
func (m *BaseKingsoftDB) GetTables() ([]string, error) {
	return m.dialect.GetTables()
}

// GetTableSchema 获取表结构
func (m *BaseKingsoftDB) GetTableSchema(tableName string) (string, error) {
	return m.dialect.GetTableSchema(tableName)
}

// GetTableColumns 获取表的列信息
func (m *BaseKingsoftDB) GetTableColumns(tableName string) ([]ColumnInfo, error) {
	return m.dialect.GetTableColumns(tableName)
}

// ExecuteQuery 执行查询
func (m *BaseKingsoftDB) ExecuteQuery(query string) ([]map[string]interface{}, error) {
	var rows []map[string]interface{}
	if err := m.db.Query(&rows, query); err != nil {
		return nil, fmt.Errorf("执行查询失败: %w", err)
	}
	for i := range rows {
		for k, v := range rows[i] {
			if b, ok := v.([]byte); ok {
				rows[i][k] = string(b)
			}
		}
	}
	return rows, nil
}

// ExecuteUpdate 执行更新
func (m *BaseKingsoftDB) ExecuteUpdate(query string) (int64, error) {
	res, err := m.db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("执行更新失败: %w", err)
	}
	return res.RowsAffected, nil
}

// ExecuteDelete 执行删除
func (m *BaseKingsoftDB) ExecuteDelete(query string) (int64, error) {
	res, err := m.db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("执行删除失败: %w", err)
	}
	return res.RowsAffected, nil
}

// ExecuteInsert 执行插入
func (m *BaseKingsoftDB) ExecuteInsert(query string) (int64, error) {
	res, err := m.db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("执行插入失败: %w", err)
	}
	return res.RowsAffected, nil
}

// GetTableData 获取表数据（分页）
func (m *BaseKingsoftDB) GetTableData(tableName string, page, pageSize int) ([]map[string]interface{}, int64, error) {
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s`", tableName)
	total, err := m.db.QueryInt64(countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("查询总数失败: %w", err)
	}
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("SELECT * FROM `%s` LIMIT %d OFFSET %d", tableName, pageSize, offset)
	var rows []map[string]interface{}
	if err := m.db.Query(&rows, query); err != nil {
		return nil, 0, fmt.Errorf("查询数据失败: %w", err)
	}
	for i := range rows {
		for k, v := range rows[i] {
			if b, ok := v.([]byte); ok {
				rows[i][k] = string(b)
			} else if v == nil {
				rows[i][k] = nil
			}
		}
	}
	return rows, total, nil
}
