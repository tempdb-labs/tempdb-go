package lib

import (
	"encoding/json"
	"fmt"
	"strings"
)

type QueryBuilder struct {
	operations []string
}

func NewQuery() *QueryBuilder {
	return &QueryBuilder{operations: make([]string, 0)}
}

func (ab *QueryBuilder) Count() *QueryBuilder {
	ab.operations = append(ab.operations, "COUNT")
	return ab
}

func (ab *QueryBuilder) Sum(field string) *QueryBuilder {
	ab.operations = append(ab.operations, fmt.Sprintf("SUM /%s", field))
	return ab
}

func (ab *QueryBuilder) Average(field string) *QueryBuilder {
	ab.operations = append(ab.operations, fmt.Sprintf("AVG /%s", field))
	return ab
}

func (ab *QueryBuilder) GroupBy(field string) *QueryBuilder {
	ab.operations = append(ab.operations, fmt.Sprintf("GROUPBY /%s", field))
	return ab
}

func (ab *QueryBuilder) Filter(field, operator, value string) *QueryBuilder {
	ab.operations = append(ab.operations, fmt.Sprintf("FILTER /%s %s %s", field, operator, value))
	return ab
}

func (qb *QueryBuilder) Min(field string) *QueryBuilder {
	qb.operations = append(qb.operations, fmt.Sprintf("MIN /%s", field))
	return qb
}

func (qb *QueryBuilder) Max(field string) *QueryBuilder {
	qb.operations = append(qb.operations, fmt.Sprintf("MAX /%s", field))
	return qb
}

func (qb *QueryBuilder) Distinct(field string) *QueryBuilder {
	qb.operations = append(qb.operations, fmt.Sprintf("DISTINCT /%s", field))
	return qb
}

func (qb *QueryBuilder) TopN(n int, field string) *QueryBuilder {
	qb.operations = append(qb.operations, fmt.Sprintf("TOPN %d /%s", n, field))
	return qb
}

func (qb *QueryBuilder) BottomN(n int, field string) *QueryBuilder {
	qb.operations = append(qb.operations, fmt.Sprintf("BOTTOMN %d /%s", n, field))
	return qb
}

// Enhanced filter operations
func (qb *QueryBuilder) FilterEquals(field, value string) *QueryBuilder {
	return qb.Filter(field, "eq", value)
}

func (qb *QueryBuilder) FilterNotEquals(field, value string) *QueryBuilder {
	return qb.Filter(field, "neq", value)
}

func (qb *QueryBuilder) FilterGreaterThan(field, value string) *QueryBuilder {
	return qb.Filter(field, "gt", value)
}

func (qb *QueryBuilder) FilterLessThan(field, value string) *QueryBuilder {
	return qb.Filter(field, "lt", value)
}

func (qb *QueryBuilder) FilterStartsWith(field, value string) *QueryBuilder {
	return qb.Filter(field, "startswith", value)
}

func (qb *QueryBuilder) FilterEndsWith(field, value string) *QueryBuilder {
	return qb.Filter(field, "endswith", value)
}

func (qb *QueryBuilder) FilterContains(field, value string) *QueryBuilder {
	return qb.Filter(field, "contains", value)
}

func (qb *QueryBuilder) FilterIn(field string, values []string) *QueryBuilder {
	jsonArray, _ := json.Marshal(values)
	return qb.Filter(field, "in", string(jsonArray))
}

func (qb *QueryBuilder) FilterNotIn(field string, values []string) *QueryBuilder {
	jsonArray, _ := json.Marshal(values)
	return qb.Filter(field, "notin", string(jsonArray))
}

func (qb *QueryBuilder) FilterExists(field string) *QueryBuilder {
	return qb.Filter(field, "exists", "true")
}

func (qb *QueryBuilder) FilterNotExists(field string) *QueryBuilder {
	return qb.Filter(field, "notexists", "true")
}

func (qb *QueryBuilder) FilterRegex(field, pattern string) *QueryBuilder {
	return qb.Filter(field, "regex", pattern)
}

func (ab *QueryBuilder) Build() string {
	return strings.Join(ab.operations, " ")
}

func (c *TempDBClient) Query(pipeline string) (interface{}, error) {
	return c.sendCommand(fmt.Sprintf("QUERY %s", pipeline))
}

func (c *TempDBClient) QueryWithBuilder(builder *QueryBuilder) (interface{}, error) {
	return c.Query(builder.Build())
}
