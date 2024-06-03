
package main

import (
)

var scripts map[string]func() = map[string]func(){}

func script(name string, code func()) func() {
    scripts[name] = code
    return code
}

var postLogContent []string = []string{}
func postLog(args ...any) {
    postLogContent = append(postLogContent, Sprintf(args...))
}

func main() {
    code, found := scripts[GetOptions("script", "")]
    if found {
        code()
    } else {
       
    }
    if GetSwitch("post-log") {
        for _, msg := range postLogContent {
            Printfc[180]("%s\n", msg)
        }
    }
}