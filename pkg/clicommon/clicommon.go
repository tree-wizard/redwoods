package clicommon

import (
	"strconv"

	"github.com/manifoldco/promptui"
)

//PromptString requires a console input from string
func PromptString(name string, defaultval string) (string, error) {
	prompt := promptui.Prompt{
		Label:   name,
		Default: defaultval,
	}

	return prompt.Run()
}

//PromptNotEmptyString requires a console input from string that can not be empty
func PromptNotEmptyString(name string, defaultval string) (string, error) {
	prompt := promptui.Prompt{
		Label:    name,
		Validate: ValidateEmptyInput,
		Default:  defaultval,
	}

	return prompt.Run()
}

//PromptInteger requires a console input from string that can not be empty
func PromptInteger(name string) (int64, error) {
	prompt := promptui.Prompt{
		Label:    name,
		Validate: ValidateIntegerNumberInput,
	}

	promptResult, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	parseInt, _ := strconv.ParseInt(promptResult, 0, 64)
	return parseInt, nil
}

//PromptConfirm requires a console conform input from string that can not be empty
func PromptConfirm(name string) (string, error) {
	prompt := promptui.Prompt{
		Label:     name,
		IsConfirm: true,
	}

	return prompt.Run()
}

//PromptGOOS is a select function for GOOS
func PromptGOOS(name string) (string, error) {
	prompt := promptui.Select{
		Label: name,
		Items: []string{"aix", "android", "darwin", "dragonfly", "freebsd", "hurd", "illumos", "ios", "js", "linux", "nacl", "netbsd", "openbsd", "plan9", "solaris", "windows", "zos"},
	}
	_, result, err := prompt.Run()

	return result, err
}

//PromptGOARCH is a select function for GOARCH
func PromptGOARCH(name string) (string, error) {
	prompt := promptui.Select{
		Label: name,
		Items: []string{"386", "amd64", "amd64p32", "arm", "arm64", "arm64be", "armbe", "loong64", "mips", "mips64", "mips64le", "mips64p32", "mips64p32le", "mipsle", "ppc", "ppc64", "ppc64le", "riscv", "riscv64", "s390", "s390x", "sparc", "sparc64", "wasm"},
	}
	_, result, err := prompt.Run()

	return result, err
}
