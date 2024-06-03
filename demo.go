
package main

import "time"

type Triplet struct {
	source string
	relation string
	target string
} 

type Inference struct {
	causes []Triplet
	consequences []Triplet
}

func (self Inference) String() string {
	result := ""
	for i, cause := range self.causes {
		if i == 0 {
			result += "  "
		} else {
			result += "∧ "
		}
		result += Sprintf("%v\n", cause)
	}
	if len(self.consequences) == 1 {
		result += Sprintf("→ %v", self.consequences[0])
	} else {
		result += "→ \n"
		for consequence := range self.consequences {
			result += Sprintf("%v\n", consequence)
		}
	} 
	return result
}

type Justification struct {
	causes []Justification
	consequence Triplet
}

func (self Justification) repr(n int) string {
	indent := "    "
	result := ""
	if len(self.causes) == 0 {
		result += Sprintf("%sFACT --> %v", StrMul(indent, n), self.consequence)
	} else {
		for _, cause := range self.causes {
			result += Sprintf("%s", cause.repr(n+1)) + "\n"
		}
		result += Sprintf("%s --> %v", StrMul(indent, n), self.consequence)
	}
	return result
}

func (self Justification) String() string {
	return self.repr(0)
}

var _c = 0
func g(goal Triplet, inferences []Inference) (bool, string) {
	applyMapping := func(mapping map[string]string, inference *Inference) []string {
		variables := []string{}
		for i := range inference.causes {
			if mapping[inference.causes[i].source] == goal.source {
				inference.causes[i].source = goal.source
			} else {
				inference.causes[i].source = "$" + inference.causes[i].source
				variables = append(variables, "$" + inference.causes[i].source) 
			}
			if mapping[inference.causes[i].target] == goal.target {
				inference.causes[i].target = goal.target
			} else {
				inference.causes[i].target = "$" + inference.causes[i].target
				variables = append(variables, "$" + inference.causes[i].target) 
			}
		}
		for i, variable := range variables {
			variables[i] = StrReplace(variable, "$$", "$")
		}
		return RemoveDuplicates(variables)
	}
	_c += 1
	// Printfc[42]("%d %v \n", _c, goal)
	doResearch := func(a, b, c string) []*Relation {
		result, t := AsTimedStep(func()[]*Relation{ 
			result := research(a, b, c)
			result = Sorted(func(r *Relation)float64{return -r.Weight}, result)
			return result
		})
		loadings = append(loadings, t)

		return result
	}
	for _, inference := range inferences {
		for csq, consequence := range inference.consequences {
			if consequence.relation == goal.relation {
				mapping := map[string]string{}
				mapping[consequence.source] = goal.source
				mapping[consequence.target] = goal.target
				variables := applyMapping(mapping, &inference)
				// Printf("%v\n", variables)
				// var candidates map[string]map[int][]string{}
				candidates := Map(
					func(candidate string) (string, map[int][]string) {
						return candidate, Map(
							func(i int) (int, []string) {
								return i, []string{}
							},
							Range[int](len(inference.causes)),
						)
					},
					variables,
				)
				for c, cause := range inference.causes {
					// Printfc[38]("%d : %v\n", c, cause)
					if IsIn(cause.source, variables) && IsIn(cause.target, variables) {
						search := doResearch("...", cause.relation, "...")

						candidates[cause.source][c] = append(candidates[cause.source][c], Apply(
							func(r *Relation) string {
								return "?"+r.Source.Name
							}, 
							search,
						)...)
						candidates[cause.target][c] = append(candidates[cause.target][c], Apply(
							func(r *Relation) string {
								return "?"+r.Target.Name
							}, 
							search,
						)...)
					} else if IsIn(cause.source, variables) {
						candidates[cause.source][c] = append(candidates[cause.source][c], Apply(
							func(r *Relation) string {
								return "?"+r.Source.Name
							}, 
							doResearch("...", cause.relation, cause.target),
						)...)
					} else if IsIn(cause.target, variables) {
						candidates[cause.target][c] = append(candidates[cause.target][c], Apply(
							func(r *Relation) string {
								return "?"+r.Target.Name
							}, 
							doResearch(cause.source, cause.relation, "..."),
						)...)
					}
				}
				for _, possibilites := range candidates {
					for i, words := range possibilites {
						possibilites[i] = RemoveDuplicates(words)
					}
				}
				// Printfc[202]("variables : \n");
				// for variable, possibilites := range candidates {
				// 	Printfc[202]("  %s : \n", variable)
				// 	for c, words := range possibilites {
				// 		Printfc[202]("  %d : %d \n", c, len(words))
				// 		if len(words) < 50 {
				// 			Printfc[202]("    %v\n", words)
				// 		}
				// 	}
				// }
				elections := map[string][]string{}
				results := map[string]string{}
				for variable, possibilites := range candidates {
					elections[variable] = []string{}
					values := [][]string{}
					for _, words := range possibilites {
						values = append(values, words)
					}
					for i, words := range values {
						for _, word := range words {
							found := false
							for j := range values {
								if i == j {
									continue
								}
								for _, w := range values[j] {
									if w == word {
										found = true
										break
									}
								}
								if found {
									break
								}
							}
							if found {
								elections[variable] = append(elections[variable], word)
							}
						}
					}
				}
				getOut := false
				for variable, words := range elections {
					if len(words) == 0 {
						getOut = true
					} else {
						results[variable] = words[0]
					}
				}
				if getOut {
					break
				}
				// Printfc[9]("%v\n", inference)
				for i, cause := range inference.causes {
					if IsIn(cause.source, variables) {
						inference.causes[i].source = results[cause.source]
					}
					if IsIn(cause.target, variables) {
						inference.causes[i].target = results[cause.target]
					}
				}
				for i, consequence := range inference.consequences {
					if IsIn(consequence.source, variables) {
						inference.consequences[i].source = results[consequence.source]
					}
					if IsIn(consequence.target, variables) {
						inference.consequences[i].target = results[consequence.target]
					}
				}
				inference.consequences[csq].source = goal.source
				inference.consequences[csq].target = goal.target
				return true, Sprintf("%v", inference)
			} 
		}
	}
	return false, ""
}

