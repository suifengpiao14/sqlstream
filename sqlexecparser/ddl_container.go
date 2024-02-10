package sqlexecparser

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

var tablePool sync.Map

func getTablePoolKey(database string, tableName string) (key string) {
	return fmt.Sprintf("%s_%s", database, tableName)
}
func RegisterTable(database string, table Table) {
	key := getTablePoolKey(database, table.TableName)
	tablePool.Store(key, &table)
}

var (
	ERROR_NOT_FOUND_TABLE = errors.New("not found table")
	ERROR_INVALID_TYPE    = errors.New("invalid type, except *parser.Table")
)

func GetTable(database string, tableName string) (table *Table, err error) {
	key := getTablePoolKey(database, tableName)
	v, ok := tablePool.Load(key)
	if !ok {
		err = errors.WithMessagef(ERROR_NOT_FOUND_TABLE, "%s", key)
		return nil, err
	}
	table, ok = v.(*Table)
	if !ok {
		return nil, ERROR_INVALID_TYPE
	}
	return table, nil
}

// RegisterTableByDDL 通过ddl语句注册表结构,避免依赖db连接,方便本地化启动模块
func RegisterTableByDDL(database string, ddlStatements string) (err error) {
	tables, err := ParseCreateDDL(ddlStatements)
	if err != nil {
		return err
	}
	for _, table := range tables {
		RegisterTable(database, table)
	}
	return nil
}