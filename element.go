
package main

import (
	"fmt"
)

type Element struct {
	Data          *Dataset
	Id            int
	Name          string
	TypeId        int
	TypeName      string
	Weight        float64
	FormatedName  string
	BestName      string
}

func (self Element) String() string {
	format := `
		Element %d :
		BestName     | %s
		TypeId       | %d
		TypeName     | %s
		Weight       | %v
	`
	return fmt.Sprintf(
		StrJoin("\n", Apply(StrStrip, StrSplit(StrStrip(format), '\n'))),
		self.Id,
		self.BestName,
		self.TypeId,
		self.TypeName,
		self.Weight,
	)
}

func NewElement(row []string, types map[int]string) *Element {
	// fmt.Printf("E> %v\n", row)
	result := Element{}
	result.Id = PanicParseInt(row[0])
	result.Name = row[1]
	if len(result.Name) >= 2 {
		result.Name = StrCrop(row[1], 1, 1)
	}
	// println("------")
	// println(">"+row[0]+"<")
	// println(">"+row[1]+"<")
	// println(">"+row[2]+"<")
	// println(">"+row[3]+"<")
	// println("------")
	result.TypeId = PanicParseInt(row[2])
	result.TypeName = types[PanicParseInt(row[2])]
	if len(result.TypeName) >= 2 {
		result.TypeName = StrCrop(result.TypeName, 1, 1)
	}
	result.Weight = PanicParseFloat64(row[3])
	if len(row) == 5 {
		result.FormatedName = StrCrop(row[4], 1, 1)
		result.BestName = result.FormatedName
	} else {
		result.BestName = result.Name
		result.FormatedName = ""
	}
	return &result
}