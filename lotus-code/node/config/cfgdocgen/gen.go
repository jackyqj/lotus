package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func run() error {
	tfb, err := ioutil.ReadFile("./node/config/types.go")
	if err != nil {
		return err
	}

	// could use the ast lib, but this is simpler

	type st int
	const (
		stGlobal st = iota // looking for typedef
		stType   st = iota // in typedef
	)

	lines := strings.Split(string(tfb), "\n")
	state := stGlobal

	type field struct {
		Name    string
		Type    string
		Comment string
	}

	var currentType string
	var currentComment []string

	out := map[string][]field{}

	for l := range lines {
		line := strings.TrimSpace(lines[l])

		switch state {
		case stGlobal:
			if strings.HasPrefix(line, "type ") {
				currentType = line
				currentType = strings.TrimPrefix(currentType, "type")
				currentType = strings.TrimSuffix(currentType, "{")
				currentType = strings.TrimSpace(currentType)
				currentType = strings.TrimSuffix(currentType, "struct")
				currentType = strings.TrimSpace(currentType)
				currentComment = nil
				state = stType
				continue
			}
		case stType:
			if strings.HasPrefix(line, "// ") {
				cline := strings.TrimSpace(strings.TrimPrefix(line, "//"))
				currentComment = append(currentComment, cline)
				continue
			}

			comment := currentComment
			currentComment = nil

			if strings.HasPrefix(line, "}") {
				state = stGlobal
				continue
			}

			f := strings.Fields(line)
			if len(f) < 2 { // empty or embedded struct
				continue
			}

			name := f[0]
			typ := f[1]

			out[currentType] = append(out[currentType], field{
				Name:    name,
				Type:    typ,
				Comment: strings.Join(comment, "\n"),
			})
		}
	}

	var outt []string
	for t := range out {
		outt = append(outt, t)
	}
	sort.Strings(outt)

	fmt.Print(`// Code generated by github.com/filecoin-project/lotus/node/config/cfgdocgen. DO NOT EDIT.

package config

type DocField struct {
	Name    string
	Type    string
	Comment string
}

var Doc = map[string][]DocField{
`)

	for _, typeName := range outt {
		typ := out[typeName]

		fmt.Printf("\t\"%s\": []DocField{\n", typeName)

		for _, f := range typ {
			fmt.Println("\t\t{")
			fmt.Printf("\t\t\tName: \"%s\",\n", f.Name)
			fmt.Printf("\t\t\tType: \"%s\",\n\n", f.Type)
			fmt.Printf("\t\t\tComment: `%s`,\n", f.Comment)
			fmt.Println("\t\t},")
		}

		fmt.Printf("\t},\n")
	}

	fmt.Println(`}`)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
