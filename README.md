# outfmt

`outfmt` 是一个面向 CLI 场景的 Go 输出库，用来把结构化数据渲染为 `table`、`yaml`、`json`，方便人类和 agent 阅读。

## 特性

- 统一入口渲染多种输出格式
- 默认支持 `struct`、`[]struct`、`map[string]any`、`[]map[string]any`
- `table` 支持嵌套 `struct` / `map` 扁平化输出，列名使用 `Parent.Child`
- `json` 默认格式化输出
- `table` 针对终端可读性做简单对齐
- 支持通过 `outfmt` tag 控制 table 列名或忽略字段

## 安装

```bash
go get github.com/kiry163/outfmt
```

发布后按版本安装：

```bash
go get github.com/kiry163/outfmt@v0.1.0
```

## 适用场景

- CLI 工具统一支持 `--output table|yaml|json`
- 同一份结构化数据同时服务人类阅读和 agent 消费
- 避免每个命令重复写一套输出格式转换逻辑

## 快速开始

```go
package main

import (
	"fmt"
	"os"

	"github.com/kiry163/outfmt"
)

type User struct {
	ID    int    `json:"id" yaml:"id" outfmt:"ID"`
	Name  string `json:"name" yaml:"name" outfmt:"Name"`
	Email string `json:"email" yaml:"email" outfmt:"Email"`
}

func main() {
	users := []User{
		{ID: 1, Name: "alice", Email: "alice@example.com"},
		{ID: 2, Name: "bob", Email: "bob@example.com"},
	}

	if err := outfmt.Render(os.Stdout, users, outfmt.Table); err != nil {
		panic(err)
	}

	fmt.Println()

	raw, err := outfmt.Marshal(users, outfmt.JSON)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(raw))
}
```

## 输出格式

### `table`

- 面向终端阅读
- 适合 `struct`、`[]struct`、`map[string]any`、`[]map[string]any`
- 会把嵌套 `struct` / `map` 扁平化成 `Parent.Child`

### `yaml`

- 更适合人类查看复杂嵌套结构
- 保留原始层级结构

### `json`

- 更适合程序、agent 或管道处理
- 默认使用缩进格式输出

## CLI 集成示例

```go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kiry163/outfmt"
)

type User struct {
	ID     int    `json:"id" yaml:"id" outfmt:"ID"`
	Name   string `json:"name" yaml:"name" outfmt:"Name"`
	Status string `json:"status" yaml:"status" outfmt:"Status"`
}

func main() {
	var output string

	flag.StringVar(&output, "output", string(outfmt.Table), "output format: table|yaml|json")
	flag.Parse()

	users := []User{
		{ID: 1, Name: "alice", Status: "active"},
		{ID: 2, Name: "bob", Status: "inactive"},
	}

	if err := outfmt.Render(os.Stdout, users, outfmt.Format(output)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

## 示例

直接看三种输出效果：

```bash
go run ./examples/basic
```

模拟 CLI 的 `--output` 切换：

```bash
go run ./examples/cli --output table
go run ./examples/cli --output yaml
go run ./examples/cli --output json
```

查看嵌套结构在 `table` / `yaml` / `json` 下的效果：

```bash
go run ./examples/nested
```

## `outfmt` tag

- `outfmt:"-"` 忽略字段
- `outfmt:"Display Name"` 设置 table 列名

## Options

### `WithJSONIndent`

自定义 JSON 缩进：

```go
raw, err := outfmt.Marshal(data, outfmt.JSON, outfmt.WithJSONIndent("    "))
```

### `WithEmptyValue`

自定义 table 中空单元格的占位值：

```go
err := outfmt.Render(os.Stdout, data, outfmt.Table, outfmt.WithEmptyValue("(empty)"))
```

## 嵌套结构说明

对于下面这种结构：

```go
type Profile struct {
	City   string `outfmt:"City"`
	Active bool   `outfmt:"Active"`
}

type User struct {
	ID      int            `outfmt:"ID"`
	Profile *Profile       `outfmt:"Profile"`
	Meta    map[string]any `outfmt:"Meta"`
	Tags    []string       `outfmt:"Tags"`
}
```

`table` 输出会类似：

```text
ID  Profile.City  Profile.Active  Meta.region  Meta.zone  Tags
--  ------------  --------------  -----------  ---------  ---------
1   shanghai      true            cn           east       [dev ops]
2   -             -               us           -          -
```

其中：

- 嵌套 `struct` 会展开成多列
- 嵌套 `map[string]any` 会展开成多列
- `slice` 当前仍作为单个单元格显示
- `nil` 指针字段会显示为空值占位符，默认是 `-`

## API

- `Render(w io.Writer, data any, format Format, opts ...Option) error`
- `Marshal(data any, format Format, opts ...Option) ([]byte, error)`
- `Format`: `outfmt.Table`、`outfmt.YAML`、`outfmt.JSON`

## 发布相关

- 变更记录见 [CHANGELOG.md](./CHANGELOG.md)
- 发布步骤见 [RELEASING.md](./RELEASING.md)
- 许可证见 [LICENSE](./LICENSE)

## 当前边界

- `table` 目前只展开嵌套 `struct` / `map`，不会展开 `slice` 为多列或多行
- `table` 面向终端阅读，不承诺复杂布局能力
- `table` 列宽目前按简单字符宽度计算，对全角字符未做专门优化
