
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
	"strings"
    "path/filepath"
)


func IfThenElse[T any](condition bool, positive T, negative T) T {
	if condition {
		return positive
	} else {
		return negative
	}
}
var Printfc map[int]func(...any) = Map(
    func(i int) (int, func(...any)) {
        return i, func(args ...any) {
            a := make([]any, len(args))
            format, _ := args[0].(string)
            for i, v := range args {
                a[i] = v
            }            
            a[0] = any(Sprintf("\x1b[38;5;%dm", i)+format+"\x1b[39m")
            Printf(a...)
        }
    },
    Range[int](256),
)

func RemoveDuplicates[T comparable](array []T) []T {
	result := make([]T, len(array))
	length := 0
	for _, object := range array {
		if !IsIn(object, result) {
			result[length] = object
			length += 1
		}
	}
	return result[:length]
}

func Range[T any](object any) []int {
    var n int
    if m, ok := object.(int); ok {
        n = m
    } else if iterable, ok := object.([]T); ok {
        n = len(iterable)
    }
    result := make([]int, n)
    for i:=0; i<n; i++ {
        result[i] = i
    }
    return result
}

func Sprintf(args ...any) string {
	if len(args) == 0 {
		panic("Sprintf shortcut must take at least one argument")
	}
	format, ok := args[0].(string)
	if !ok && len(args) == 1 {
		fmt.Printf("%v\n", args)
	}
	values := args[1:]
	return fmt.Sprintf(format, values...)
}

func Printf(args ...any) {
	if len(args) == 0 {
		panic("Printf shortcut must take at least one argument")
	}
	format, ok := args[0].(string)
	if !ok && len(args) == 1 {
		fmt.Printf("%v\n", args)
	}
	values := args[1:]
	fmt.Printf(format, values...)
}

func SliceSubIndex[T comparable](elements []T,  subelements[]T) int {
	for i := range elements {
		if SliceEq(elements[i:i+len(subelements)], subelements) {
			return i
		}
	}
	return -1
}

func SliceEq[T comparable](l1 []T, l2 []T) bool {
	if len(l1) != len(l2) {
		return false
	}
	for i := range l1 {
		if l1[i] != l2[i] {
			return false
		}
	}
	return true
}

func IsIn[T comparable](x T, list []T) bool {
	for _, obj := range list {
		if obj == x {
			return true
		}
	}
	return false
}

func Step(name string, code func()()) {
	AsStep(name, func()int{code(); return 0})
}

func AsStep[T any](name string, code func()T) T {
	start := time.Now()
	result := code()
	duration := time.Now().Sub(start)
	fmt.Printf("%s : %s\n", name, duration)
	return result
}

func AsTimedStep[T any](code func()T) (T, time.Duration) {
	start := time.Now()
	result := code()
	duration := time.Now().Sub(start)
	return result, duration
} 

func StrIndex(str string, substr string) int {
	content := []rune(str)
	subContent := []rune(substr)
	for i := range content {
		if SliceEq(content[i:i+len(subContent)], subContent) {
			return i
		}
	}
	return -1
}

func StrStrip(line string) string {
	SKIPABLES := []rune{' ', '\r', '\n', '\t'}
	content := []rune(line)
	for len(content) > 0 && IsIn(content[0], SKIPABLES) {
		content = content[1:]
	}
	for len(content) > 0 && IsIn(content[len(content)-1], SKIPABLES) {
		content = content[:len(content)-1]
	}
	return string(content)
}

func Filter[T any](mapper func(T)bool, iterable []T) []T {
    result := make([]T, 0)
    for _, x := range iterable {
        if mapper(x) {
            result = append(result, x)
        }
    }
    return result
}

func IFilter[T any](mapper func(int, T)bool, iterable []T) []T {
    result := make([]T, 0)
    for i, x := range iterable {
        if mapper(i, x) {
            result = append(result, x)
        }
    }
    return result
}

func Apply[T1, T2 any](mapper func(T1)T2, iterable []T1) []T2 {
    result := make([]T2, len(iterable))
    for i, x := range iterable {
        result[i] = mapper(x)
    }
    return result
}

func IApply[T1, T2 any](mapper func(int, T1)T2, iterable []T1) []T2 {
    result := make([]T2, len(iterable))
    for i, x := range iterable {
        result[i] = mapper(i, x)
    }
    return result
}

func _EMap[X any, K comparable, V any](mapper func(int, X)(K, V), iterable []X) map[K]V {
    result := map[K]V{}
    for i, x := range iterable {
        k, v := mapper(i, x)
        result[k] = v
    }
    return result
}

func KMap[K comparable, V any](mapper func(K)V, iterable []K) map[K]V {
    result := map[K]V{}
    for _, k := range iterable {
        v := mapper(k)
        result[k] = v
    }
    return result
}


func Map[X any, K comparable, V any](mapper func(X)(K, V), iterable []X) map[K]V {
    result := map[K]V{}
    for _, x := range iterable {
        k, v := mapper(x)
        result[k] = v
    }
    return result
}

