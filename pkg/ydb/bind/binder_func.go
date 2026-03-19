package bind

var _ Binder = Func(nil)

type Func func(sql string, args ...any) (string, []any, error)

func (f Func) Rebind(sql string, args ...any) (string, []any, error) {
	return f(sql, args...)
}
