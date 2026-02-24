package migrator

type DialectRecursiveCTE interface {
	Dialect

	RecursiveQueriesAreSupported() (bool, error)
}
