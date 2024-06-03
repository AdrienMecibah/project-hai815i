 
package main

import (
	"fmt"
)

const INWARD_RELATION = "inward"
const IR = INWARD_RELATION
const OUTWARD_RELATION = "outward"
const OR = OUTWARD_RELATION

type Relation struct {
	Data          *Dataset
	Id             int
	sourceId       int
	destinationId  int
	Source        *Element
	Target        *Element
	TypeId         int
	TypeName       string
	Mode           string
	Weight         float64
	NormedWeight   float64
	Rank           int
}

func (self Relation) String() string {
	format := `
		Relation %d :
		Mode         | %s
		SourceName   | %s
		TargetName   | %s
		TypeName     | %s
		SourceId     | %d
		TargetId     | %d
		TypeId       | %d
		Weight       | %v
		NormedWeight | %v
		Rank         | %d
	`
	return fmt.Sprintf(
		StrJoin("\n", Apply(StrStrip, StrSplit(StrStrip(format), '\n'))),
		self.Id,
		self.Mode,
		self.Source.BestName,
		self.Target.BestName,
		self.TypeName,
		self.Source.Id,
		self.Target.Id,
		self.TypeId,
		self.Weight,
		self.NormedWeight,
		self.Rank,
	)
}

func NewRelation(row []string, types map[int]string, debug ...any) *Relation {
	// fmt.Printf("R> %v %v\n", debug[0], row)
	result := Relation{}
	result.Id = PanicParseInt(row[0])
	result.sourceId = PanicParseInt(row[1])
	result.destinationId = PanicParseInt(row[2])
	result.TypeId = PanicParseInt(row[3])
	result.TypeName = StrCrop(types[result.TypeId], 1, 1)
	result.Weight = PanicParseFloat64(row[4])
	// println(len(row) == 7, StrStrip(row[5]) != "-", StrStrip(row[6]) != "-")
	// 	for i, l := range row {
	// 		println(i, ">"+(l)+"< | >"+StrStrip(l)+"<")
	// 	}
	if len(row) == 7 && StrStrip(row[5]) != "-" && StrStrip(row[6]) != "-" {
		result.Mode = OR
		result.NormedWeight = PanicParseFloat64(row[5])
		result.Rank = PanicParseInt(row[6])
	} else if len(row) == 5 || len(row) == 7 && StrStrip(row[5]) == "-" && StrStrip(row[6]) == "-" {
		result.Mode = IR
		result.NormedWeight = -1
		result.Rank = -1
	} else {
		for i, l := range row {
			println(i, ">"+(l)+"<")
		}
		panic(fmt.Sprintf("unknow line model of length : %d\n", len(row)))
	}
	return &result
}

