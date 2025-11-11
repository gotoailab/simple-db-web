package database

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"ksogit.kingsoft.net/kgo/mysql"
)

func GetDialectByType(dbType string, db mysql.DBAdapter) Dialect {
	switch dbType {
	case "dameng":
		return NewDamengDialect(db)
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

type DamengDialect struct {
	*BaseDialect
}

func NewDamengDialect(db mysql.DBAdapter) *DamengDialect {
	return &DamengDialect{BaseDialect: NewBaseDialect(db)}
}

func (m *DamengDialect) GetTableColumns(tableName string) ([]ColumnInfo, error) {
	schema, err := m.BaseDialect.GetTableSchema(tableName)
	if err != nil {
		return nil, err
	}
	return m.getColumnsFromSchema(schema)
}

func (m *DamengDialect) getColumnsFromSchema(schema string) ([]ColumnInfo, error) {
	columns := []ColumnInfo{}
	
	// 匹配 IDENTITY 模式，用于识别主键
	identityReg := regexp.MustCompile(`IDENTITY\s*\([^)]+\)`)
	
	// 匹配所有列定义（包括有 DEFAULT 和没有 DEFAULT 的）
	// 格式1: "column_name" TYPE IDENTITY(1, 1) NOT NULL,
	// 格式2: "column_name" TYPE DEFAULT 'value' NOT NULL,
	// 格式3: "column_name" TYPE DEFAULT 0 NOT NULL,
	// 注意：使用 (.+?) 非贪婪匹配，匹配到 DEFAULT 或 NOT NULL 之前
	reg := regexp.MustCompile(`"([^"]+)"\s+(.+?)(?:\s+DEFAULT\s+(?:'([^']*)'|(\S+)))?\s+NOT\s+NULL`)
	matches := reg.FindAllStringSubmatch(schema, -1)
	
	for _, match := range matches {
		col := ColumnInfo{
			Name: match[1],
		}
		
		typeAndAttrs := strings.TrimSpace(match[2])
		
		// 检查是否包含 IDENTITY，如果包含则标记为主键
		if identityReg.MatchString(typeAndAttrs) {
			col.Key = "PRI"
			// 从类型字符串中移除 IDENTITY 部分
			typeAndAttrs = identityReg.ReplaceAllString(typeAndAttrs, "")
			typeAndAttrs = strings.TrimSpace(typeAndAttrs)
			// 移除 IDENTITY 后的多余空格
			typeAndAttrs = strings.Trim(typeAndAttrs, " ")
		}
		col.Type = typeAndAttrs
		
		// 默认值可能是单引号内的值（match[3]）或数字值（match[4]）
		if match[3] != "" {
			col.DefaultValue = match[3]
		} else if match[4] != "" {
			col.DefaultValue = match[4]
		}
		
		columns = append(columns, col)
	}
	
	return columns, nil
}