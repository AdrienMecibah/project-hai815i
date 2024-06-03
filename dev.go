
package main

const DEFAULT_DEV_ELEMENT = "chat"
const DEFAULT_DEV_RELATION = "r_agent-1"

const DEFAULT_DEMO_SUBJECT = "chat"
const DEFAULT_DEMO_RELATION = "r_syn"
const DEFAULT_DEMO_OBJECT = "chat"

func demo() {
    subjectName := GetOptions("s", DEFAULT_DEMO_SUBJECT)
    relationName := GetOptions("r", DEFAULT_DEMO_RELATION)
    objectName := GetOptions("o", DEFAULT_DEMO_OBJECT)
    
    dsHyperOfSubject := AsStep("Opening dsHyperOfSubject", func()*Dataset {
        return newDataset(subjectName, "r_isa")
    })
    Step("Infering stuff", func(){
        for _, relation := range dsHyperOfSubject.Relations[:10] {
            // dsTarget := newDataset(relation.Target.)
            print(relation)
        }
    })
    println(len(dsHyperOfSubject.Relations))
    println(relationName, objectName)
}



func dev() {
    elementName := GetOptions("e", DEFAULT_DEV_ELEMENT)
    relationName := GetOptions("r", DEFAULT_DEV_RELATION)
    ds := AsStep("Opening ds", func()*Dataset {
        return newDataset(elementName, relationName)
    })
    print(len(ds.Relations), "->")
    ds.Relations = Filter(func(r *Relation)bool{return r.TypeName==relationName}, ds.Relations)
    println(len(ds.Relations))
    Step("Infering stuff", func(){
        mainElementCandidates := IFilter(
            func(i int, e *Element) bool {
                return elementName == e.BestName
            },
            ds.Elements,
        )
        if len(mainElementCandidates) != 1 {
            panic(Sprintf("too many or not enough elements : %v", mainElementCandidates))
        }
        mainElement := mainElementCandidates[0]
        for i, r := range ds.Relations {
            if r.Weight > 0 { 
                continue
            }
            color := -1
            if r.Source.Id == mainElement.Id {
                color = 38
            } else if r.Target.Id == mainElement.Id {
                color = 42
            }
            Printfc[color]("%d %s \n\n", i, r)
        }
    })
}