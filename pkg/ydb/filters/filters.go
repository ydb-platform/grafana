package filters

import "github.com/grafana/grafana/pkg/util/xorm/core"

var Filters = []core.Filter{
	newPatches( // first conversion by lookup
		&AddColumn{},
		&RenameColumn{},
		&CreateIndex{},
		&DropIndex{},
		&RemoveOrderByFromUpdate{},
		&Replace{},
		&Lower{},
		&ReduceDuplicateIdInSelect{},
		&ConvertILikeToLikeLowerCase{},
		&ConvertLikeToStartsWithEndsWith{},
		&ConvertSubstrToSubstring{},
		&ReplaceIdPlaceholderToTableId{},
	),
	//&Quote{},
	//&ConcatOrgID{},
	&EnableCostBasedOptimizer{},
	//&ExtractLiteralToArgs{},
	&ConvertNumbersToInt64{},
	//&ConvertBoolToUint8{},
	&ConvertStringToDatetime64{},
	&ConvertTimeToDatetime64{},
	&ConvertDurationToInt64{},
	&ConvertLimitArgToUint64{},
	&ConvertInArgsToList{},
	&ConvertPositionalArgsToYdbNamedParameters{}, // must be a final conversion
}

func isWordByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
