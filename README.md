# A Golang runtime Python plugin library

## Introduction
This library is a Golang plugin library that allows executing Python code within a Golang runtime. It is based on the [execjs](https://github.com/cokeBeer/execjs) library, which provides a way to execute JavaScript code within a Golang runtime.

[中文](https://github.com/qmdemon/goExecPy/blob/master/README_ZH.md)

## Requirement
The environment variable `PATH` must contain the path to the Python executable.

## Install      
You can install the `goExecPy` library using the following command:
 ```
 go get -u github.com/qmdemon/goExecPy
 ```
## Usage
You can use the `Eval` method to get the value of an expression, and the `Compile` method to compile a Python code block and call it. Here's an example:
 ```go
 output, err := execjs.Eval(`"Hello " + "World"`)
 if err != nil {
     log.Fatal(err)
 }
 fmt.Println(output)
 ```

When using the `Compile` method, it is important to ensure that the Python code block is indented correctly.

And here's an example of using the `Compile` method:
 ```go
 c, _ := execjs.Compile(`
 def add(x, y):
     return x + y
 `)
 output, err := c.Call("add", 1, 2)
 if err != nil {
     log.Fatal(err)
 }
 fmt.Println(output)
 ```
More usage examples can be found in the `execpy_test.go` file.

Please note that the returned value is of type `interface{}`, so you may need to perform type assertions when using the returned value. For example:
 ```go
 output.(string) //This will convert the value to string type
 output.([]interface{}) //This will convert the value to slice type
 ```
## Thanks to
* [execjs](https://github.com/cokeBeer/execjs) for providing the basic implementation of executing JavaScript code within a Golang runtime.