package database

import (
	"context"

	"ksogit.kingsoft.net/chat/lib/xmysql"
	xmysqlv2 "ksogit.kingsoft.net/chat/lib/xmysql/v2"
	"ksogit.kingsoft.net/kgo/mysql"
)

type KingsoftDB struct {
	db mysql.DBAdapter
}

// Connect 建立MySQL连接
func (m *KingsoftDB) Connect(dsn string) error {
	dbConfig, err := xmysql.GetDatabaseFromDSN(dsn)
	if err != nil {
		return err
	}
	db, err := xmysqlv2.NewDBBuilder(dbConfig, &xmysqlv2.ServiceInfo{}).WithNameSuffix("master").Build(context.Background())
	if err != nil {
		return err
	}
	m.db = db
	return nil
}
