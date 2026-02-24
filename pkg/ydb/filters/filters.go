package filters

import "github.com/grafana/grafana/pkg/util/xorm/core"

var Filters = []core.Filter{
	newPatches(),
	&ID{},
	&Quote{},
	&Replace{},
	&Lower{},
	&DuplicateID{},
	&ConcatOrgID{},
	&BoolIsColumn{},
	&Substr{},
	&ILike{},
	&CostBasedOptimizer{},
	&Numbers{},
	&Duration{},
	&Limit{},
	&IN{},
	&Args{},
}
