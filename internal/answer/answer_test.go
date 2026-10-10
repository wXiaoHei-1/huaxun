package answer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFillRejectsMismatch(t *testing.T) {
	template := File{
		PackageID: "pkg-08e82c8d",
		Answers:   []Item{{ID: "M001-S009"}, {ID: "M002-S004"}},
	}
	questions := []Question{
		{ID: "M001-S009", Text: "第一题"},
		{ID: "M002-S004", Text: "第二题"},
	}

	t.Run("缺题", func(t *testing.T) {
		_, err := Fill(template, questions[:1], func(string) string { return "ok" }, "")
		if err == nil {
			t.Fatal("缺题应报错")
		}
	})
	t.Run("未知题号", func(t *testing.T) {
		extra := append([]Question{}, questions...)
		extra = append(extra, Question{ID: "M999-S000", Text: "多出来"})
		_, err := Fill(template, extra, func(string) string { return "ok" }, "")
		if err == nil {
			t.Fatal("未知题号应报错")
		}
	})
	t.Run("题号重复", func(t *testing.T) {
		dup := []Question{questions[0], questions[0]}
		_, err := Fill(template, dup, func(string) string { return "ok" }, "")
		if err == nil {
			t.Fatal("重复题号应报错")
		}
	})
}

func TestFillKeepsPackageID(t *testing.T) {
	template := File{
		Name:      "",
		PackageID: "pkg-08e82c8d",
		Answers:   []Item{{ID: "M001-S009"}},
	}
	questions := []Question{{ID: "M001-S009", Text: "查 IP"}}
	got, err := Fill(template, questions, func(text string) string { return "ip==\"" + text + "\"" }, "张三")
	if err != nil {
		t.Fatal(err)
	}
	if got.PackageID != "pkg-08e82c8d" {
		t.Fatalf("参赛包编号被改成 %s", got.PackageID)
	}
	if got.Name != "张三" {
		t.Fatalf("选手名称 got %s", got.Name)
	}
	if got.Answers[0].Query != `ip=="查 IP"` {
		t.Fatalf("查询 got %s", got.Answers[0].Query)
	}

	kept, err := Fill(template, questions, func(string) string { return "port=\"80\"" }, "")
	if err != nil {
		t.Fatal(err)
	}
	if kept.Name != "" {
		t.Fatal("未传入姓名时应保留模板中的空字符串")
	}
}

func TestRealTemplateShape(t *testing.T) {
	root := filepath.Join("..", "..", "readme")
	template, err := LoadTemplate(filepath.Join(root, "答案模板.txt"))
	if err != nil {
		t.Fatal(err)
	}
	questions, err := LoadQuestions(filepath.Join(root, "题目.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if template.PackageID != "pkg-08e82c8d" {
		t.Fatalf("模板参赛包编号 = %s", template.PackageID)
	}
	got, err := Fill(template, questions, func(string) string { return "ok" }, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Answers) != 100 || len(questions) != 100 {
		t.Fatalf("题数 template=%d questions=%d", len(got.Answers), len(questions))
	}
	for i, item := range got.Answers {
		if item.ID != template.Answers[i].ID || item.Query != "ok" {
			t.Fatalf("第 %d 条回填错位: %+v", i, item)
		}
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "答案.json")
	if err := Write(path, got); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadTemplate(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.PackageID != template.PackageID || len(loaded.Answers) != 100 {
		t.Fatalf("写出后读回失败: %+v", loaded)
	}
	if len(body) == 0 || body[0] == 0xEF {
		t.Fatal("答案文件不应为空，也不强制写 BOM")
	}
}
