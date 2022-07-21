package rwconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type RWConfiguration struct {
	Name                string
	Repository          string
	ContainerRepository string
	CanUseDocker        bool
	SudoDocker          bool
	Fuzzsettings        RWFuzzsettings
	Gosettings          RWgosettings
	Dependencysettings  RWdependencysettings
	Sourcesettings      RWsourceconfig
}

type RWFuzzsettings struct {
	FUZZ_TEST_RUN_TIME int64
	FUZZ_TEST_TIMEOUT  int64
	FUZZ_NUM_PROCESSES int64
}

type RWgosettings struct {
	GO111MODULE string
	GOPROXY     string
	GONOSUMDB   string
	GOOS        string
	GOARCH      string
	GO_BIN      string
}

type RWdependencysettings struct {
	GOFUZZ_COMMIT            string
	GOFUZZ_BUILD_PKG         string
	GOFUZZ_BUILD_PKG_VERSION string
	GOFUZZ_PKG               string
	GOFUZZ_PKG_VERSION       string
	GOFUZZ_DEP_PKG           string
	GOFUZZ_DEP_PKG_VERSION   string
	GOCYCLO_PKG              string
	GOCYCLO_PKG_VERSION      string
}

type RWsourceconfig struct {
	Commit       string
	Branch       string
	FuzzPackages []Fuzzpackage
}

type Fuzzpackage struct {
	FuncName     string
	PackageName  string
	BuildArchive string
}

func DefaultConfig() RWConfiguration {
	golangCfg := RWgosettings{
		GO111MODULE: "on",
		GOPROXY:     "https://gomodules.cbhq.net/",
		GONOSUMDB:   "github.cbhq.net",
		GOOS:        "linux",
		GOARCH:      "amd64",
		GO_BIN:      os.Getenv("GOPATH") + "/bin",
	}
	dependenyCfg := RWdependencysettings{
		GOFUZZ_COMMIT:            "latest",
		GOFUZZ_BUILD_PKG:         "github.com/dvyukov/go-fuzz/go-fuzz-build",
		GOFUZZ_BUILD_PKG_VERSION: "latest",
		GOFUZZ_PKG:               "github.com/dvyukov/go-fuzz/go-fuzz",
		GOFUZZ_PKG_VERSION:       "latest",
		GOFUZZ_DEP_PKG:           "github.com/dvyukov/go-fuzz/go-fuzz-dep",

		GOFUZZ_DEP_PKG_VERSION: "latest",
		GOCYCLO_PKG:            "github.com/fzipp/gocyclo/cmd/gocyclo",
		GOCYCLO_PKG_VERSION:    "latest",
	}

	fuzzCfg := RWFuzzsettings{
		FUZZ_TEST_RUN_TIME: 30,
		FUZZ_TEST_TIMEOUT:  20,
		FUZZ_NUM_PROCESSES: 4,
	}

	sourcecfg := RWsourceconfig{
		Commit:       "",
		Branch:       "main",
		FuzzPackages: make([]Fuzzpackage, 0),
	}

	config := RWConfiguration{
		Name:                "",
		Repository:          "",
		ContainerRepository: "docker.io",
		CanUseDocker:        false,
		SudoDocker:          false,
		Fuzzsettings:        fuzzCfg,
		Gosettings:          golangCfg,
		Dependencysettings:  dependenyCfg,
		Sourcesettings:      sourcecfg,
	}
	return config
}

// configRead configuration from json file
// configpath = "./config/config.json"
func ConfigRead(configpath string) (*RWConfiguration, error) {

	file, err := ioutil.ReadFile(configpath)
	if err != nil {
		fmt.Printf("Could not find config file in path: %s because %s", configpath, err)
		return nil, err
	}

	configuration := RWConfiguration{}
	err = json.Unmarshal(file, &configuration)
	if err != nil {
		fmt.Printf("Could not Unmarshal config file in path: %s because %s", configpath, err)
		return nil, err
	}
	return &configuration, nil
}

// configExists checks if the config file exists
// configpath = "./config/config.json"
func ConfigExists(configpath string) bool {
	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		fmt.Printf("Could not find config file in path: %s because %s", configpath, err)
		return false
	}
	return true
}

// Write configuration to json file
// config = Configuration object
// configpath = "./config/config.json"
func ConfigWrite(config RWConfiguration, configpath string) error {
	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("Could not convert the Configuration to []byte because %s", err)
		return err
	}
	err = ioutil.WriteFile(configpath, bytes, 0644)
	if err != nil {
		fmt.Printf("Could not write config file to path: %s because %s", configpath, err)
		return err
	}
	return nil
}
