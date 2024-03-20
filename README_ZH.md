# A Golang runtime Python plugin library

# goExecPy
- [goExecPy](#goExecPy)
  - [Introduction](#introduction)
  - [Requirement](#requirement)
  - [Intall](#intall)
  - [Usage](#usage)
  - [Thanks to](#thanks-to)
## Introduction
参考`github.com/cokeBeer/execjs`, 实现了一个golang的执行python代码的插件库

## Requirement
环境变量path中存在python。

## Install
```
go get -u github.com/qmdemon/goExecPy
```
## Usage
可以使用Eval方法获取表达式的值，这将输出`Hello World`
```go
output, err := execjs.Eval(`"Hello " + "World"`)
if err != nil {
    log.Fatal(err)
}
fmt.Println(output)
```
可以使用Compile方法编译一个Context，然后调用。这将输出`3`,使用Compile方法时一定要确保Python代码块缩进正确
```go
c, _ := Compile(`
def add(x, y):
    return x + y
`)
output, err := c.Call("add", 1, 2)
if err != nil {
    log.Fatal(err)
}
fmt.Println(output)
```
更多用法参见测试文件`execpy_test.go`
> 注意：因为返回的是`interface{}`类型的变量，使用时要进行类型断言，例如
```go
output.(string) //这将输出值变为string类型
output.([]interface{}) //这将输出值变为slice类型
```
## Thanks to
[execjs](https://github.com/cokeBeer/execjs)
