package db

import (
	"database/sql/driver"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/pkg/errors"
)

var driverValuerType = reflect.TypeOf((*driver.Valuer)(nil)).Elem()

type indexTag struct {
	Name    string
	Type    string
	Columns []string
}

type tagInfo struct {
	Key string
	Val string
}

// CreateDatabase 如果数据库不存在，则创建
func CreateDatabase(dbConfig Config) (err error) {
	dbName, dsnWithoutDBName, err := splitDSN(dbConfig.Driver, dbConfig.DSN)
	if err != nil || len(dbName) == 0 || len(dsnWithoutDBName) == 0 {
		return
	}

	conn, err := dbr.Open(dbConfig.Driver, dsnWithoutDBName, NewEventLogger(true))
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	return createDatabase(conn.NewSession(nil), dbName)
}

func createDatabase(session *dbr.Session, dbName string) (err error) {
	switch getBaseDialect(session) {
	case dialect.MySQL:
		sqlstr := "CREATE DATABASE IF NOT EXISTS " +
			session.QuoteIdent(dbName) +
			" DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci"
		_, err = session.InsertBySql(sqlstr).Exec()
		return errors.WithStack(err)

	case dialect.PostgreSQL:
		count, _ := session.Select("count(*)").From("pg_database").Where("datname=?", dbName).ReturnInt64()
		if count == 0 {
			sqlstr := "CREATE DATABASE " + session.QuoteIdent(dbName)
			_, err = session.InsertBySql(sqlstr).Exec()
			return errors.WithStack(err)
		}

	}

	return
}

// CreateDatabase 如果表不存在，则创建
func CreateTable(session *dbr.Session, table string, schema interface{}) (err error) {
	baseDialect := getBaseDialect(session)
	rows, err := session.Query("SELECT * FROM " + session.QuoteIdent(table) + " WHERE 1 != 1")
	if err != nil {
		createSQL := schema2CreateTableSQL(baseDialect, table, schema)
		_, err = session.InsertBySql(createSQL).Exec()
		if err != nil {
			return errors.WithStack(err)
		}
		for _, indexSQL := range listTableIndexSQL(baseDialect, table, schema) {
			if _, err = session.InsertBySql(indexSQL).Exec(); err != nil {
				return errors.WithStack(err)
			}
		}
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return errors.WithStack(err)
	}
	quotedTable := session.QuoteIdent(table)
	tp := reflect.Indirect(reflect.ValueOf(schema)).Type()

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		columnName := dbr.NameMapping(field.Name)

		if field.Tag.Get("db") == "-" {
			continue
		}

		if index := findIndex(columns, columnName); index < 0 {
			sqlstr := fmt.Sprintf("ALTER TABLE %s ADD %s", quotedTable, field2SQL(baseDialect, field))
			if _, err = session.InsertBySql(sqlstr).Exec(); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return
}

func schema2CreateTableSQL(d dbr.Dialect, table string, schema interface{}) string {
	sqlstr := "CREATE TABLE IF NOT EXISTS " + d.QuoteIdent(table)
	sqlstr += "(\n"

	tp := reflect.Indirect(reflect.ValueOf(schema)).Type()
	columns := []string{}

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)

		if field.Tag.Get("db") == "-" {
			continue
		}
		columns = append(columns, field2SQL(d, field))
	}

	sqlstr += strings.Join(columns, ",\n") + "\n)"

	return sqlstr
}

func field2SQL(d dbr.Dialect, field reflect.StructField) string {
	buf := strings.Builder{}
	columnName := dbr.NameMapping(field.Name)
	tag := field.Tag
	kind := fieldKind(field)

	buf.WriteString(d.QuoteIdent(columnName))
	buf.WriteString(" ")

	sqlTag := tag.Get("sql")
	tags := parseTags2Map(sqlTag)
	defaultLen, defaultVal := "", ""

	if strings.ToLower(field.Name) == "id" &&
		len(tags["index"]) == 0 &&
		kindType(kind) == kindType_Integer {
		pKeySQL := "PRIMARY KEY"

		switch d {
		case dialect.PostgreSQL:
			pKeySQL = "SERIAL " + pKeySQL
		case dialect.MySQL:
			pKeySQL = "BIGINT UNSIGNED NOT NULL AUTO_INCREMENT " + pKeySQL
		default:
			pKeySQL = "INTEGER " + pKeySQL
		}

		buf.WriteString(pKeySQL)

		return buf.String()
	}

	if typ := tags["type"]; len(typ) > 0 {
		buf.WriteString(strings.ToUpper(typ))
	} else {
		switch kindType(kind) {
		case kindType_Boolean:
			if d == dialect.PostgreSQL {
				buf.WriteString("BOOLEAN")
				defaultVal = "'false'"
			} else {
				buf.WriteString("TINYINT")
				defaultVal = "'0'"
			}

		case kindType_Integer:
			buf.WriteString("INTEGER")
			defaultVal = "'0'"

		case kindType_Float:
			buf.WriteString("DECIMAL")
			defaultLen, defaultVal = "(20,2)", "'0'"

		case kindType_String:
			buf.WriteString("VARCHAR")
			defaultLen, defaultVal = "(255)", "''"

		case kindType_Object:
			if isDateTimeType(fieldType(field)) {
				if d == dialect.PostgreSQL {
					buf.WriteString("TIMESTAMP")
				} else {
					buf.WriteString("DATETIME")
				}
			} else {
				buf.WriteString("TEXT")
			}
		}
	}

	if length := tags["length"]; len(length) > 0 {
		buf.WriteString("(" + length + ")")
	} else if len(defaultLen) > 0 {
		buf.WriteString(defaultLen)
	}

	buf.WriteString(" ")

	if isUnsignedNumber(kind) || tags["unsigned"] == "true" {
		buf.WriteString("UNSIGNED ")
	}

	if strings.ToLower(field.Name) == "id" && len(tags["index"]) == 0 {
		buf.WriteString("PRIMARY KEY")

		return buf.String()
	}

	if field.Type.Kind() == reflect.Ptr ||
		field.Type.Implements(driverValuerType) {
		buf.WriteString("NULL")
		defaultVal = ""
	} else {
		buf.WriteString("NOT NULL")
	}

	if def, ok := tags["default"]; ok {
		buf.WriteString(" DEFAULT " + def)
	} else if len(defaultVal) > 0 {
		buf.WriteString(" DEFAULT " + defaultVal)
	}

	return buf.String()
}

