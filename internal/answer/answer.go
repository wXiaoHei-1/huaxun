// Package answer 只负责题目和答案文件，不理解自然语言。
package answer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// Question 是题目文件中的一道题。
type Question struct {
	ID   string `json:"题号"`
	Text string `json:"自然语言输入"`
}

// Item 是答案文件中的一条查询。
type Item struct {
	ID    string `json:"题号"`
	Query string `json:"查询语句"`
}

// File 是可提交的答案文件。
type File struct {
	Name      string `json:"选手名称"`
	PackageID string `json:"参赛包编号"`
	Answers   []Item `json:"答案"`
}

// LoadQuestions 读取题目 JSON 数组。
func LoadQuestions(path string) ([]Question, error) {
	body, err := readJSON(path)
	if err != nil {
		return nil, err
	}
	var questions []Question
	if err := json.Unmarshal(body, &questions); err != nil {
		return nil, fmt.Errorf("解析题目: %w", err)
	}
	return questions, nil
}

// LoadTemplate 读取答案模板。
func LoadTemplate(path string) (File, error) {
	body, err := readJSON(path)
	if err != nil {
		return File{}, err
	}
	var file File
	if err := json.Unmarshal(body, &file); err != nil {
		return File{}, fmt.Errorf("解析答案模板: %w", err)
	}
	return file, nil
}

// Fill 按模板题号回填查询。参赛包编号原样保留。
// name 为空时保留模板里的选手名称。
func Fill(template File, questions []Question, answerOf func(string) string, name string) (File, error) {
	questionIDs, questionText, err := indexQuestions(questions)
	if err != nil {
		return File{}, err
	}
	templateIDs, err := indexTemplate(template)
	if err != nil {
		return File{}, err
	}
	if err := sameSet(questionIDs, templateIDs); err != nil {
		return File{}, err
	}

	out := template
	if name != "" {
		out.Name = name
	}
	out.PackageID = template.PackageID
	out.Answers = make([]Item, len(template.Answers))
	for i, item := range template.Answers {
		text, ok := questionText[item.ID]
		if !ok {
			return File{}, fmt.Errorf("模板题号 %s 在题目中不存在", item.ID)
		}
		query := answerOf(text)
		// 关键步骤：空查询不能写入答案，否则整份文件会被拒绝。
		if query == "" {
			return File{}, fmt.Errorf("题号 %s 的查询语句为空", item.ID)
		}
		out.Answers[i] = Item{ID: item.ID, Query: query}
	}
	return out, nil
}

// Write 以 UTF-8、两个空格缩进写出答案，不转义中文。
func Write(path string, file File) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(file); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func readJSON(path string) ([]byte, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes.TrimPrefix(body, []byte{0xEF, 0xBB, 0xBF}), nil
}

func indexQuestions(questions []Question) (map[string]struct{}, map[string]string, error) {
	ids := make(map[string]struct{}, len(questions))
	text := make(map[string]string, len(questions))
	for _, question := range questions {
		if question.ID == "" {
			return nil, nil, fmt.Errorf("题目存在空题号")
		}
		if _, ok := ids[question.ID]; ok {
			return nil, nil, fmt.Errorf("题目题号重复: %s", question.ID)
		}
		ids[question.ID] = struct{}{}
		text[question.ID] = question.Text
	}
	return ids, text, nil
}

func indexTemplate(file File) (map[string]struct{}, error) {
	ids := make(map[string]struct{}, len(file.Answers))
	for _, item := range file.Answers {
		if item.ID == "" {
			return nil, fmt.Errorf("模板存在空题号")
		}
		if _, ok := ids[item.ID]; ok {
			return nil, fmt.Errorf("模板题号重复: %s", item.ID)
		}
		ids[item.ID] = struct{}{}
	}
	return ids, nil
}

func sameSet(questionIDs, templateIDs map[string]struct{}) error {
	for id := range templateIDs {
		if _, ok := questionIDs[id]; !ok {
			return fmt.Errorf("未知题号: %s", id)
		}
	}
	for id := range questionIDs {
		if _, ok := templateIDs[id]; !ok {
			return fmt.Errorf("缺题: %s", id)
		}
	}
	if len(questionIDs) != len(templateIDs) {
		return fmt.Errorf("题号数量不一致")
	}
	return nil
}
