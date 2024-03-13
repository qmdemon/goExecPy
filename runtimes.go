package goExecPy

const (
	Python string = "python"
)

type Tuple struct {
	Name    string
	Runtime RuntimeInterface
}

func Init() *Tuple {
	return Register(Python, python())
}

func Register(name string, runtime RuntimeInterface) *Tuple {
	return &Tuple{Name: name, Runtime: runtime}
}

func GetRuntime(myruntime *Tuple) (RuntimeInterface, error) {
	runtime, err := find_available_runtime(myruntime)
	if err != nil {
		return nil, err
	} else {
		return runtime, nil
	}

}

func find_available_runtime(myruntime *Tuple) (RuntimeInterface, error) {
	var runtime RuntimeInterface
	runtime = nil
	if myruntime.Runtime.Is_available() {
		runtime = myruntime.Runtime
	}
	if runtime != nil {
		return runtime, nil
	} else {
		return nil, RuntimeUnavailableError{Message: "Could not find an available Python runtime."}
	}
}

func Eval(source string) (interface{}, error) {
	myruntime := Init()
	runtime, err := GetRuntime(myruntime)
	if err != nil {
		return "", err
	}
	return runtime.Eval(source)
}

//func Exec_(source string) (interface{}, error) {
//	myruntime := Init()
//	runtime, err := GetRuntime(myruntime)
//	if err != nil {
//		return "", err
//	}
//	return runtime.Exec_(source)
//}

func Compile(source string) (RuntimeContextInterface, error) {
	myruntime := Init()
	runtime, err := GetRuntime(myruntime)
	if err != nil {
		return nil, err
	}
	return runtime.Compile(source), nil
}
