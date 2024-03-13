package goExecPy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type ExternalRuntime struct {
	name          string
	command       []string
	runner_source string
	encoding      string
	tempfile      bool
	available     bool
	binary_cache  []string
}

func BuildExternalRuntime(name string, command []string, runner_source string) *ExternalRuntime {
	r := &ExternalRuntime{name: name, command: command, runner_source: runner_source}
	r.available = (r.binary() != nil)
	return r
}

func (r *ExternalRuntime) Eval(source string) (interface{}, error) {
	return r.Compile("").Eval(source)
}

func (r *ExternalRuntime) Compile(source string) RuntimeContextInterface {
	if r.Is_available() {
		r.runner_source = strings.Replace(r.runner_source, "#{source}", source, 1)
		return r.compile(source)
	} else {
		return nil
	}
}

func (r *ExternalRuntime) Is_available() bool {
	return r.available
}

func (r *ExternalRuntime) compile(source string) RuntimeContextInterface {
	return &Context{runtime: r, source: source, tempfile: r.tempfile}
}

func (r *ExternalRuntime) binary() []string {
	if r.binary_cache == nil {
		r.binary_cache = which(r.command)
	}
	return r.binary_cache
}

type Context struct {
	runtime  *ExternalRuntime
	source   string
	cwd      string
	tempfile bool
}

func (c *Context) Exec_(source string) (interface{}, error) {
	if !c.Is_available() {
		return "", RuntimeUnavailableError{Message: fmt.Sprintf("runtime is not available on this system")}
	}
	output, err := c.exec_(source)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (c *Context) Eval(source string) (interface{}, error) {
	if !c.Is_available() {
		return "", RuntimeUnavailableError{Message: fmt.Sprintf("runtime is not available on this system")}
	}
	source_bytes, _ := json.Marshal(source)
	output, err := c.Exec_(string(source_bytes))
	if err != nil {
		return "", err
	}
	return output, nil
}

func (c *Context) Call(name string, args ...interface{}) (interface{}, error) {
	if !c.Is_available() {
		return "", RuntimeUnavailableError{Message: fmt.Sprintf("runtime is not available on this system")}
	}
	output, err := c.call(name, args...)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (c *Context) Is_available() bool {
	return c.runtime.Is_available()
}

//func (c *Context) eval(source string) (interface{}, error) {
//	var data string
//	if len(strings.TrimSpace(source)) == 0 {
//		data = "''"
//	} else {
//		data = "'('+'" + source + "'+')'"
//	}
//	code := fmt.Sprintf("return eval(%s)", data)
//	return c.Exec_(code)
//}

func (c *Context) exec_(source string) (interface{}, error) {
	//if c.source != "" {
	//	source = c.source + "\n" + source
	//}
	var (
		output string
		err    error
	)
	if c.tempfile {
		output, err = c.exec_with_tempfile(source)
		if err != nil {
			return "", err
		}
	} else {
		output, err = c.exec_with_pipe(source)
		if err != nil {
			return "", RuntimeUnavailableError{Message: fmt.Sprintf("%v \nPython Error:\n %s", err, output)}
		}
	}
	return c.extract_result(output)
}

func (c *Context) call(fun string, args ...interface{}) (interface{}, error) {
	stringArgs := make([]string, 0)
	var fun_args string
	fun2, _ := json.Marshal(fun)
	stringArgs = append(stringArgs, string(fun2))
	arg, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	if string(arg) != "null" {
		stringArgs = append(stringArgs, string(arg))
		fun_args = strings.Join(stringArgs, ",")
	} else {
		fun_args = string(fun2)
	}

	output, err := c.Exec_(fun_args)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (c *Context) exec_with_tempfile(_ string) (string, error) {
	return "", nil
}

func (c *Context) exec_with_pipe(encoded_source string) (string, error) {
	command := c.runtime.command
	cmd := exec.Command(command[0], command[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	input := c.compile(encoded_source)
	//fmt.Println(input)
	_, err = stdin.Write([]byte(input))
	if err != nil {
		return "", err
	}
	stdin.Close()
	err = cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		return errStr, err
	}
	return outStr, nil
}

func (c *Context) extract_result(output string) (interface{}, error) {
	ret := strings.Split(output, "\n")
	data := ret[len(ret)-2]
	var res []interface{}
	json.Unmarshal([]byte(data), &res)
	if len(res) == 1 {
		if res[0].(string) == "ok" {
			return "", nil
		} else {
			return "", RuntimeUnavailableError{Message: fmt.Sprintf("%v", res[1])}
		}
	} else {
		return res[1], nil
	}
}

func (c *Context) compile(encoded_source string) string {
	runner_source := c.runtime.runner_source
	return strings.Replace(runner_source, "#{encoded_source}", encoded_source, 1)
}

func which(command []string) []string {
	name := command[0]
	path := CheckCommandExists(name)
	binary := make([]string, len(command))
	if path {
		copy(binary, command)
		//binary[0] = path
		return binary
	} else {
		return nil
	}
}

func CheckCommandExists(cmd string) bool {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return false
	}
	return len(path) > 0
}

//func find_executable(prog string) string {
//	pathlist := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
//	filename := ""
//	if runtime.GOOS == "windows" {
//		prog += ".exe"
//	}
//	for _, dir := range pathlist {
//		filename = path.Join(dir, prog)
//		filename = strings.ReplaceAll(filename, "\\", "/")
//		filename = strings.ReplaceAll(filename, "//", "/")
//		//fmt.Println(filename)
//		info, err := os.Stat(filename)
//		if err != nil {
//			fmt.Println(filename, err)
//			continue
//		}
//		if info.Mode()&0111 == 0111 {
//			return filename
//		}
//	}
//	return ""
//
//}

func python() *ExternalRuntime {
	r := python_python_3()
	if r.Is_available() {
		return r
	}
	return python_python()
}

func python_python_3() *ExternalRuntime {
	return BuildExternalRuntime("python", []string{"python"}, Python_source)
}

func python_python() *ExternalRuntime {
	return BuildExternalRuntime("python3", []string{"python3"}, Python_source)
}
