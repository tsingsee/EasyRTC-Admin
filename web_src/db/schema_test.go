package db

import (
	"reflect"
	"testing"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var configs = []Config{
	{
		Driver: "sqlite3",
		DSN:    ":memory:",
	},
	// {
	// 	Driver:   "postgres",
	// 	DSN:      "postgres://postgres:123456@localhost:5432/testdb?sslmode=disable",
	// 	Timezone: "Asia/Shanghai",
	// },
}
var sessions []*dbr.Session

func createSession(dbConfig Config) *dbr.Session {
	return NewSQLDB(dbConfig, true)
}

func init() {
	// time.Local = time.UTC

	for _, config := range configs {
		if err := CreateDatabase(config); err != nil {
			panic(err)
		}
		sessions = append(sessions, createSession(config))
	}

	for _, sess := range sessions {
		if err := dropTable(sess, "people"); err != nil {
			panic(err)
		}
		if err := CreateTable(sess, "people", People{}); err != nil {
			panic(err)
		}
	}
}

type TestTable struct {
	A string
	B int `sql:"default:'1'"`
	C struct{}
	D NullTime
	E *string
	F string `sql:"type:text"`
}

type TestTable2 struct {
	ID int
	A  string `sql:"index:a_b,unique"`
	B  string `sql:"index:a_b"`
	C  string `sql:"index:c"`
}

func dropTable(session *dbr.Session, table string) error {
	_, err := session.InsertBySql("Drop Table If Exists " + session.QuoteIdent(table)).Exec()
	return err
}

func TestCreateTable(t *testing.T) {
	for _, session := range sessions {
		dropTable(session, "test_table")

		err := CreateTable(session, "test_table", TestTable{})
		require.NoError(t, err)

		err = CreateTable(session, "test_table", TestTable{})
		require.NoError(t, err)

		dropTable(session, "test_table2")

		err = CreateTable(session, "test_table2", TestTable2{})
		require.NoError(t, err)
	}
}

func TestSplitDSN(t *testing.T) {
	dbName, dsn, err := splitDSN("mysql", "user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true")
	require.NoError(t, err)
	require.Equal(t, "dbname", dbName)
	require.Equal(t, "user:password@tcp(localhost:5555)/?tls=skip-verify&autocommit=true", dsn)

	dbName, dsn, err = splitDSN("postgres", "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")
	require.NoError(t, err)
	require.Equal(t, "pqgotest", dbName)
	require.Equal(t, "postgres://pqgotest:password@localhost/?sslmode=verify-full", dsn)
}

func TestSchema2CreateTableSQL(t *testing.T) {
	testCases := []struct {
		Driver  string
		Dialect dbr.Dialect
		Want    string
	}{
		{
			Driver:  "sqlite3",
			Dialect: dialect.SQLite3,
			Want: `CREATE TABLE IF NOT EXISTS "test_table"(
"a" VARCHAR(255) NOT NULL DEFAULT '',
"b" INTEGER NOT NULL DEFAULT '1',
"c" TEXT NOT NULL,
"d" DATETIME NULL,
"e" VARCHAR(255) NULL,
"f" TEXT NOT NULL
)`,
		},
		{
			Driver:  "mysql",
			Dialect: dialect.MySQL,
			Want: "CREATE TABLE IF NOT EXISTS `test_table`(\n" +
				"`a` VARCHAR(255) NOT NULL DEFAULT '',\n" +
				"`b` INTEGER NOT NULL DEFAULT '1',\n" +
				"`c` TEXT NOT NULL,\n" +
				"`d` DATETIME NULL,\n" +
				"`e` VARCHAR(255) NULL,\n" +
				"`f` TEXT NOT NULL\n)",
		},
		{
			Driver:  "postgres",
			Dialect: dialect.PostgreSQL,
			Want: `CREATE TABLE IF NOT EXISTS "test_table"(
"a" VARCHAR(255) NOT NULL DEFAULT '',
"b" INTEGER NOT NULL DEFAULT '1',
"c" TEXT NOT NULL,
"d" TIMESTAMP NULL,
"e" VARCHAR(255) NULL,
"f" TEXT NOT NULL
)`,
		},
	}

	for _, c := range testCases {
		sqlstr := schema2CreateTableSQL(c.Dialect, "test_table", TestTable{})
		require.Equal(t, c.Want, sqlstr, c.Driver)
	}
}

func TestListTableIndexSQL(t *testing.T) {
	type Test struct {
		A string `sql:"index:a_b,unique"`
		B string `sql:"index:a_b"`
	}
	sqls := listTableIndexSQL(dialect.MySQL, "test", Test{})
	require.Len(t, sqls, 1)

	require.Equal(t, "CREATE UNIQUE INDEX `a_b` ON `test` (`a`,`b`)", sqls[0])
}

func TestParseTags(t *testing.T) {
	type People struct {
		Id    int
		Name  string    `sql:"index:ne,unique"`
		Email string    `sql:"index:ne,unique"`
		Ctime time.Time `sql:"index:ctime"`
	}
	indexSQLs := listTableIndexSQL(dialect.MySQL, "people", People{})
	require.Equal(t, []string{
		"CREATE UNIQUE `ne` ON `people` (`name`,`email`)",
		"CREATE INDEX `ctime` ON `people` (`ctime`)",
	}, indexSQLs)
}

func TestTypeImplements(t *testing.T) {
	tp := reflect.TypeOf(NullTime{})
	assert.True(t, tp.Implements(driverValuerType))
}
