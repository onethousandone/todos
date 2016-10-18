package parser

import (
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"regexp"

	"golang.org/x/tools/go/loader"
)

type Todo struct {
	Text string
	Pos  token.Position
}

func ParsePackage(dir string) (*loader.Program, error) {
	relDir, err := filepath.Rel(filepath.Join(build.Default.GOPATH, "src"), dir)
	if err != nil {
		return nil, fmt.Errorf("provided directory not under GOPATH (%s): %v",
			build.Default.GOPATH, err)
	}

	conf := loader.Config{
		TypeChecker: types.Config{
			FakeImportC:      true,
			IgnoreFuncBodies: true,
		},
		ParserMode: parser.ParseComments,
	}
	conf.Import(relDir)
	program, err := conf.Load()
	if err != nil {
		return nil, fmt.Errorf("couldn't load package: %v", err)
	}

	return program, nil
}

func GetTodos(prg *loader.Program) (todos []Todo) {
	re := regexp.MustCompile(`\r?\n`)
	for _, pkg := range prg.InitialPackages() {
		for _, file := range pkg.Files {
			for _, cgroup := range file.Comments {
				for _, todo := range parseTodos(cgroup.Text()) {
					todos = append(todos, Todo{
						Text: re.ReplaceAllString(todo, ""),
						Pos:  prg.Fset.Position(cgroup.Pos()),
					})
				}
			}
		}
	}
	return
}

func parseTodos(str string) []string {
	re := regexp.MustCompile(`(?:TODO|FIXME):((.|\n)*)`)
	return re.FindAllString(str, -1)
}
