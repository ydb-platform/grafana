package filters

import "github.com/grafana/grafana/pkg/util/xorm/core"

var Filters = []core.Filter{
	&AddColumn{},
	&RenameColumn{},
	&CreateIndex{},
	&DropIndex{},
	newPatches(),
	&ID{},
	&Quote{},
	&Replace{},
	&Lower{},
	&DuplicateID{},
	&ConcatOrgID{},
	&Substr{},
	&ILike{},
	&CostBasedOptimizer{},
	&Literal{},
	&Numbers{},
	&Bool{},
	&String{},
	&Duration{},
	&Limit{},
	&IN{},
	&Args{},
}

func isWordByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
