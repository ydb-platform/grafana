package bind

type Binder interface {
	Rebind(sql string, args ...any) (string, []any, error)
}
