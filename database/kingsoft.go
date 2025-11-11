package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"ksogit.kingsoft.net/chat/lib/xmysql"
	xmysqlv2 "ksogit.kingsoft.net/chat/lib/xmysql/v2"
	"ksogit.kingsoft.net/kgo/mysql"
)

type BaseKingsoftDB struct {
	db      mysql.DBAdapter
	dialect Dialect

	dbConfig    *xmysql.Database
	dialectType string
}

func NewKingsoftDB(dialectType string) *BaseKingsoftDB {
	return &BaseKingsoftDB{dialectType: dialectType}
}

// Connect 建立MySQL连接
func (m *BaseKingsoftDB) Connect(dsn string) error {
	dbConfig, err := xmysql.GetDatabaseFromDSN(dsn)
	if err != nil {
		return err
	}
	return m.ConnectWithConfig(dbConfig)
}

func formatDSN(dbConfig *xmysql.Database) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.UserName, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
}

func (m *BaseKingsoftDB) ConnectWithConfig(dbConfig *xmysql.Database) error {
	db, err := xmysqlv2.NewDBBuilder(dbConfig, &xmysqlv2.ServiceInfo{
		LocalEnv:    "local",
		DeployEnv:   "prod",
		ServiceName: "logic",
	}).WithNameSuffix("master").Build(context.Background())
	if err != nil {
		return err
	}
	if m.db != nil {
		m.db.Close()
	}
	m.db = db
	m.dialect = GetDialectByType(m.dialectType, db)
	m.dbConfig = dbConfig
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
	var rows = make([]map[string]interface{}, 0)
	sqlxDB, err := sqlx.Open("mysql", formatDSN(m.dbConfig))
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接失败: %w", err)
	}
	defer sqlxDB.Close()
	scanRows, err := sqlxDB.Queryx(query)
	if err != nil {
		return nil, fmt.Errorf("执行查询失败: %w", err)
	}
	for scanRows.Next() {
		row := make(map[string]interface{})
		err = scanRows.MapScan(row)
		if err != nil {
			return nil, fmt.Errorf("扫描数据失败: %w", err)
		}
		for k := range row {
			if value, ok := row[k].([]byte); ok {
				row[k] = string(value)
			}
		}
		rows = append(rows, row)
	}
	return rows, scanRows.Err()
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

func (m *BaseKingsoftDB) ExecuteUpdate(query string) (int64, error) {
	res, err := m.db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("执行更新失败: %w", err)
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
	var rows = make([]map[string]interface{}, 0)
	sqlxDB, err := sqlx.Open("mysql", formatDSN(m.dbConfig))
	if err != nil {
		return nil, 0, fmt.Errorf("打开数据库连接失败: %w", err)
	}
	defer sqlxDB.Close()
	scanRows, err := sqlxDB.Queryx(query)
	if err != nil {
		return nil, 0, fmt.Errorf("查询数据失败: %w", err)
	}
	for scanRows.Next() {
		row := make(map[string]interface{})
		err = scanRows.MapScan(row)
		if err != nil {
			return nil, 0, fmt.Errorf("扫描数据失败: %w", err)
		}
		for k := range row {
			if value, ok := row[k].([]byte); ok {
				row[k] = string(value)
			}
		}
		rows = append(rows, row)
	}
	return rows, total, nil
}

// GetDatabases 获取所有数据库名称
func (m *BaseKingsoftDB) GetDatabases() ([]string, error) {
	return m.dialect.GetDatabases()
}

// SwitchDatabase 切换当前使用的数据库
func (m *BaseKingsoftDB) SwitchDatabase(databaseName string) error {
	m.dbConfig.DBName = databaseName
	return m.ConnectWithConfig(m.dbConfig)
}
