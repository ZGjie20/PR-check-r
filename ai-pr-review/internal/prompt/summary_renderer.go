package prompt

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

type SummaryRenderer struct {
	system   string
	userTmpl *template.Template
}

func NewSummaryRenderer(templates *SummaryTemplates) (*SummaryRenderer, error) {
	tmpl, err := template.New("summary_user").Parse(templates.User)
	if err != nil {
		return nil, fmt.Errorf("parse summary user prompt template: %w", err)
	}
	return &SummaryRenderer{
		system:   templates.System,
		userTmpl: tmpl,
	}, nil
}

func (r *SummaryRenderer) SystemMessage() string {
	return r.system
}

func (r *SummaryRenderer) RenderUser(input model.SummaryInput) (string, error) {
	var buf bytes.Buffer
	if err := r.userTmpl.Execute(&buf, input); err != nil {
		return "", fmt.Errorf("render summary user prompt: %w", err)
	}
	return buf.String(), nil
}