var loadings []time.Duration

var _ = script("demo", func(){

	s := GetArgs(0)
	r := GetArgs(1)
	t := GetArgs(2)

	inferences := []Inference{}
	transitiveRelations := []string {
	    "r_isa",      
	    "r_hypo",     
	    "r_domain",   
	    "r_holo",     
	    "r_has_part", 
	    "r_lieu",     
	    "r_lieu-1",
	}
	for _, relation := range transitiveRelations {
		inferences = append(inferences, Inference{
			causes: []Triplet {
				Triplet{"A", relation, "B"},
				Triplet{"B", relation, "C"},
			},
			consequences: []Triplet {
				Triplet{"A", relation, "C"},
			},
		})
	}
	synonymCompatibleRelations := []string {
		"r_associated",
		"r_raff_sem",
		"r_raff_morpho",
		"r_domain",
		"r_pos",
		"r_isa",
		"r_anto",
		"r_hypo",
		"r_has_part",
		"r_holo",
	}
	for _, relation := range synonymCompatibleRelations {
		inferences = append(inferences, Inference{
			causes: []Triplet {
				Triplet{"A", "r_syn", "B"},
				Triplet{"B", relation, "C"},
			},
			consequences: []Triplet {
				Triplet{"A", relation, "C"},
			},
		})
	}
	if !GetSwitch("no-square-syn") {
		inferences = append(
			[]Inference {
				Inference {
					causes: []Triplet {
						Triplet{"A", r, "B"},
						Triplet{"C", "r_syn", "A"},
						Triplet{"D", "r_syn", "B"},
					},
					consequences: []Triplet {
						Triplet{"C", r, "D"},
					},
				},
			},
			inferences...,
		)
	}
	if !GetSwitch("no-isa-based") {
		inferences = append(
			[]Inference {
				Inference {
					causes: []Triplet {
						Triplet{"A", "r_isa", "B"},
						Triplet{"B", r, "C"},
					},
					consequences: []Triplet {
						Triplet{"A", r, "C"},
					},
				},
			},
			inferences...,
		)
	}

	var res bool
	var justification string
	_, duration := AsTimedStep(func()int{
		res, justification = g(Triplet{s, r, t}, inferences)
		return 0
	})
	for _, d := range loadings {
		duration = duration - (d)
	}
	Printf("\x1b[4m%v\x1b[0m\n", res)
	Printfc[226]("%s\n", justification)
	Printfc[42]("Inference duration : %v\nLoading durations : %v\n", duration, loadings)

})