func listTableIndexSQL(d dbr.Dialect, table string, schema interface{}) (sqls []string) {
	for _, index := range listIndexTags(schema) {
		indexType := "INDEX"

		switch strings.ToLower(index.Type) {
		case "unique":
			indexType = "UNIQUE INDEX"
		case "primary":
			indexType = "PRIMARY KEY"
		}

		columns := []string{}

		for _, col := range index.Columns {
			columns = append(columns, d.QuoteIdent(col))
		}

		indexSQL := fmt.Sprintf("CREATE %s %s ON %s (%s)",
			indexType, d.QuoteIdent(index.Name), d.QuoteIdent(table), strings.Join(columns, ","))
		sqls = append(sqls, indexSQL)
	}

	return
}

func getBaseDialect(session *dbr.Session) dbr.Dialect {
	baseDialect := session.Dialect

	if ld, ok := baseDialect.(*localDialect); ok {
		baseDialect = ld.Dialect
	}

	return baseDialect
}

func splitDSN(driver, dsn string) (dbName, dsnWithoutDBName string, err error) {
	switch driver {
	case "mysql":
		cfg, err1 := mysql.ParseDSN(dsn)
		if err1 != nil {
			err = errors.WithStack(err1)
			return
		}

		dbName = cfg.DBName
		cfg.DBName = ""
		dsnWithoutDBName = cfg.FormatDSN()

	case "postgres":
		dsnURL, err1 := url.Parse(dsn)
		if err1 != nil {
			err = errors.WithStack(err1)
			return
		}

		if len(dsnURL.Path) > 0 {
			dbName = strings.Trim(dsnURL.Path, "/")
			dsnURL.Path = "/"
			dsnWithoutDBName = dsnURL.String()
		}
	}

	return
}

func parseTags(sqlTag string) (tags []tagInfo) {
	if len(sqlTag) == 0 {
		return
	}

	for _, tagField := range strings.Fields(sqlTag) {
		parts := strings.SplitN(tagField, ":", 2)
		if len(parts) < 2 {
			parts = append(parts, "true")
		}
		tags = append(tags, tagInfo{
			Key: strings.TrimSpace(parts[0]),
			Val: strings.TrimSpace(parts[1]),
		})
	}

	return
}

func parseTags2Map(sqlTag string) map[string]string {
	tags := make(map[string]string)

	for _, tag := range parseTags(sqlTag) {
		tags[tag.Key] = tag.Val
	}

	return tags
}

func listIndexTags(schema interface{}) (tags []*indexTag) {
	tp := reflect.Indirect(reflect.ValueOf(schema)).Type()

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		sqlTags := parseTags(field.Tag.Get("sql"))
		column := dbr.NameMapping(field.Name)

		for _, tag := range sqlTags {
			if tag.Key != "index" {
				continue
			}
			parts := strings.SplitN(tag.Val, ",", 2)
			indexName := parts[0]
			indexType := "index"

			if len(parts) > 1 {
				indexType = parts[1]
			}

			if index := findIndexTag(tags, indexName); index >= 0 {
				tags[index].Columns = append(tags[index].Columns, column)
			} else {
				tags = append(tags, &indexTag{
					Name:    indexName,
					Type:    indexType,
					Columns: []string{column},
				})
			}
		}
	}
	return
}

const (
	kindType_Integer = "integer"
	kindType_Float   = "float"
	kindType_String  = "string"
	kindType_Boolean = "boolean"
	kindType_Object  = "object"
)

func kindType(kind reflect.Kind) string {
	switch kind {
	case reflect.Bool:
		return kindType_Boolean

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return kindType_Integer

	case reflect.Float32, reflect.Float64:
		return kindType_Float

	case reflect.String:
		return kindType_String

	default:
		return kindType_Object
	}
}

func isUnsignedNumber(kind reflect.Kind) bool {
	return kind == reflect.Uint ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64
}

// isDateTimeType check if typ is datetime. typ should be struct type
func isDateTimeType(typ reflect.Type) bool {
	return typ == reflect.TypeOf(dbr.NullTime{}) ||
		typ == reflect.TypeOf(NullTime{}) ||
		typ == reflect.TypeOf(time.Time{})
}

func fieldType(field reflect.StructField) reflect.Type {
	fk := field.Type.Kind()

	if fk == reflect.Ptr {
		return field.Type.Elem()
	}

	return field.Type
}

func fieldKind(field reflect.StructField) reflect.Kind {
	return fieldType(field).Kind()
}

func findIndex(list []string, e string) int {
	for i, s := range list {
		if s == e {
			return i
		}
	}
	return -1
}

func findIndexTag(tags []*indexTag, name string) int {
	for i, tag := range tags {
		if tag.Name == name {
			return i
		}
	}
	return -1
}
