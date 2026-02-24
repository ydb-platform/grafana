package core

// FilterWithArgs is an interface to filter SQL with args
type FilterWithArgs interface {
	DoWithArgs(sql string, dialect Dialect, table *Table, args ...any) (string, []any)
}

func QueryPreprocess(
	dialect Dialect,
	table *Table,
	sql string,
	args ...any,
) (string, []any) {
	for _, filter := range dialect.Filters() {
		if filterWithArgs, has := filter.(FilterWithArgs); has {
			sql, args = filterWithArgs.DoWithArgs(sql, dialect, table, args...)
		} else {
			sql = filter.Do(sql, dialect, table)
		}
	}

	return sql, args
}
