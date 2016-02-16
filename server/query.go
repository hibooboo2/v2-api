package server

import "bytes"

type Query struct {
	Table      string
	Fields     string
	Conditions []string
	ByAccount  bool
	RemoveNull bool
}

func NewQuery(table string) Query {
	return Query{
		Table:      table,
		Fields:     "*",
		ByAccount:  true,
		RemoveNull: true,
	}
}

func (q Query) Buffer() *bytes.Buffer {
	buf := &bytes.Buffer{}
	buf.WriteString("SELECT ")
	buf.WriteString(q.Fields)
	buf.WriteString(" FROM ")
	buf.WriteString(q.Table)
	buf.WriteString(" WHERE ")
	if q.ByAccount {
		buf.WriteString("account_id = :account_id ")
	}
	if q.RemoveNull {
		buf.WriteString("removed IS NULL")
	}
	for _, cond := range q.Conditions {
		buf.WriteString(cond)
		buf.WriteString(" ")
	}
	return buf
}

func (q Query) String() string {
	return q.Buffer().String()
}
