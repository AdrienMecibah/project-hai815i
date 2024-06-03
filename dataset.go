
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"html"
	"os"
)

type Dataset struct {
	Elements []*Element
	Relations []*Relation
}



var _NEWDATASET_MODE_JUMPER_SWITCH = false
var _NEWDATASET_MODE_JUMPER_RELIN = false
var _NEWDATASET_MODE_JUMPER_RELOUT = false

func newDataset(element string, relation string) *Dataset {
	if _NEWDATASET_MODE_JUMPER_SWITCH {
		postLog("Bypass !!!")
		return newDatasetFromElementAndRelation(element, relation, _NEWDATASET_MODE_JUMPER_RELOUT, 	_NEWDATASET_MODE_JUMPER_RELIN)
	}
	return newDatasetFromElementAndRelation(element, relation, true, true)
}

func newDatasetFromElementAndRelation(element string, relation string, relationOut bool, relationIn bool) *Dataset {
	url := fmt.Sprintf("https://www.jeuxdemots.org/rezo-dump.php?gotermsubmit=Chercher&gotermrel=%s&rel=%s", element, relation)
	if !relationOut {
		url += "&relout=norelout"
	}
	if !relationIn {
		url += "&relin=norelin"
	}
	return newDatasetFromURL(url)
}

func newDatasetFromURL(url string) *Dataset {
	// encoded_url := StrJoin("-", Apply(func(r rune) string { return string([]rune(fmt.Sprintf("%U", r))[2:]) }, []rune(url)[41:]))
	encoded_url := StrJoin("", Apply(func(r rune)string{ 
		if IsIn(r, []rune(`\/:*?"<>|`)) {
			return string([]rune(fmt.Sprintf("%U", r))[0:])
		} else {
			return string(r)
		}},
		[]rune(url)[40:],
	))
	// encoded_url := StrJoin("-", Apply(func(r rune) string { return string([]rune(fmt.Sprintf("%U", r))[2:]) }, []rune(url)[41:]))
	for _, name := range listDir("history") {
		if name == "history"+string(os.PathSeparator)+encoded_url {
			return newDatasetFromContent(ReadFile(name))
		}
	}
	const startMark = "</def>"
	const stopMark = "</CODE>"
	response, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("Error making GET request : %v", err))
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(fmt.Sprintf("Error reading response body : %v", err))
	}
	text := string(body)
	text = html.UnescapeString(text)
	content := []rune(text)
	if SliceSubIndex(content, []rune(startMark)) == -1 {
		println("Can not read HTML")
		println(text)
		println("Can not read HTML")
	}
	content = content[SliceSubIndex(content, []rune(startMark))+len([]rune(startMark)):]
	// fmt.Printf("\x1b[38;5;202m%s\x1b[39m\n", string(content))
	content = content[:SliceSubIndex(content, []rune(stopMark))]
	// fmt.Printf("\x1b[38;5;226m%s\x1b[39m\n", string(content))
	newPath := ".\\history"+string(os.PathSeparator)+encoded_url
	postLog("storing in "+newPath)
	WriteFile(newPath, string(content))
	return newDatasetFromContent(string(content))
}

func newDatasetFromContent(content string) *Dataset {
	grid := Apply(
		func(line string) []string { // séparer les éléments séparés par les ";" en prenant en compte les apostrophes
			return StrSplit(line, ';')
			result := []string{""}
			lit := false
			for i, r := range []rune(line) {
				if r == ';' && !lit {
					result = append(result, "")
				} else if r == '\'' {
					if i > 0 && []rune(line)[i-1] == '\\' {
						if i > 1 && []rune(line)[i-2] == '\\' {
							panic("Reading litteral back slashes is not implemented")
						}
						result[len(result)-1] += string(r)
					} else {
						lit = !lit
					}
				} else {
					result[len(result)-1] += string(r)
				}
			}
			return result
		},
		Filter( // filter les lignes vides
			func(line string) bool {
				content := []rune(StrStrip(line))
				if len(content) == 0 || (len(content) >= 2 && content[0] == '/' && content[1] == '/') {
					return false
				}
				return true
			},
			StrSplit(content, '\n'),
		),
	)
	nodeTypes := Map(
		func(line []string) (int, string) {
			return PanicParseInt(line[1]), line[2]
		},
		Filter(
			func(line []string) bool {
				return line[0] == "nt"
			},
			grid,
		),
	)
	relationTypes := Map(
		func(line []string) (int, string) {
			return PanicParseInt(line[1]), line[2]
		},
		Filter(
			func(line []string) bool {
				return line[0] == "rt"
			},
			grid,
		),
	)
	elements := Apply(
		func(line []string) *Element {
			// fmt.Printf("\nNew[%d]%v\n", len(line), line)
			return NewElement(line[1:], nodeTypes)
		},
		Filter(
			func(block []string) bool {
				return block[0] == "e"
			},
			grid,
		),
	)
	i := -1
	relations := Apply(
		func(line []string) *Relation {
			i += 1
			return NewRelation(line[1:], relationTypes,i)
		},
		Filter(
			func(block []string) bool {
				return block[0] == "r"
			},
			grid,
		),
	)
	for _, r := range relations {
		for _, e := range elements {
			if r.sourceId == e.Id {
				r.Source = e
			}
			if r.destinationId == e.Id {
				r.Target = e
			}
		}
		if r.Source == nil {
			panic(fmt.Sprintf("no node with id %d", r.Id))
		}
		if r.Target == nil {
			panic(fmt.Sprintf("no node with id %d", r.Id))
		}
	}
	result := Dataset{}
	result.Elements = elements
	result.Relations = relations
	return &result
}