func Sorted[T any, X ~int|~float64](sorter func(T)X, array []T) []T {
	if len(array) == 0 {
		return array
	}
	scores := make([]struct{index int; score X}, len(array))
	var min X
	for i := range array {
		s := sorter(array[i])
		if i == 0 || s < min {
			min = s
		}
		scores[i] = struct{index int; score X}{i, s}
	}
	for i := range scores[1:] {
		if scores[i].score > scores[i+1].score {
			sc := scores[i+1]
			scores[i+1] = scores[i]
			scores[i] = sc
		}
	}
	result := make([]T, len(array))
	for i, score := range scores {
		result[i] = array[score.index]
	}
	return result
}

func All(values []bool) bool {
	for _, value := range values {
		if !value {
			return false
		}
	}
	return true
}

func Any(values []bool) bool {
	for _, value := range values {
		if value {
			return true
		}
	}
	return false
}

func StrJoin(join string, lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	if len(lines) == 1 {
		return lines[0]
	}
	result := lines[0]
	for _, line := range lines[1:] {
		result += join + line
	}	
	return result
}

func GetArgs(index int) string {
	args := []string{}
	for i:=1; i<len(os.Args); i++ {
		runes := []rune(os.Args[i])
		if len(runes) > 1 && runes[0] == '-' && runes[1] == '-' {
			// ...
		} else if len(runes) > 0 && runes[0] == '-' {
			i += 1
		} else {
			args = append(args, os.Args[i])
		}
	}
	if index >= len(args) {
		panic(Sprintf("there are not %d arguments, only %d : %v", index, len(args), args))
	}
	return args[index]
}

func GetOptions(name string, defaultValue string) string {
	if len(os.Args) == 1 {
		return defaultValue
	}
	for i, arg := range os.Args {
		if i == 0 || i == len(os.Args) - 1 {
			continue
		}
		if arg == "-" + name {
			return os.Args[i+1]
		}
	}
	return defaultValue
}

func GetEqOptions(name string, defaultValue string) string {
	if len(os.Args) == 1 {
		return defaultValue
	}
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if IsIn('=', []rune(name)) {
			chunks := StrSplit(name, '=')
			arg = chunks[0]
			val := chunks[1]
			if arg == "--" + name {
				return val
			}
		}
	}
	return defaultValue
}

func GetSwitch(name string, defaultValues ...bool) bool {
	if len(defaultValues) > 1 {
		panic(Sprintf("agument defaultValues of GetSwitch is optional and can only be given one value, not %d", len(defaultValues)))
	}
	defaultValue := false
	if len(defaultValues) == 1 {
		defaultValue = defaultValues[0]
	}
	if len(os.Args) == 1 {
		return defaultValue
	}
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if arg == "--" + name {
			return true
		}
	}
	return defaultValue
}

func StrMul(chain string, count int) string {
	if count < 0 {
		panic("string can not be multiplied a negative number of time")
	}
	result := ""
	for i:=0; i<count; i++ {
		result += chain
	}
	return result
}

func StrReplace(chain string, targetAndReplacements ...string) string {
	result := chain
	for i:=0; i<len(targetAndReplacements)-1; i++ {
		result = strings.Replace(chain, targetAndReplacements[i], targetAndReplacements[i+1], -1)
	}
	return result
}

func StrCrop(chain string, begining int, end int) string {
    runes := []rune(chain)
    return string(runes[begining:len(runes)-end])
}

func StrSplit(text string, sep rune) []string {
		result := []string{}
		chain := []rune{}
		runes := []rune(text)
		for i:=0; i < len(runes); i++ {
			if runes[i] == sep {
				result = append(result, string(chain))
				chain = []rune{}
			} else if runes[i] != '\r' {
				chain = append(chain, runes[i])
			}
		}
		if len(chain) > 0 {
			result = append(result, string(chain))
		} else {
			result = append(result, "")
		}
		return result
	}

func PanicParseFloat64(s string) float64 {
	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(fmt.Sprintf("Can't parse float64 from \"%s\"", s))
	}
	return result
}

func PanicParseInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("Can't parse float64 from \"%s\"", s))
	}
	return result
}

func listDir(directory string) []string {
    var fileNames []string
    _ = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            fileNames = append(fileNames, path)
        }
        return nil
    })
    return fileNames
}

func ReadFile(filepath string) string {
	output, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(fmt.Sprintf("Error reading file : %v", err))
	}
	return string(output)
}

func WriteFile(filepath string, content string) {
	err := ioutil.WriteFile(filepath, []byte(content), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func ReadCSV(filepath string) ([]string, []map[string]string) {
	output, err := ioutil.ReadFile(filepath)
	content := string(output)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	result := []map[string]string{}
	lines := StrSplit(content, '\n')
	lines = Filter(func(line string)bool{return len(line)>0}, lines)
	{
		n := -1
		for _, line := range lines {
			m := len(StrSplit(line, ','))
			if n != m {
				if n == - 1 {
					n = m
				} else {
					panic(fmt.Sprintf("%d != %d\n\"%s\"\n", n, m, line))
				}
			}
		} 
	}
	keys := StrSplit(lines[0], ','	)
	lines = lines[1:]
	for _, line := range lines {
		l := StrSplit(line, ',')
		obj := map[string]string{}
		for i, key := range keys {
			obj[key] = l[i]
		}
		result = append(result, obj)
	}
	l := -1
	for _, line := range result {
		if l == -1 {
			l = len(line)
		} else {
			if len(line) != l {
				panic("all lines do not have the same number of lines")
			}
		}
	}
	return keys, result
}