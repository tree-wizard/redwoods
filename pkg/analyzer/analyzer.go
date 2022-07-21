package analyzer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func RunGoCyclo(path string) []CyclomaticComplexity {
	var stdBuffer bytes.Buffer
	e := exec.Command("gocyclo", path)
	e.Stdout = &stdBuffer
	e.Stderr = &stdBuffer
	err := e.Run()
	if err != nil {
		panic(err)
	}
	return parseGocyclo(stdBuffer.Bytes())
}

//takes the gocyclo output and converts it into an object
func parseGocyclo(input []byte) []CyclomaticComplexity {
	complexcities := make([]CyclomaticComplexity, 0)
	for _, line := range strings.Split(string(input), "\n") {
		x := CyclomaticComplexity{
			Package:    "",
			Function:   "",
			Complecity: 0,
		}
		for index, v := range strings.Split(line, " ") {
			switch index {
			case 0: //complexity
				x.Complecity, _ = strconv.Atoi(v)
			case 1: //package
				x.Package = v
			case 2: //function
				x.Function = v
				complexcities = append(complexcities, x)
				x = CyclomaticComplexity{
					Package:    "",
					Function:   "",
					Complecity: 0,
				}
			}

		}
	}
	return complexcities
}

//NewAnalysis Will find all subdirectories of the project in order to analyze all go files
func NewAnalysis(root string) Project {
	set := token.NewFileSet()
	project := Project{
		Name:             "",
		Repo:             root,
		Packages:         make([]Package, 0),
		FuzzPackages:     make([]Package, 0),
		LinesOfCode:      0,
		FuzzFileCount:    0,
		TestFileCount:    0,
		TotalFileCount:   0,
		PackageCount:     0,
		FuzzPackageCount: 0,
	}
	subDirectories, err := getSubDirectories(root)
	if err != nil {
		panic(err)
	}

	subDirectories = append(subDirectories, root)

	for _, sub := range subDirectories {
		project.Packages = append(project.Packages, ScanAst(set, sub, filterFiles)...)
	}
	return project
}

//filterFiles makes sure we only filter .go files
func filterFiles(fi os.FileInfo) bool {
	return path.Ext(fi.Name()) == ".go"
}

//Analyze will take the results of a run and merge them to the desired format
func Analyze(project Project, onlyfuzz bool) Project {

	for _, pkg := range project.Packages {
		project.PackageCount++
		project.LinesOfCode += pkg.LinesOfCode
		project.TestFileCount += pkg.TestFileCount
		project.TotalFileCount += pkg.TotalFileCount
		if pkg.HasFuzzTests || pkg.HasGoFuzzTests {
			project.FuzzPackages = append(project.FuzzPackages, pkg)
			project.FuzzFileCount += pkg.FuzzFileCount
			project.FuzzPackageCount++
		}

	}
	if onlyfuzz {
		project.Packages = nil
	}
	return project
}

//ToFuzzable returns the package path of all fuzzable packages in the run
func ToGofuzz(project Project, path string) string {
	results := ""
	for _, pkg := range project.FuzzPackages {

		if pkg.HasGoFuzzTests {
			results += strings.Replace(pkg.Location, path, "", -1) + " "
		}
	}
	return results

}

//ToConsole returns a project output to console. abandoned because of tmi
func ToConsole(project Project) {
	fmt.Println("########")
	fmt.Println("Project: ", project.Repo)
	fmt.Println("########")
	fmt.Println("Go-Files: ", project.TotalFileCount)
	fmt.Println("Testfiles ", project.TestFileCount)
	fmt.Println("Fuzzfiles: ", project.FuzzFileCount)
	fmt.Println("Fuzzable Packages: ", project.FuzzPackageCount)
	pkgs := project.Packages
	if len(project.Packages) == 0 {
		pkgs = project.FuzzPackages
	}
	for _, pkg := range pkgs {
		fmt.Println("########")
		fmt.Println("Package: ", pkg.Name)
		fmt.Println("LOC: ", pkg.LinesOfCode)
		fmt.Println("Location: ", pkg.Location)
		if pkg.HasFuzzTests {
			fmt.Println("Found std-Fuzztests!")
		}
		if pkg.HasGoFuzzTests {
			fmt.Println("Found go-fuzztests!")
		}

		fmt.Println("Found ", len(pkg.Functions), " Functions (including test ", pkg.TestFileCount, "and", pkg.FuzzFileCount, "fuzz)")
		fmt.Println("########")
		for _, f := range pkg.Functions {
			fmt.Println("---")
			fmt.Println("Name:", f.Name, "LOC:", f.LinesOfCode, "Cyclomatic Complexity:", f.CyclomaticComplexity)
			fmt.Println("Parameter:")
			for k, v := range f.Parameter {
				fmt.Println(k, " ", v)
			}

		}
	}

}

//FromJson marshals a project to file
func ToJson(project Project, outputPath string) {
	json, err := json.MarshalIndent(project, "", " ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(outputPath, json, 0644)
	if err != nil {
		panic(err)
	}
}

//FromJson unmarshals a file to project
func FromJson(path string) (Project, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Could not Unmarshal config analyis", err)
		return Project{}, err
	}
	project := Project{}
	err = json.Unmarshal(file, &project)
	if err != nil {
		fmt.Printf("Could not Unmarshal config analyis", err)
		return Project{}, err
	}
	return project, nil
}

