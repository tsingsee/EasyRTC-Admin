package db

import (
	"reflect"

	"github.com/gocraft/dbr/v2"
)

// 通用的查询接口
type Selector struct {
	session dbr.SessionRunner
	where   []dbr.Builder

	Table      interface{} `json:"table,omitempty"`
	JoinTables []JoinTable `json:"joinTables,omitempty"`
	Cols       []string    `json:"cols,omitempty"`
	Conditions []Condition `json:"conditions,omitempty"`
	Groups     []string    `json:"groups,omitempty"`
	Havings    []Condition `json:"havings,omitempty"`
	Orders     []Order     `json:"orders,omitempty"`
	Pagination Pagination  `json:"pagination,omitempty"`
}

type JoinType string

const (
	JoinType_Full  JoinType = "full"
	JoinType_Left           = "left"
	JoinType_Right          = "right"
)

type JoinTable struct {
	Table    interface{} `json:"table,omitempty"`
	JoinType JoinType    `json:"joinType,omitempty"`
	On       interface{} `json:"on,omitempty"`
}

type Condition struct {
	Col     string      `json:"col,omitempty"`
	Cmp     string      `json:"cmp,omitempty"`
	Val     interface{} `json:"val,omitempty"`
	Escapes []string    `json:"escapes,omitempty"`
}

const (
	CmpEq  = "eq"
	CmpGte = "gte"
	CmpLte = "lte"
)

func (c Condition) Build() dbr.Builder {
	switch c.Cmp {
	case "eq":
		return dbr.Eq(c.Col, c.Val)
	case "neq":
		return dbr.Neq(c.Col, c.Val)
	case "gt":
		return dbr.Gt(c.Col, c.Val)
	case "gte":
		return dbr.Gte(c.Col, c.Val)
	case "lt":
		return dbr.Lt(c.Col, c.Val)
	case "lte":
		return dbr.Lte(c.Col, c.Val)
	case "like":
		return dbr.Like(c.Col, c.Val.(string), c.Escapes...)
	case "notLike":
		return dbr.NotLike(c.Col, c.Val.(string), c.Escapes...)
	default:
		return nil
	}
}

type Order struct {
	Col string `json:"col"`
	Asc bool   `json:"asc"`
}

type Column struct {
	Name     string       `json:"name,omitempty"`
	Type     string       `json:"type,omitempty"`
	Nullable bool         `json:"nullable,omitempty"`
	Length   int64        `json:"length,omitempty"`
	Decimal  *DecimalSize `json:"decimal,omitempty"`
}

type DecimalSize struct {
	Precision int64 `json:"precision,omitempty"`
	Scale     int64 `json:"scale,omitempty"`
}

type SelectResult struct {
	Columns []Column        `json:"columns,omitempty"`
	Data    [][]interface{} `json:"data,omitempty"`
	Count   int64           `json:"count,omitempty"`
}

func (s SelectResult) Rows() (mm []map[string]interface{}) {
	for _, data := range s.Data {
		m := make(map[string]interface{})
		for i, val := range data {
			m[s.Columns[i].Name] = val
		}
		mm = append(mm, m)
	}

	return
}

func (s SelectResult) First() (m map[string]interface{}) {
	if mm := s.Rows(); len(mm) > 0 {
		return mm[0]
	}

	return
}

func NewSelector(session dbr.SessionRunner) *Selector {
	return &Selector{session: session}
}

func (s *Selector) DisableJoin() {
	s.JoinTables = nil
}

func (s *Selector) From(table interface{}) *Selector {
	s.Table = table

	return s
}

func (s *Selector) Join(table interface{}, on string) *Selector {
	s.JoinTables = append(s.JoinTables, JoinTable{
		Table: table,
		On:    on,
	})

	return s
}

func (s *Selector) LeftJoin(table interface{}, on string) *Selector {
	s.JoinTables = append(s.JoinTables, JoinTable{
		Table:    table,
		JoinType: JoinType_Left,
		On:       on,
	})

	return s
}

func (s *Selector) RightJoin(table interface{}, on string) *Selector {
	s.JoinTables = append(s.JoinTables, JoinTable{
		Table:    table,
		JoinType: JoinType_Right,
		On:       on,
	})

	return s
}

func (s *Selector) FullJoin(table interface{}, on string) *Selector {
	s.JoinTables = append(s.JoinTables, JoinTable{
		Table:    table,
		JoinType: JoinType_Full,
		On:       on,
	})

	return s
}

func (s *Selector) Paginate(page, perPage uint64) *Selector {
	s.Pagination = NewPagination(page, perPage)

	return s
}

func (s *Selector) Where(builders ...dbr.Builder) *Selector {
	s.where = append(s.where, builders...)

	return s
}

func (s *Selector) OrderDesc(col string) *Selector {
	s.Orders = append(s.Orders, Order{
		Col: col,
		Asc: false,
	})
	return s
}

func (s *Selector) OrderAsc(col string) *Selector {
	s.Orders = append(s.Orders, Order{
		Col: col,
		Asc: true,
	})
	return s
}

func (s Selector) Load(value interface{}) (count int, err error) {
	return s.Stmt().Load(value)
}

func (s Selector) Count() (count int64, err error) {
	stmt := s.session.Select("COUNT(*)").From(s.Table)

	for _, condition := range s.Conditions {
		if builder := condition.Build(); builder != nil {
			stmt.Where(builder)
		}
	}

	err = stmt.LoadOne(&count)

	return
}

func (s Selector) LoadPage(items interface{}, columns ...string) (result *PageResult, err error) {
	if tp := reflect.TypeOf(items); tp.Kind() != reflect.Ptr {
		panic("items is not pointer")
	}

	s.Cols = append(s.Cols, columns...)
	stmt := s.Stmt()

	if _, err = stmt.Load(items); err != nil {
		return
	}

	stmt.Column = []interface{}{"COUNT(*)"}
	stmt.HavingCond = nil
	stmt.Order = nil
	stmt.LimitCount = -1
	stmt.OffsetCount = -1
	stmt.Suffixes = nil

	var count dbr.NullInt64

	if _, err = stmt.Load(&count); err != nil {
		return
	}

	result = &PageResult{
		Items: reflect.Indirect(reflect.ValueOf(items)).Interface(),
		Count: count.Int64,
	}

	return
}

func (s Selector) Stmt() *dbr.SelectStmt {
	var cols []string

	for _, col := range s.Cols {
		cols = append(cols, col)
	}

	if len(cols) == 0 {
		cols = []string{"*"}
	}

	stmt := s.session.Select(cols...).From(s.Table)

	for _, join := range s.JoinTables {
		table := join.Table

		switch join.JoinType {
		case JoinType_Full:
			stmt.FullJoin(table, join.On)
		case JoinType_Left:
			stmt.LeftJoin(table, join.On)
		case JoinType_Right:
			stmt.RightJoin(table, join.On)
		default:
			stmt.Join(table, join.On)
		}
	}

	where := s.where[:]

	for _, condition := range s.Conditions {
		if builder := condition.Build(); builder != nil {
			where = append(where, builder)
		}
	}

	for _, builder := range where {
		stmt.Where(builder)
	}

	for _, col := range s.Groups {
		stmt.GroupBy(col)
	}

	for _, having := range s.Havings {
		if builder := having.Build(); builder != nil {
			stmt.Having(builder)
		}
	}

	for _, order := range s.Orders {
		stmt.OrderDir(order.Col, order.Asc)
	}

	return stmt.Paginate(s.Pagination.filter())
}
