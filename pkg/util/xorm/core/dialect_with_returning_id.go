package core

type DialectWithReturningID interface {
	Dialect

	WithReturningID()
}
