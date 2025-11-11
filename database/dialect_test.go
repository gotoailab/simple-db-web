package database

import (
	"testing"
)

func TestDamengDialect_getColumnsFromSchema(t *testing.T) {
	// 创建一个 DamengDialect 实例用于测试（不需要真实的数据库连接）
	dialect := &DamengDialect{BaseDialect: &BaseDialect{}}

	tests := []struct {
		name          string
		schema        string
		expectedCols  []ColumnInfo
		expectedError bool
	}{
		{
			name: "测试用例：kim_partition_dbinst 表",
			schema: `CREATE TABLE "test_kim_10132903552522188748"."kim_partition_dbinst"
(
"id" BIGINT IDENTITY(1, 1) NOT NULL,
"src_database" VARCHAR(64 CHAR) DEFAULT '' NOT NULL,
"src_table" VARCHAR(64 CHAR) DEFAULT '' NOT NULL,
"mgr_id" BIGINT DEFAULT 0 NOT NULL,
"db_inst" VARCHAR(64 CHAR) DEFAULT '' NOT NULL,
"dst_table" VARCHAR(64 CHAR) DEFAULT '' NOT NULL,
"start_pos" BIGINT DEFAULT 0 NOT NULL,
"end_pos" BIGINT DEFAULT 0 NOT NULL,
NOT CLUSTER PRIMARY KEY("id")) STORAGE(ON "MAIN", CLUSTERBTR) ;`,
			expectedCols: []ColumnInfo{
				{Name: "id", Type: "BIGINT", DefaultValue: "", Key: "PRI"},
				{Name: "src_database", Type: "VARCHAR(64 CHAR)", DefaultValue: "", Key: ""},
				{Name: "src_table", Type: "VARCHAR(64 CHAR)", DefaultValue: "", Key: ""},
				{Name: "mgr_id", Type: "BIGINT", DefaultValue: "0", Key: ""},
				{Name: "db_inst", Type: "VARCHAR(64 CHAR)", DefaultValue: "", Key: ""},
				{Name: "dst_table", Type: "VARCHAR(64 CHAR)", DefaultValue: "", Key: ""},
				{Name: "start_pos", Type: "BIGINT", DefaultValue: "0", Key: ""},
				{Name: "end_pos", Type: "BIGINT", DefaultValue: "0", Key: ""},
			},
			expectedError: false,
		},
		{
			name: "测试包含字符串 DEFAULT 值的列",
			schema: `CREATE TABLE "test_table"
(
"id" BIGINT IDENTITY(1, 1) NOT NULL,
"name" VARCHAR(100) DEFAULT 'test' NOT NULL,
"email" VARCHAR(255) DEFAULT '' NOT NULL,
"description" VARCHAR(500) DEFAULT 'default description' NOT NULL
) ;`,
			expectedCols: []ColumnInfo{
				{Name: "id", Type: "BIGINT", DefaultValue: "", Key: "PRI"},
				{Name: "name", Type: "VARCHAR(100)", DefaultValue: "test", Key: ""},
				{Name: "email", Type: "VARCHAR(255)", DefaultValue: "", Key: ""},
				{Name: "description", Type: "VARCHAR(500)", DefaultValue: "default description", Key: ""},
			},
			expectedError: false,
		},
		{
			name: "测试数字 DEFAULT 值（带单引号）",
			schema: `CREATE TABLE "test_table"
(
"id" BIGINT IDENTITY(1, 1) NOT NULL,
"age" INT DEFAULT '18' NOT NULL,
"score" DECIMAL(10,2) DEFAULT '0.00' NOT NULL
) ;`,
			expectedCols: []ColumnInfo{
				{Name: "id", Type: "BIGINT", DefaultValue: "", Key: "PRI"},
				{Name: "age", Type: "INT", DefaultValue: "18", Key: ""},
				{Name: "score", Type: "DECIMAL(10,2)", DefaultValue: "0.00", Key: ""},
			},
			expectedError: false,
		},
		{
			name: "测试空表结构",
			schema: `CREATE TABLE "empty_table"
(
) ;`,
			expectedCols: []ColumnInfo{},
			expectedError: false,
		},
		{
			name: "测试没有 DEFAULT 值的列",
			schema: `CREATE TABLE "test_table"
(
"id" BIGINT IDENTITY(1, 1) NOT NULL,
"name" VARCHAR(100) NOT NULL,
"created_at" TIMESTAMP NOT NULL
) ;`,
			expectedCols: []ColumnInfo{
				{Name: "id", Type: "BIGINT", DefaultValue: "", Key: "PRI"},
				{Name: "name", Type: "VARCHAR(100)", DefaultValue: "", Key: ""},
				{Name: "created_at", Type: "TIMESTAMP", DefaultValue: "", Key: ""},
			},
			expectedError: false,
		},
		{
			name: "测试混合情况：有 DEFAULT 和没有 DEFAULT 的列",
			schema: `CREATE TABLE "test_table"
(
"id" BIGINT IDENTITY(1, 1) NOT NULL,
"name" VARCHAR(100) NOT NULL,
"status" VARCHAR(20) DEFAULT 'active' NOT NULL,
"created_at" TIMESTAMP NOT NULL,
"updated_at" TIMESTAMP DEFAULT '2024-01-01' NOT NULL
) ;`,
			expectedCols: []ColumnInfo{
				{Name: "id", Type: "BIGINT", DefaultValue: "", Key: "PRI"},
				{Name: "name", Type: "VARCHAR(100)", DefaultValue: "", Key: ""},
				{Name: "status", Type: "VARCHAR(20)", DefaultValue: "active", Key: ""},
				{Name: "created_at", Type: "TIMESTAMP", DefaultValue: "", Key: ""},
				{Name: "updated_at", Type: "TIMESTAMP", DefaultValue: "2024-01-01", Key: ""},
			},
			expectedError: false,
		},
		{
			name: "测试复杂类型定义",
			schema: `CREATE TABLE "test_table"
(
"id" BIGINT IDENTITY(1, 1) NOT NULL,
"data" VARCHAR(100 CHAR) DEFAULT 'default' NOT NULL,
"value" NUMERIC(10,2) DEFAULT '0.00' NOT NULL
) ;`,
			expectedCols: []ColumnInfo{
				{Name: "id", Type: "BIGINT", DefaultValue: "", Key: "PRI"},
				{Name: "data", Type: "VARCHAR(100 CHAR)", DefaultValue: "default", Key: ""},
				{Name: "value", Type: "NUMERIC(10,2)", DefaultValue: "0.00", Key: ""},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 直接调用私有方法 getColumnsFromSchema
			columns, err := dialect.getColumnsFromSchema(tt.schema)

			// 验证错误
			if tt.expectedError {
				if err == nil {
					t.Errorf("期望返回错误，但实际没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望返回错误，但实际返回: %v", err)
				return
			}

			// 验证列数量
			if len(columns) != len(tt.expectedCols) {
				t.Errorf("列数量不匹配: 期望 %d, 实际 %d", len(tt.expectedCols), len(columns))
				t.Logf("期望的列: %+v", tt.expectedCols)
				t.Logf("实际的列: %+v", columns)
				return
			}

			// 验证每列的信息
			for i, expectedCol := range tt.expectedCols {
				if i >= len(columns) {
					t.Errorf("列索引 %d 超出范围", i)
					continue
				}

				actualCol := columns[i]
				if actualCol.Name != expectedCol.Name {
					t.Errorf("列 %d 名称不匹配: 期望 %s, 实际 %s", i, expectedCol.Name, actualCol.Name)
				}
				if actualCol.Type != expectedCol.Type {
					t.Errorf("列 %d 类型不匹配: 期望 %s, 实际 %s", i, expectedCol.Type, actualCol.Type)
				}
				if actualCol.DefaultValue != expectedCol.DefaultValue {
					t.Errorf("列 %d 默认值不匹配: 期望 %s, 实际 %s", i, expectedCol.DefaultValue, actualCol.DefaultValue)
				}
			}
		})
	}
}

