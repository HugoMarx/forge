package forgetable

import (
	"charm.land/bubbles/v2/table"
	"hugom/forge/components"
)

type Rowable interface {
	ToRow() []string
}

type ColConfig struct {
	Title string
	Width int
}

type ForgeTable struct {
	headers []ColConfig
	Table   table.Model
	HasData bool
}

var MainTable = &ForgeTable{
	headers: []ColConfig{
		{Title: "Project", Width: 75},
		{Title: "Modified", Width: 15},
		{Title: "Size", Width: 10},
	},
}

var DockerTable = &ForgeTable{
	headers: []ColConfig{
		{Title: "Container", Width: 20},
		{Title: "Image", Width: 15},
		{Title: "State", Width: 10},
		{Title: "Status", Width: 15},
		{Title: "Port", Width: 15},
	},
}

func ToRowable[T Rowable](items []T) []Rowable {
	entries := make([]Rowable, 0, len(items))
	for _, item := range items {
		entries = append(entries, item)
	}
	return entries
}

func (t *ForgeTable) BuildTable(entries []Rowable) {
	var columns []table.Column
	for _, col := range t.headers {
		columns = append(columns, table.Column{Title: col.Title, Width: col.Width})
	}

	var rows []table.Row
	for _, entry := range entries {
		rows = append(rows, entry.ToRow())
	}

	t.HasData = len(rows) != 0

	tableModel := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	tableModel.SetStyles(getStyle())
	t.Table = tableModel
}

func (t *ForgeTable) ResizeColumns(termWidth int) {
	var cols []table.Column
	for key, col := range t.Table.Columns() {
		cols = append(cols, table.Column{
			Title: col.Title,
			Width: int(float64(termWidth/2) * float64(t.headers[key].Width) / 100),
		})
	}
	t.Table.SetColumns(cols)
}

func (t *ForgeTable) Render() string {
	return components.BaseStyle.Render(t.Table.View())
}
