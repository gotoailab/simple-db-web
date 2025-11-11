package database

import (
	"database/sql"
	"fmt"
	"strings"

	"ksogit.kingsoft.net/kgo/mysql"
)

func GetDialectByType(dbType string, db mysql.DBAdapter) Dialect {
	switch dbType {
	case "dameng":
		return NewDamengDialect(db)
	case "openguass":
		return NewOpenguassDialect(db)
	case "vastbase":
		return NewVastbaseDialect(db)
	case "kingbase":
		return NewKingbaseDialect(db)
	case "oceandb":
		return NewOceandbDialect(db)
	}
	return nil
}

type Dialect interface {
	GetDatabases() ([]string, error)
	GetTables() ([]string, error)
	GetTableSchema(tableName string) (string, error)
	GetTableColumns(tableName string) ([]ColumnInfo, error)
}

type BaseDialect struct {
	db mysql.DBAdapter
}

func NewBaseDialect(db mysql.DBAdapter) *BaseDialect {
	return &BaseDialect{db: db}
}

func (m *BaseDialect) GetDatabases() ([]string, error) {
	var databases []string
	if err := m.db.Query(&databases, "SHOW DATABASES"); err != nil {
		return nil, fmt.Errorf("查询数据库列表失败: %w", err)
	}
	return databases, nil
}

// GetTables 获取所有表名
func (m *BaseDialect) GetTables() ([]string, error) {
	var tables []string
	if err := m.db.Query(&tables, "SHOW TABLES"); err != nil {
		return nil, fmt.Errorf("查询表列表失败: %w", err)
	}
	return tables, nil
}

// GetTableSchema 获取表结构
func (m *BaseDialect) GetTableSchema(tableName string) (string, error) {
	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
	var schemas []string
	if err := m.db.Query(&schemas, query); err != nil {
		return "", fmt.Errorf("查询表结构失败: %w", err)
	}
	if len(schemas) == 0 {
		return "", fmt.Errorf("表 %s 不存在", tableName)
	}
	return schemas[0], nil
}

// GetTableColumns 获取表的列信息
func (m *BaseDialect) GetTableColumns(tableName string) ([]ColumnInfo, error) {
	query := fmt.Sprintf("DESCRIBE `%s`", tableName)
	var rows []map[string]interface{}
	if err := m.db.Query(&rows, query); err != nil {
		return nil, fmt.Errorf("查询列信息失败: %w", err)
	}
	var columns []ColumnInfo
	for _, r := range rows {
		col := ColumnInfo{}
		if v, ok := r["Field"]; ok {
			if b, ok2 := v.([]byte); ok2 {
				col.Name = string(b)
			} else {
				col.Name = fmt.Sprint(v)
			}
		}
		if v, ok := r["Type"]; ok {
			if b, ok2 := v.([]byte); ok2 {
				col.Type = string(b)
			} else {
				col.Type = fmt.Sprint(v)
			}
		}
		if v, ok := r["Null"]; ok {
			var s string
			if b, ok2 := v.([]byte); ok2 {
				s = string(b)
			} else {
				s = fmt.Sprint(v)
			}
			col.Nullable = (strings.EqualFold(s, "YES"))
		}
		if v, ok := r["Key"]; ok {
			if b, ok2 := v.([]byte); ok2 {
				col.Key = string(b)
			} else {
				col.Key = fmt.Sprint(v)
			}
		}
		if v, ok := r["Default"]; ok {
			switch dv := v.(type) {
			case nil:
			case []byte:
				col.DefaultValue = string(dv)
			case sql.NullString:
				if dv.Valid {
					col.DefaultValue = dv.String
				}
			default:
				col.DefaultValue = fmt.Sprint(dv)
			}
		}
		columns = append(columns, col)
	}
	return columns, nil
}
