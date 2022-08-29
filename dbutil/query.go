package storage

import (
	"errors"
	"fmt"
	"strings"
)

// chain call method to create SQL
type Query struct {
	TableName   string
	selectExpr  map[string]string
	filterExpr  []string
	groupByExpr map[string]string
	aggExpr     []string
	havingExpr  []string
	orderByExpr []string
}

func (q *Query) From(name string) *Query {
	if q.TableName == "" {
		q.TableName = name
	}
	return q
}

func (q *Query) SelectAs(exprs map[string]string) *Query {
	if q.selectExpr == nil {
		q.selectExpr = make(map[string]string)
	}

	for field, alias := range exprs {
		if strings.TrimSpace(field) == "" {
			continue
		}

		q.selectExpr[strings.TrimSpace(field)] = strings.TrimSpace(alias)
	}

	return q
}

func (q *Query) Select(exprs ...string) *Query {
	if q.selectExpr == nil {
		q.selectExpr = make(map[string]string)
	}

	for _, expr := range exprs {
		if strings.TrimSpace(expr) == "" {
			continue
		}

		var field, alias string
		arr := strings.Split(expr, " as ")
		if len(arr) == 2 {
			field = strings.TrimSpace(arr[0])
			alias = strings.TrimSpace(arr[1])
		} else if len(arr) == 1 {
			field = strings.TrimSpace(arr[0])
			alias = ""
		}

		q.selectExpr[field] = alias
	}

	return q
}

func (q *Query) Selected() string {
	var selected []string
	for field, alias := range q.selectExpr {
		if alias == "" {
			selected = append(selected, field)
		} else {
			selected = append(selected, field+" as "+alias)
		}
	}

	return strings.Join(selected, ", ")
}

func (q *Query) Where(exprs ...string) *Query {
	return q.Filter(exprs...)
}

func (q *Query) WhereIn(exprs map[string][]string) *Query {
	return q.FilterIn(exprs)
}

func (q *Query) Filter(exprs ...string) *Query {
	q.filterExpr = q.append(q.filterExpr, exprs...)
	return q
}

// exprs key-value: field-inValues
// each key-value means an independent filter
func (q *Query) FilterIn(exprs map[string][]string) *Query {
	filters := make([]string, 0)

	for field, values := range exprs {
		filters = append(filters, fmt.Sprintf("%s in ('%s')", field, strings.Join(values, "','")))
	}

	return q.Filter(filters...)
}

func (q *Query) GroupBy(exprs ...string) *Query {
	if q.groupByExpr == nil {
		q.groupByExpr = make(map[string]string)
	}

	for _, expr := range exprs {
		if strings.TrimSpace(expr) == "" {
			continue
		}

		var field, alias string
		arr := strings.Split(expr, " as ")
		if len(arr) == 2 {
			field = strings.TrimSpace(arr[0])
			alias = strings.TrimSpace(arr[1])
		} else if len(arr) == 1 {
			field = strings.TrimSpace(arr[0])
			alias = ""
		}
		q.groupByExpr[field] = alias
	}
	return q
}

func (q *Query) GroupByKeys() string {
	var keys []string
	for field := range q.groupByExpr {
		keys = append(keys, field)
	}
	return strings.Join(keys, ", ")
}

func (q *Query) Agg(exprs ...string) *Query {
	q.aggExpr = q.append(q.aggExpr, exprs...)
	return q
}

func (q *Query) Having(exprs ...string) *Query {
	q.havingExpr = q.append(q.havingExpr, exprs...)
	return q
}

func (q *Query) OrderBy(exprs ...string) *Query {
	q.orderByExpr = q.append(q.orderByExpr, exprs...)
	return q
}

func (q *Query) String() (string, error) {
	// GroupBy() must be called explicitly while calling Agg()
	if len(q.aggExpr) > 0 && len(q.groupByExpr) == 0 {
		return "", errors.New("group by keys are expected when accept agg expr")
	}

	q.SelectAs(q.groupByExpr)
	q.Select(q.aggExpr...)

	// select expression is expected
	if len(q.selectExpr) == 0 {
		return "", errors.New("nothing selected")
	}

	template := fmt.Sprintf("\nselect %s \nfrom %s\n", q.Selected(), q.TableName)
	if len(q.filterExpr) > 0 {
		template += fmt.Sprintf("where (%s)\n", strings.Join(q.filterExpr, ")\n  and ("))
	}
	if len(q.groupByExpr) > 0 {
		template += fmt.Sprintf("group by %s\n", q.GroupByKeys())
	}
	if len(q.havingExpr) > 0 {
		template += fmt.Sprintf("having (%s)\n", strings.Join(q.havingExpr, ")\n   and ("))
	}
	if len(q.orderByExpr) > 0 {
		template += fmt.Sprintf("order by %s\n", strings.Join(q.orderByExpr, ", "))
	}

	return template, nil
}

func (q *Query) append(arr []string, elems ...string) []string {
	for _, elem := range elems {
		elem = strings.TrimSpace(elem)

		if elem != "" && !q.contains(arr, elem) {
			arr = append(arr, elem)
		}
	}

	return arr
}

func (q *Query) contains(arr []string, elem string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}
