package bind

var _ Binder = Prepend("")

type Prepend string

func (p Prepend) Rebind(sql string, args ...any) (string, []any, error) {
	return string(p) + ";\n\n" + sql, args, nil
}
