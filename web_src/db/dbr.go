package db

import (
	"log"
	"time"

	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Driver   string `json:"driver,omitempty"`
	DSN      string `json:"dsn,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

type localDialect struct {
	dbr.Dialect
	loc *time.Location
}

func newLocalDialect(parentDialect dbr.Dialect, loc *time.Location) dbr.Dialect {
	return &localDialect{
		Dialect: parentDialect,
		loc:     loc,
	}
}

func (d localDialect) EncodeTime(t time.Time) string {
	return `'` + t.In(d.loc).Format("2006-01-02 15:04:05") + `'`
}

func loadLocation(dbConfig Config) (loc *time.Location, err error) {
	loc = time.UTC

	if len(dbConfig.Timezone) > 0 {
		loc, err = time.LoadLocation(dbConfig.Timezone)
	}

	return
}

func NewSQLDB(dbConfig Config, debug bool) (session *dbr.Session) {
	EventLogger := NewEventLogger(debug)
	loc, err := loadLocation(dbConfig)
	if err != nil {
		panic(err)
	}

	// TODO: 设置数据库的时区，无法处理多数据库不同时区
	Local = loc

	conn, err := dbr.Open(dbConfig.Driver, dbConfig.DSN, EventLogger)
	if err != nil {
		panic(err)
	}
	conn.Dialect = newLocalDialect(conn.Dialect, loc)

	return conn.NewSession(nil)
}

type EventLogger struct {
	dbr.NullEventReceiver
	debug bool
}

func NewEventLogger(debug bool) dbr.EventReceiver {
	return &EventLogger{debug: debug}
}

// Event receives a simple notification when various events occur.
func (n *EventLogger) Event(eventName string) {
	if n.debug {
		log.Printf("event: %s", eventName)
	}
}

// EventKv receives a notification when various events occur along with
// optional key/value data.
func (n *EventLogger) EventKv(eventName string, kvs map[string]string) {
	if n.debug {
		log.Printf("event: %s, sql: %v", eventName, kvs["sql"])
	}
}

// EventErr receives a notification of an error if one occurs.
func (n *EventLogger) EventErr(eventName string, err error) error {
	if err != dbr.ErrNotFound {
		log.Printf("event: %s, err: %s", eventName, err)
	}
	return err
}

// EventErrKv receives a notification of an error if one occurs along with
// optional key/value data.
func (n *EventLogger) EventErrKv(eventName string, err error, kvs map[string]string) error {
	if err != dbr.ErrNotFound {
		log.Printf("event: %s, err: %s, sql: %v", eventName, err, kvs["sql"])
	}
	return err
}

// Timing receives the time an event took to happen.
func (n *EventLogger) Timing(eventName string, nanoseconds int64) {
	if n.debug {
		log.Printf("event: %s, timing: %s", eventName, time.Duration(nanoseconds))
	}
}

// TimingKv receives the time an event took to happen along with optional key/value data.
func (n *EventLogger) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	if n.debug {
		log.Printf("event: %s, timing: %s, sql: %v", eventName, time.Duration(nanoseconds), kvs["sql"])
	}
}
