package searchstore

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
)

// Builder defaults to returning a SQL query to get a list of all dashboards
// in default order, but can be modified by applying filters.
type Builder struct {
	// List of FilterWhere/FilterGroupBy/FilterOrderBy/FilterLeftJoin
	// to modify the query.
	Filters []interface{}
	Dialect migrator.Dialect

	params []interface{}
	sql    bytes.Buffer
}

// ToSQL builds the SQL query and returns it as a string, together with the SQL parameters.
func (b *Builder) ToSQL(limit, page int64) (string, []interface{}) {
	b.params = make([]interface{}, 0)
	b.sql = bytes.Buffer{}

	b.buildSelect()

	b.sql.WriteString("( ")
	orderQuery := b.applyFilters()

	b.sql.WriteString(b.Dialect.LimitOffset(limit, (page-1)*limit) + `) AS ids
		INNER JOIN dashboard ON ids.id = dashboard.id`)
	b.sql.WriteString("\n")

	b.sql.WriteString(
		`LEFT OUTER JOIN dashboard AS folder ON folder.id = dashboard.folder_id
		LEFT OUTER JOIN dashboard_tag ON dashboard.id = dashboard_tag.dashboard_id`)
	b.sql.WriteString("\n")
	b.sql.WriteString(orderQuery)

	return b.sql.String(), b.params
}

func (b *Builder) buildSelect() {
	var recQuery string
	var recQueryParams []interface{}

	selectList := `SELECT
			dashboard.id,
			dashboard.uid,
			dashboard.title,
			dashboard.slug,
			dashboard_tag.term,
			dashboard.is_folder,
			dashboard.folder_id,
			folder.uid AS folder_uid,
			folder.slug AS folder_slug,
			folder.title AS folder_title
		`
	if b.Dialect.DriverName() == migrator.YDB {
		selectList += `, dashboard.title AS title `
	}
	b.sql.WriteString(selectList)
	for _, f := range b.Filters {
		if f, ok := f.(FilterSelect); ok {
			b.sql.WriteString(fmt.Sprintf(", %s", f.Select()))
		}

		if f, ok := f.(FilterWith); ok {
			recQuery, recQueryParams = f.With()
		}
	}

	b.sql.WriteString(` FROM `)

	if recQuery == "" {
		return
	}

	// prepend recursive queries
	var bf bytes.Buffer
	bf.WriteString(recQuery)
	bf.WriteString(b.sql.String())

	b.sql = bf
	b.params = append(recQueryParams, b.params...)
}

func (b *Builder) applyFilters() (ordering string) {
	joins := []string{}
	orderJoins := []string{}

	wheres := []string{}
	whereParams := []interface{}{}

	groups := []string{}
	groupParams := []interface{}{}

	orders := []string{}

	for _, f := range b.Filters {
		if f, ok := f.(FilterLeftJoin); ok {
			s := f.LeftJoin()
			if s != "" {
				joins = append(joins, fmt.Sprintf(" LEFT OUTER JOIN %s ", s))
			}
		}

		if f, ok := f.(FilterWhere); ok {
			sql, params := f.Where()
			if sql != "" {
				wheres = append(wheres, sql)
				whereParams = append(whereParams, params...)
			}
		}

		if f, ok := f.(FilterGroupBy); ok {
			sql, params := f.GroupBy()
			if sql != "" {
				groups = append(groups, sql)
				groupParams = append(groupParams, params...)
			}
		}

		if f, ok := f.(FilterOrderBy); ok {
			if f, ok := f.(FilterLeftJoin); ok {
				orderJoins = append(orderJoins, fmt.Sprintf(" LEFT OUTER JOIN %s ", f.LeftJoin()))
			}
			orders = append(orders, f.OrderBy())
		}
	}

	forceIndex := ""
	if b.Dialect.DriverName() == migrator.MySQL {
		forceIndex = " FORCE INDEX (IDX_dashboard_title) "
	}

	// YDB: subquery ORDER BY runs over the subquery result; we always add default title sort below if empty, so include title
	if len(orders) < 1 {
		orders = append(orders, TitleSorter{}.OrderBy())
	}
	subquerySelect := "SELECT dashboard.id AS id FROM dashboard"
	if b.Dialect.DriverName() == migrator.YDB {
		for _, o := range orders {
			if strings.Contains(o, "dashboard.title") {
				subquerySelect = "SELECT dashboard.id AS id, dashboard.title AS title FROM dashboard"
				break
			}
		}
	}
	b.sql.WriteString(fmt.Sprintf("%s %s", subquerySelect, forceIndex))
	b.sql.WriteString(strings.Join(joins, ""))

	if len(wheres) > 0 {
		b.sql.WriteString(fmt.Sprintf(" WHERE %s", strings.Join(wheres, " AND ")))
		b.params = append(b.params, whereParams...)
	}

	if len(groups) > 0 {
		cols := make([]string, 0, len(orders)+len(groups))
		for _, o := range orders {
			o := strings.TrimSuffix(o, " DESC")
			o = strings.TrimSuffix(o, " ASC")
			exists := false
			for _, g := range groups {
				if g == o {
					exists = true
					break
				}
			}
			if !exists {
				cols = append(cols, o)
			}
		}
		cols = append(cols, groups...)
		b.sql.WriteString(fmt.Sprintf(" GROUP BY %s", strings.Join(cols, ", ")))
		b.params = append(b.params, groupParams...)
	}

	orderByCols := []string{}
	for _, o := range orders {
		orderByCols = append(orderByCols, b.Dialect.OrderBy(o))
	}

	orderBy := fmt.Sprintf(" ORDER BY %s", strings.Join(orderByCols, ", "))
	b.sql.WriteString(orderBy)

	// YDB: outer query has both dashboard.title and alias "title"; ORDER BY title is ambiguous ("column name: title conflicted")
	orderByOuter := orderBy
	if b.Dialect.DriverName() == migrator.YDB {
		orderByOuter = strings.ReplaceAll(orderByOuter, " title ", " dashboard.title ")
	}
	order := strings.Join(orderJoins, "")
	order += orderByOuter
	return order
}
