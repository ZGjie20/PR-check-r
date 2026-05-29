package prompt

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

type Renderer struct {
	system   string
	userTmpl *template.Template
}

func NewRenderer(templates *ReviewTemplates) (*Renderer, error) {
	tmpl, err := template.New("review_user").Parse(templates.User)
	if err != nil {
		return nil, fmt.Errorf("parse user prompt template: %w", err)
	}
	return &Renderer{
		system:   templates.System,
		userTmpl: tmpl,
	}, nil
}

func (r *Renderer) SystemMessage() string {
	return r.system
}

func (r *Renderer) RenderUser(input model.ReviewInput) (string, error) {
	var buf bytes.Buffer
	if err := r.userTmpl.Execute(&buf, input); err != nil {
		return "", fmt.Errorf("render user prompt: %w", err)
	}
	return buf.String(), nil
}
