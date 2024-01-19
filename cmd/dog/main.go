package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/yahuian/gox/filex"
)

var (
	//go:embed tmpl/model_tmpl.txt
	modelTmplText string
	//go:embed tmpl/api_tmpl.txt
	apiTmplText string
	//go:embed tmpl/service/create_tmpl.txt
	svcCreateTmplText string
	//go:embed tmpl/service/delete_tmpl.txt
	svcDeleteTmplText string
	//go:embed tmpl/service/detail_tmpl.txt
	svcDetailTmplText string
	//go:embed tmpl/service/list_tmpl.txt
	svcListTmplText string
	//go:embed tmpl/service/update_tmpl.txt
	svcUpdateTmplText string
)

type Param struct {
	PackageName string
	ModelName   string

	// 最终生成的文件
	Filename string

	// 模板
	TmplText string
}

func main() {
	packageName := flag.String("p_name", "", "package name，eg: user, search_history")
	modelName := flag.String("m_name", "", "model name, eg: User, SearchHistory")
	flag.Parse()

	if *packageName == "" && *modelName == "" {
		panic("p_name and m_name can't be empty")
	}

	params := []Param{
		// model
		{
			ModelName: *modelName,
			Filename:  fmt.Sprintf("./model/%s.go", *packageName),
			TmplText:  modelTmplText,
		},

		// api
		{
			PackageName: *packageName,
			ModelName:   *modelName,
			Filename:    fmt.Sprintf("./api/%s.go", *packageName),
			TmplText:    apiTmplText,
		},

		// service
		{
			PackageName: *packageName,
			ModelName:   *modelName,
			Filename:    fmt.Sprintf("./service/%s/%s.go", *packageName, "create"),
			TmplText:    svcCreateTmplText,
		},
		{
			PackageName: *packageName,
			ModelName:   *modelName,
			Filename:    fmt.Sprintf("./service/%s/%s.go", *packageName, "delete"),
			TmplText:    svcDeleteTmplText,
		},
		{
			PackageName: *packageName,
			ModelName:   *modelName,
			Filename:    fmt.Sprintf("./service/%s/%s.go", *packageName, "detail"),
			TmplText:    svcDetailTmplText,
		},
		{
			PackageName: *packageName,
			ModelName:   *modelName,
			Filename:    fmt.Sprintf("./service/%s/%s.go", *packageName, "list"),
			TmplText:    svcListTmplText,
		},
		{
			PackageName: *packageName,
			ModelName:   *modelName,
			Filename:    fmt.Sprintf("./service/%s/%s.go", *packageName, "update"),
			TmplText:    svcUpdateTmplText,
		},
	}

	// 创建 service 目录
	svcPath := "./service/" + *packageName
	exist, err := filex.Exist(svcPath)
	if err != nil {
		panic(err)
	}
	if !exist {
		if err := os.Mkdir(svcPath, 0755); err != nil {
			panic(err)
		}
	}

	for _, v := range params {
		if err := gen(v); err != nil {
			panic(err)
		}
		fmt.Printf("[OK] %s\n", v.Filename)
	}
}

func gen(p Param) error {
	modelTmpl, err := template.New("tmpl").Parse(p.TmplText)
	if err != nil {
		return fmt.Errorf("parse tmpl text err: %w", err)
	}

	exist, err := filex.Exist(p.Filename)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("%s already exists", p.Filename)
	}

	file, err := os.Create(p.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := modelTmpl.Execute(file, p); err != nil {
		return fmt.Errorf("tmpl execute err: %w", err)
	}

	return nil
}