//getSubDirectories lists all the subdirectories of a root
func getSubDirectories(root string) ([]string, error) {
	dir, err := os.Open(root)
	if err != nil {
		return nil, err
	}

	infos, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}
	subDirectories := make([]string, 0)

	for _, fi := range infos {
		//we assume all "." dirs to be configs so we ignore them
		if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
			subDirectories = append(subDirectories, path.Join(root, fi.Name()))
			children, _ := getSubDirectories(path.Join(root, fi.Name()))
			if len(children) > 0 {
				subDirectories = append(subDirectories, children...)
			}
		}
	}
	return subDirectories, nil
}

//ScanAst creates an AST assasment based on an input folder which usually is the package
func ScanAst(set *token.FileSet, folder string, filter func(fi fs.FileInfo) bool) []Package {
	pkgs := make([]Package, 0)
	packs, err := parser.ParseDir(set, folder, filter, parser.ParseComments)
	if err != nil {
		fmt.Println("Failed to parse package:", err)
		os.Exit(1)
	}

	//fmt.Println("found: ", len(packs), " Packages")
	for _, pack := range packs {
		pkg := Package{
			Name:           pack.Name,
			Location:       folder,
			Files:          []string{},
			Functions:      make(map[string]Function),
			LinesOfCode:    0,
			FuzzFileCount:  0,
			TestFileCount:  0,
			TotalFileCount: 0,
			HasTests:       false,
			HasFuzzTests:   false,
			HasGoFuzzTests: false,
			Imports:        make([]string, 0),
		}
		//fmt.Println("Package Name :", pack.Name)
		//fmt.Println("##")
		for fName, f := range pack.Files {
			complexcities := make([]CyclomaticComplexity, 0)
			if strings.HasSuffix(fName, ".go") {
				complexcities = RunGoCyclo(fName)
			}
			pkg.Files = strAppendIfMissing(pkg.Files, fName)
			inFuzzFile := false
			inGoFuzzFile := false
			if strings.HasSuffix(fName, "_test.go") {
				pkg.TestFileCount++
				pkg.HasTests = true
			} else if strings.HasSuffix(fName, "fuzz.go") {
				gofuzzTest := false
				for _, c := range f.Comments {
					for _, c1 := range c.List {
						if c1.Text == "//+build gofuzz" {
							gofuzzTest = true
						}
					}
				}
				pkg.FuzzFileCount++
				if gofuzzTest {
					pkg.HasGoFuzzTests = true
					inGoFuzzFile = true
				} else {
					pkg.HasFuzzTests = true
					inFuzzFile = true
				}

			} else {
				pkg.TotalFileCount++
			}

			//For All decals in the AST
			for _, d := range f.Decls {
				//if it is a function
				if fn, isFn := d.(*ast.FuncDecl); isFn {
					theFunc := Function{
						Name:                 fn.Name.Name,
						FileName:             fName,
						Documentation:        "",
						Parameter:            make(map[string]string),
						LinesOfCode:          0,
						HasDocumentation:     false,
						IsExported:           false,
						CyclomaticComplexity: 0,
						HasFuzzTests:         inFuzzFile,
						HasGoFuzzTests:       inGoFuzzFile,
					}
					if f.Doc != nil {
						theFunc.HasDocumentation = true
					}
					if fn.Type != nil {
						if fn.Type.Params != nil {
							for _, param := range fn.Type.Params.List {
								for _, name := range param.Names {
									theFunc.Parameter[name.Name] = FormatNode(param.Type)
								}
							}
						}
					}
					for _, c := range complexcities {
						if c.Function == theFunc.Name {
							theFunc.CyclomaticComplexity = c.Complecity
						}
					}

					theFunc.IsExported = fn.Name.IsExported()

					if fn.Doc != nil {
						theFunc.Documentation = fn.Doc.Text()
					}
					bodylength, err := getFunctionBodyLength(fn, set)
					if err != nil {
						fmt.Println(err)
					}
					theFunc.LinesOfCode = bodylength
					pkg.LinesOfCode += bodylength
					pkg.Functions[theFunc.Name] = theFunc
				}
			}
			for _, imported := range f.Imports {
				trimmedImport := strings.Trim(imported.Path.Value, "\"")
				pkg.Imports = strAppendIfMissing(pkg.Imports, trimmedImport)
			}

		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}

//strAppendIfMissing is a helper function to keep our slices unique
func strAppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

//getFunctionBodyLength returns the body length of the function
func getFunctionBodyLength(f *ast.FuncDecl, fs *token.FileSet) (int, error) {
	if fs == nil {
		return 0, errors.New("FileSet is nil")
	}
	if f.Body == nil {
		return 0, nil
	}
	if !f.Body.Lbrace.IsValid() || !f.Body.Rbrace.IsValid() {
		return 0, fmt.Errorf("function %s is not syntactically valid", f.Name.String())
	}
	length := fs.Position(f.Body.Rbrace).Line - fs.Position(f.Body.Lbrace).Line - 1
	if length > 0 {
		return length, nil
	}
	return 0, nil
}

//formats a nodeput to string
func FormatNode(node ast.Node) string {
	buf := new(bytes.Buffer)
	_ = format.Node(buf, token.NewFileSet(), node)
	return buf.String()
}
