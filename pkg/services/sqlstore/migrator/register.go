package migrator

func Register(driverName string, fn dialectFunc) {
	if _, exist := supportedDialects[driverName]; exist {
		panic("Already register database type: " + driverName)
	}

	supportedDialects[driverName] = fn
}
