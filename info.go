
package main

func unsortedResearch(src, rel, trg string) []*Relation {

	isNull := func(chain string)bool{return(chain=="...")}
	result := []*Relation{}

	postLog("%s %s %s\n%v %v %v", src, rel, trg, isNull(src), isNull(rel), isNull(trg))
	
	if !isNull(src) && !isNull(rel) && !isNull(trg) {
		ds, duration := AsTimedStep(func()*Dataset{return newDataset(src, rel)})
		postLog("Using search mode #%d\n", 0)
		postLog("Dataset loading time : %s\n", duration)
        return Filter(
        	func(relation *Relation) bool {
        		return relation.Target.Name == trg
        	},
        	ds.Relations,
        )
	}

	if !isNull(src) && !isNull(rel) &&  isNull(trg) {
		ds, duration := AsTimedStep(func()*Dataset{return newDataset(src, rel)})
		postLog("Using search mode #%d\n", 1)
		postLog("Dataset loading time : %s\n", duration)
		return ds.Relations
	}

	if !isNull(src) &&  isNull(rel) && !isNull(trg) {
		ds, duration := AsTimedStep(func()*Dataset{return newDataset(src, "")})
		postLog("Using search mode #%d\n", 2)
		postLog("Dataset loading time : %s\n", duration)
		return Filter(
        	func(relation *Relation) bool {
        		return relation.Target.Name == trg
        	},
			ds.Relations,
		)
	}

	if !isNull(src) &&  isNull(rel) &&  isNull(trg) {
		ds, duration := AsTimedStep(func()*Dataset{return newDataset(src, "")})
		postLog("Using search mode #%d\n", 3)
		postLog("Dataset loading time : %s\n", duration)
		return ds.Relations
	}

	if  isNull(src) && !isNull(rel) && !isNull(trg) {
		ds, duration := AsTimedStep(func()*Dataset{return newDataset(trg, "")})
		postLog("Using search mode #%d\n", 4)
		postLog("Dataset loading time : %s\n", duration)
		return Filter(
			func(relation *Relation) bool {
				return relation.Mode == INWARD_RELATION
			},
			ds.Relations,
		)
	}

	postLog("Research case not handled")
	return result

}

func research(src, rel, trg string) []*Relation {

	result := unsortedResearch(src, rel, trg)
	// postLog("using sorted research")
	// result = Sorted(func(r *Relation)float64{return r.Weight}, result)
	result = Filter(func(r *Relation)bool{return -1 != r.Weight && r.Weight != 1}, result)
	// result = Filter(func(r *Relation)bool{return -1 != r.NormedWeight && r.NormedWeight != 1}, result)
	// result = Filter(func(r *Relation)bool{return 0.2 <= r.NormedWeight && r.NormedWeight <= 0.8}, result)
	return result

}

func findRelations(src, trg string) []Relation {
    return AsStep("Checking", func()[]Relation {
        ds := AsStep("Opening ds", func()*Dataset {
            return newDataset(src, "")
        })
        result := []Relation{}
        Printfc[42]("|dsR| = %d\n", len(ds.Relations))
        for _, relation := range ds.Relations {
            if relation.Target.Name == trg {
            	result = append(result, *relation)
            }
        }
        Printfc[42]("|res| = %d\n", len(result))
        return result
    })
}

func checkExists(src, rel, trg string) bool {
    return AsStep("Checking", func()bool{
        ds := AsStep("Opening ds", func()*Dataset {
            return newDataset(src, rel)
        })
        for _, relation := range ds.Relations {
            if relation.Target.Name == trg {
                return true
            }
        }
        return false
    })
}