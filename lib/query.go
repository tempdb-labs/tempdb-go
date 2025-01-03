package lib

import (
	"encoding/json"
	"fmt"
)

type SortOrder string

const (
	Ascending  SortOrder = "Ascending"
	Descending SortOrder = "Descending"
)

type Operator string

const (
	Eq         Operator = "Eq"
	Gt         Operator = "Gt"
	Lt         Operator = "Lt"
	Gte        Operator = "Gte"
	Lte        Operator = "Lte"
	Contains   Operator = "Contains"
	StartsWith Operator = "StartsWith"
	EndsWith   Operator = "EndsWith"
)

type Condition struct {
	Field    string      `json:"field"`
	Operator Operator    `json:"operator"`
	Value    interface{} `json:"value"`
}

type Query struct {
	Conditions []Condition `json:"conditions"`
	SortBy     *string     `json:"sort_by,omitempty"`
	SortOrder  *SortOrder  `json:"sort_order,omitempty"`
	Limit      *int        `json:"limit,omitempty"`
	Offset     *int        `json:"offset,omitempty"`
	TimeRange  *struct {
		Start int64 `json:"start"`
		End   int64 `json:"end"`
	} `json:"time_range,omitempty"`
}

type QueryBuilder struct {
	query Query
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		query: Query{
			Conditions: make([]Condition, 0),
		},
	}
}

func (qb *QueryBuilder) WhereEqual(field string, value interface{}) *QueryBuilder {
	qb.query.Conditions = append(qb.query.Conditions, Condition{
		Field:    field,
		Operator: Eq,
		Value:    value,
	})
	return qb
}

func (qb *QueryBuilder) Sort(field string, order SortOrder) *QueryBuilder {
	qb.query.SortBy = &field
	qb.query.SortOrder = &order
	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.query.Limit = &limit
	return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.query.Offset = &offset
	return qb
}

func (qb *QueryBuilder) TimeRange(start int64, end int64) *QueryBuilder {
	qb.query.TimeRange = &struct {
		Start int64 `json:"start"`
		End   int64 `json:"end"`
	}{Start: start, End: end}
	return qb
}

func (qb *QueryBuilder) Build() Query {
	return qb.query
}

func (c *TempDBClient) Query(query Query) (interface{}, error) {
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return "", err
	}
	return c.sendCommand(fmt.Sprintf("QUERY %s", (string(queryJSON))))
}

func (c *TempDBClient) QueryPipeline(pipeline string) (interface{}, error) {
	return c.sendCommand(fmt.Sprintf("PIPELINE %s", pipeline))
}

// Convience method for common query patterns
func (c *TempDBClient) QueryWhere(field, value string) (interface{}, error) {
	pipeline := fmt.Sprintf("where %s eq %s", field, value)
	return c.QueryPipeline(pipeline)
}

// for building complex pipelines
func (c *TempDBClient) QueryBuilder() *QueryBuilder {
	return NewQueryBuilder()
}
