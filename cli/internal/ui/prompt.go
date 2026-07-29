package ui

import (
	"context"
	"io"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type Prompt struct {
	input      io.Reader
	output     io.Writer
	accessible bool
	noColor    bool
}

func NewPrompt(input io.Reader, output io.Writer, accessible, noColor bool) *Prompt {
	return &Prompt{input: input, output: output, accessible: accessible, noColor: noColor}
}

func (p *Prompt) Select(ctx context.Context, title string, options ...huh.Option[string]) (string, error) {
	var selected string
	form := huh.NewForm(huh.NewGroup(huh.NewSelect[string]().Title(title).Options(options...).Value(&selected)))
	err := p.configure(form).RunWithContext(ctx)
	return selected, err
}

func (p *Prompt) Confirm(ctx context.Context, title string, affirmative bool) (bool, error) {
	value := affirmative
	form := huh.NewForm(huh.NewGroup(huh.NewConfirm().Title(title).Affirmative("Yes").Negative("No").Value(&value)))
	err := p.configure(form).RunWithContext(ctx)
	return value, err
}

func (p *Prompt) configure(form *huh.Form) *huh.Form {
	return form.WithInput(p.input).WithOutput(p.output).WithAccessible(p.accessible).WithTheme(theme(p.noColor))
}

func theme(noColor bool) *huh.Theme {
	if noColor {
		return huh.ThemeBase()
	}

	value := huh.ThemeCharm()
	green := lipgloss.AdaptiveColor{Light: "#00875A", Dark: "#5FAF5F"}
	muted := lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#808080"}

	value.Focused.Title = value.Focused.Title.Foreground(lipgloss.NoColor{}).Bold(false)
	value.Focused.SelectSelector = value.Focused.SelectSelector.Foreground(green).SetString("● ")
	value.Focused.SelectedOption = value.Focused.SelectedOption.Foreground(green)
	value.Focused.UnselectedOption = value.Focused.UnselectedOption.Foreground(muted)
	value.Focused.Base = value.Focused.Base.BorderForeground(muted)
	value.Focused.Card = value.Focused.Base
	value.Blurred = value.Focused
	value.Group.Title = value.Focused.Title
	return value
}
