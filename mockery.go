package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ryanbrainard/mockery/mockery"
)

const regexMetadataChars = "\\.+*?()|[]{}^$"

var fName = flag.String("name", "", "name or matching regular expression of interface to generate mock for")
var fPrint = flag.Bool("print", false, "print the generated mock to stdout")
var fOutput = flag.String("output", "./mocks", "directory to write mocks to")
var fDir = flag.String("dir", ".", "directory to search for interfaces")
var fRecursive = flag.Bool("recursive", false, "recurse search into sub-directories")
var fAll = flag.Bool("all", false, "generates mocks for all found interfaces in all sub-directories")
var fIP = flag.Bool("inpkg", false, "generate a mock that goes inside the original package")
var fCase = flag.String("case", "camel", "name the mocked file using casing convention")
var fNote = flag.String("note", "", "comment to insert into prologue of each generated file")

func main() {
	flag.Parse()

	var recursive bool
	var filter *regexp.Regexp
	var err error
	var limitOne bool

	if *fName != "" && *fAll {
		fmt.Fprintln(os.Stderr, "Specify -name or -all, but not both")
		os.Exit(1)
	} else if *fName != "" {
		recursive = *fRecursive
		if strings.ContainsAny(*fName, regexMetadataChars) {
			if filter, err = regexp.Compile(*fName); err != nil {
				fmt.Fprintln(os.Stderr, "Invalid regular expression provided to -name")
				os.Exit(1)
			}
		} else {
			filter = regexp.MustCompile(fmt.Sprintf("^%s$", *fName))
			limitOne = true
		}
	} else if *fAll {
		recursive = true
		filter = regexp.MustCompile(".*")
	} else {
		fmt.Fprintln(os.Stderr, "Use -name to specify the name of the interface or -all for all interfaces found")
		os.Exit(1)
	}

	generated := walkDir(*fDir, recursive, filter, limitOne)

	if *fName != "" && !generated {
		fmt.Printf("Unable to find %s in any go files under this path\n", *fName)
		os.Exit(1)
	}
}

func walkDir(dir string, recursive bool, filter *regexp.Regexp, limitOne bool) (generated bool) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		path := filepath.Join(dir, file.Name())

		if file.IsDir() {
			if recursive {
				generated = walkDir(path, recursive, filter, limitOne) || generated
				if generated && limitOne {
					return
				}
			}
			continue
		}

		if !strings.HasSuffix(path, ".go") {
			continue
		}

		p := mockery.NewParser()

		err = p.Parse(path)
		if err != nil {
			continue
		}

		for _, iface := range p.Interfaces() {
			if !filter.MatchString(iface.Name) {
				continue
			}
			genMock(iface)
			generated = true
			if limitOne {
				return
			}
		}
	}

	return
}

func genMock(iface *mockery.Interface) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Unable to generated mock for '%s': %s\n", iface.Name, r)
			return
		}
	}()

	var out io.Writer

	pkg := "mocks"
	name := iface.Name
	caseName := iface.Name
	if *fCase == "underscore" {
		rxp := regexp.MustCompile("(.)([A-Z])")
		caseName = strings.ToLower(rxp.ReplaceAllString(caseName, "$1_$2"))
	}

	if *fPrint {
		out = os.Stdout
	} else {
		var path string

		if *fIP {
			path = filepath.Join(filepath.Dir(iface.Path), "mock_"+caseName+".go")
		} else {
			path = filepath.Join(*fOutput, caseName+".go")
			os.MkdirAll(filepath.Dir(path), 0755)
			pkg = filepath.Base(filepath.Dir(path))
		}

		f, err := os.Create(path)
		if err != nil {
			fmt.Printf("Unable to create output file for generated mock: %s\n", err)
			os.Exit(1)
		}

		defer f.Close()

		out = f

		fmt.Printf("Generating mock for: %s\n", name)
	}

	gen := mockery.NewGenerator(iface)

	gen.GeneratePrologueNote(*fNote)

	if *fIP {
		gen.GenerateIPPrologue()
	} else {
		gen.GeneratePrologue(pkg)
	}

	err := gen.Generate()
	if err != nil {
		fmt.Printf("Error with %s: %s\n", name, err)
		os.Exit(1)
	}

	err = gen.Write(out)
	if err != nil {
		fmt.Printf("Error writing %s: %s\n", name, err)
		os.Exit(1)
	}
}
