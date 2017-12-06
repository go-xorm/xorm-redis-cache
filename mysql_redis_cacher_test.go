package xormrediscache

import (
	"database/sql"
	_ "encoding/gob"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	. "github.com/go-xorm/tests"
	"github.com/go-xorm/xorm"
)

/*
CREATE DATABASE IF NOT EXISTS xorm_test CHARACTER SET
utf8 COLLATE utf8_general_ci;
*/
func TestMysqlWithCache(t *testing.T) {
	err := mysqlDdlImport()
	if err != nil {
		t.Error(err)
		return
	}

	engine, err := xorm.NewEngine("mysql", "root:@/xorm_test2?charset=utf8")
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetDefaultCacher(NewRedisCacher("localhost:6379", "", DEFAULT_EXPIRATION, engine.Logger()))
	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_DEBUG)

	BaseTestAll(engine, t)
	BaseTestAllSnakeMapper(engine, t)
	BaseTestAll2(engine, t)
}

func TestMysqlWithCacheSameMapper(t *testing.T) {
	err := mysqlDdlImport()
	if err != nil {
		t.Error(err)
		return
	}

	engine, err := xorm.NewEngine("mysql", "root:@/xorm_test3?charset=utf8")
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}
	engine.SetMapper(core.SameMapper{})
	engine.SetDefaultCacher(NewRedisCacher("localhost:6379", "", DEFAULT_EXPIRATION, engine.Logger()))
	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_DEBUG)

	BaseTestAll(engine, t)
	BaseTestAllSameMapper(engine, t)
	BaseTestAll2(engine, t)
}

func newMysqlEngine() (*xorm.Engine, error) {
	return xorm.NewEngine("mysql", "root:@/xorm_test?charset=utf8")
}

func newMysqlEngineWithCacher() (*xorm.Engine, error) {
	engine, err := newMysqlEngine()
	if err == nil {
		engine.SetDefaultCacher(NewRedisCacher("localhost:6379", "", DEFAULT_EXPIRATION, engine.Logger()))
	}
	return engine, err
}

func mysqlDdlImport() error {
	engine, err := xorm.NewEngine("mysql", "root:@/?charset=utf8")
	if err != nil {
		return err
	}
	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_DEBUG)

	sqlResults, _ := engine.ImportFile("../testdata/mysql_ddl.sql")
	engine.Logger().Debugf("sql results: %v", sqlResults)
	engine.Close()
	return nil
}

func newMysqlDriverDB() (*sql.DB, error) {
	return sql.Open("mysql", "root:@/xorm_test?charset=utf8")
}

func BenchmarkMysqlDriverInsert(t *testing.B) {
	DoBenchDriver(newMysqlDriverDB, CreateTableMySql, DropTableMySql,
		DoBenchDriverInsert, t)
}

func BenchmarkMysqlDriverFind(t *testing.B) {
	DoBenchDriver(newMysqlDriverDB, CreateTableMySql, DropTableMySql,
		DoBenchDriverFind, t)
}

func BenchmarkMysqlCacheInsert(t *testing.B) {
	engine, err := newMysqlEngineWithCacher()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}

	DoBenchInsert(engine, t)
}

func BenchmarkMysqlCacheFind(t *testing.B) {
	engine, err := newMysqlEngineWithCacher()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}

	DoBenchFind(engine, t)
}

func BenchmarkMysqlCacheFindPtr(t *testing.B) {
	engine, err := newMysqlEngineWithCacher()
	defer engine.Close()
	if err != nil {
		t.Error(err)
		return
	}

	DoBenchFindPtr(engine, t)
}
