package migrator

func NewBaseDialect(
	driverName string,
	dialect Dialect,
) BaseDialect {
	return BaseDialect{
		dialect:    dialect,
		driverName: driverName,
	}
}
