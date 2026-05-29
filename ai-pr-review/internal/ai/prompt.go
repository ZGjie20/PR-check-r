package ai

import (
	"fmt"
	"strings"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

func BuildReviewPrompt(input model.ReviewInput) string {
	var b strings.Builder

	b.WriteString("你是一名资深代码审查员。请分析以下 Pull Request 中的代码变更片段。\n\n")
	b.WriteString(fmt.Sprintf("PR Title: %s\n", input.PRTitle))
	b.WriteString(fmt.Sprintf("PR Number: %d\n", input.PRNumber))

	if len(input.Commits) > 0 {
		b.WriteString("Recent commit messages:\n")
		limit := len(input.Commits)
		if limit > 5 {
			limit = 5
		}
		for i := 0; i < limit; i++ {
			b.WriteString(fmt.Sprintf("- %s\n", input.Commits[i]))
		}
	}

	chunk := input.Chunk
	b.WriteString(fmt.Sprintf("\nFile: %s\n", chunk.FilePath))
	b.WriteString(fmt.Sprintf("Language: %s\n", chunk.Language))
	if chunk.Function != "" {
		b.WriteString(fmt.Sprintf("Function: %s\n", chunk.Function))
	}
	b.WriteString(fmt.Sprintf("Line range: %d-%d\n", chunk.StartLine, chunk.EndLine))

	if len(chunk.DeletedLines) > 0 {
		b.WriteString("\nDeleted lines:\n")
		for _, line := range chunk.DeletedLines {
			b.WriteString(fmt.Sprintf("- %s\n", line))
		}
	}

	if len(chunk.AddedLines) > 0 {
		b.WriteString("\nAdded lines:\n")
		for _, line := range chunk.AddedLines {
			b.WriteString(fmt.Sprintf("+ %s\n", line))
		}
	}

	b.WriteString(`
审查重点：
- goroutine 泄漏
- 空指针解引用
- 错误处理
- context 传递
- panic 风险
- 并发问题
- SQL 事务问题
- 性能风险

请返回如下结构的 JSON 对象：
{"issues":[{"file":"<filepath>","line":<line_number>,"severity":"high|medium|low","message":"<问题描述>","suggestion":"<修复建议>"}]}

要求：
- message 和 suggestion 必须使用简体中文
- severity 仍使用 high、medium、low（分别表示高、中、低风险）
- 若无问题，返回 {"issues":[]}
- 只返回合法 JSON，不要 markdown 或其他额外文字
- line 指新文件中、且落在上述行号范围内的行号
`)

	return b.String()
}
