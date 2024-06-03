
package main

var _ = script("find-relations", func() {
	for _ , rel := range findRelations(GetArgs(0), GetArgs(1)) {
		Printf("%s\n\n", rel)
	}
})

var _ = script("check", func() {
	println(checkExists(GetArgs(0), GetArgs(1), GetArgs(2)))
})

var _ = script("search", func() {
	// ? -script search <source> <relation> <target> [--positive-weight|--zero-weight|--negative-weight] [-mode inward|outward] [-lim <count>]
	relations := Sorted(func(r *Relation)float64{return -r.NormedWeight}, research(GetArgs(0), GetArgs(1), GetArgs(2)))
	if GetSwitch("positive-weight") {
		relations = Filter(
			func(relation *Relation) bool {
				return relation.Weight > 0
			},
			relations,
		)
	}
	if GetSwitch("zero-weight") {
		relations = Filter(
			func(relation *Relation) bool {
				return relation.Weight == 0
			},
			relations,
		)
	}
	if GetSwitch("negative-weight") {
		relations = Filter(
			func(relation *Relation) bool {
				return relation.Weight < 0
			},
			relations,
		)
	}
	if GetOptions("mode", "") != "" {
		relations = Filter(
			func(relation *Relation) bool {
				return relation.Mode == GetOptions("mode", "")
			},
			relations,
		)
	}
	if lim := GetOptions("lim", ""); lim != "" && PanicParseInt(lim) < len(relations) {
		relations = relations[:PanicParseInt(lim)]
	}
	for i, relation := range relations {
		Printf("\x1b[38;5;38m[%d]\x1b[39m %s \n", i+1, StrReplace(Sprintf("%s", relation), "\n", "\n"+StrMul(" ", 3+len(Sprintf("%d", i+1)))))
	}
	if len(relations) == 0 {
		Printfc[9]("no relations\n")
	}
})
