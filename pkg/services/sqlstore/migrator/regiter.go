package migrator

func RegisterDialect(driverName string, fn dialectFunc) {
	if _, exist := supportedDialects[driverName]; exist {
		panic("Database already registered: " + driverName)
	}
	supportedDialects[driverName] = fn
}

func HasDialect(driverName string) bool {
	_, exist := supportedDialects[driverName]
	return exist
}
