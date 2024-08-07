package utils

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func CreateStyledTable(columns []string, rows [][]string) string {
	t := table.New().Border(lipgloss.NormalBorder()).BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).Headers(columns...).Rows(rows...).StyleFunc(func(row, col int) lipgloss.Style {
		if row == 0 {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Border(lipgloss.NormalBorder()).BorderTop(false).BorderLeft(false).BorderRight(false).BorderBottom(true).Bold(true)
		}
		if row%2 == 0 {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Bold(true)
		}
		return lipgloss.NewStyle()
	})
	return t.Render()
}
