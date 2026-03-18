package forgetable

import (
	"hugom/forge/components"
	"hugom/forge/projects"

	"charm.land/bubbles/v2/table"
)

var (
	DiscoveredProjects []projects.Project = projects.DiscoverProjects()
	ProjectsTable      ForgeTable       = initTable()
)

type ColConfig struct {
	Title string
	Width int
}

type ForgeTable struct {
	Table table.Model
}

var tableHeaders = []ColConfig{
	{Title: "Project", Width: 75},
	{Title: "Modified", Width: 15},
	{Title: "Size", Width: 10},
}

func initTable() ForgeTable {
	var columns []table.Column
	for _, col := range tableHeaders {
		columns = append(columns, table.Column{Title: col.Title, Width: col.Width})
	}

	var rows []table.Row
	for _, entry := range DiscoveredProjects {
		rows = append(rows, table.Row{entry.Name, entry.Modified, entry.DirSize})
	}

	tableModel := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	tableModel.SetStyles(getStyle())
	return ForgeTable{Table: tableModel}
}

func (t *ForgeTable) ResizeColumns(termWidth int) {
	var cols []table.Column
	for key, col := range t.Table.Columns() {
		cols = append(cols, table.Column{
			Title: col.Title,
			Width: int(float64(termWidth/2) * float64(tableHeaders[key].Width) / 100),
		})
	}
	t.Table.SetColumns(cols)
}

func (t ForgeTable) Render() string {
	return components.BaseStyle.Render(t.Table.View())
}
