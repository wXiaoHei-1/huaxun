package main

import (
	"flag"
	"fmt"
	"os"

	"huaxun/internal/answer"
	"huaxun/internal/convert"
)

func main() {
	name := flag.String("name", "", "选手名称，不传则保留模板中的空字符串")
	questionsPath := flag.String("questions", "readme/题目.txt", "题目文件")
	templatePath := flag.String("template", "readme/答案模板.txt", "答案模板")
	outPath := flag.String("out", "答案.json", "输出的答案文件")
	flag.Parse()

	if err := run(*name, *questionsPath, *templatePath, *outPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(name, questionsPath, templatePath, outPath string) error {
	questions, err := answer.LoadQuestions(questionsPath)
	if err != nil {
		return err
	}
	template, err := answer.LoadTemplate(templatePath)
	if err != nil {
		return err
	}
	filled, err := answer.Fill(template, questions, convert.ToQuery, name)
	if err != nil {
		return err
	}
	return answer.Write(outPath, filled)
}
