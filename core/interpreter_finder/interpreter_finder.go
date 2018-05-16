package interpreterFinder

import (
	"../configurator"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type InterpreterFinder struct {
	CurrentDir string
}

func (f *InterpreterFinder) HaveBuiltIn() bool {
	_, e := os.Stat(filepath.Join(f.CurrentDir, builtinRelativeFilePath))
	exists := !os.IsNotExist(e)

	if exists && e == nil {
		return true
	}

	return false
}

func (f *InterpreterFinder) FindBuiltin() (path string) {
	if f.HaveBuiltIn() {
		path = filepath.Join(f.CurrentDir, builtinRelativeFilePath)
	}
	return
}

func (f *InterpreterFinder) Find() *string {
	// Built-in interpreter
	//if f.Config.UseBuiltinInterpreter {
	//	builtInPath := builtinRelativeFilePath
	//	_, e := os.Stat(builtInPath)
	//	exists := !os.IsNotExist(e)
	//
	//	if exists && e == nil {
	//		return &builtInPath
	//	}
	//}

	// External interpreter
	for _, path := range exactFilePaths() {
		_, e := os.Stat(path)
		exists := !os.IsNotExist(e)

		if exists && e == nil {
			return &path
		}
	}

	return nil
}

func (f *InterpreterFinder) Check(command string) (version string, e error) {
	out, e := exec.Command(configurator.ExpandInterpreterCommand(command), "-version").Output()
	if e != nil {
		return "", e
	}

	replacer := strings.NewReplacer("\n", "", "\r", "")
	version = replacer.Replace(string(out))

	return
}
