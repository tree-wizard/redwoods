package analyzer

type Function struct {
	Name                 string            `json:"name"`
	FileName             string            `json:"filename"`
	Documentation        string            `json:"documentation"`
	Parameter            map[string]string `json:"parameter"`
	LinesOfCode          int               `json:"lines_of_code"`
	HasDocumentation     bool              `json:"has_documentation"`
	IsExported           bool              `json:"is_exported"`
	CyclomaticComplexity int               `json:"cyclomatic_complexity"`
	HasFuzzTests         bool              `json:"has_fuzz_tests"`
	HasGoFuzzTests       bool              `json:"has_gofuzz_tests"`
}

type Package struct {
	Name           string              `json:"name"`
	Location       string              `json:"location"`
	LinesOfCode    int                 `json:"lines_of_code"`
	FuzzFileCount  int                 `json:"fuzz_file_count"`
	TestFileCount  int                 `json:"test_file_count"`
	TotalFileCount int                 `json:"file_count"`
	HasTests       bool                `json:"has_tests"`
	HasFuzzTests   bool                `json:"has_fuzz_tests"`
	HasGoFuzzTests bool                `json:"has_gofuzz_tests"`
	Imports        []string            `json:"imported_packages"`
	Files          []string            `json:"files"`
	Functions      map[string]Function `json:"functions"`
}

type Project struct {
	Name             string    `json:"name"`
	Repo             string    `json:"git_repository"`
	LinesOfCode      int       `json:"lines_of_code"`
	FuzzFileCount    int       `json:"fuzz_file_count"`
	TestFileCount    int       `json:"test_file_count"`
	TotalFileCount   int       `json:"file_count"`
	PackageCount     int       `json:"package_count"`
	FuzzPackageCount int       `json:"fuzz_package_count"`
	Packages         []Package `json:"packages"`
	FuzzPackages     []Package `json:"fuzz_packages"`
}

type CyclomaticComplexity struct {
	Package    string
	Function   string
	Complecity int
}
