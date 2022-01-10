package main

import (
	"github.com/olekukonko/tablewriter"
	"io"
	"os"
)

func (tableS *TableStruct) showTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tableS.Header)
	for _, v := range tableS.Data {
		table.Append(v)
	}
	table.Render() // Send output
}

type Table struct {
	w *tablewriter.Table
}

type Tablee interface {
	//Header() []string
	Rows() [][]string
}

func NewTable(w io.Writer) Table {
	return Table{
		tablewriter.NewWriter(w),
	}
}

func (t Table) Render(tab Tablee) {
	//t.w.SetHeader(tab.Header())
	t.w.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	t.w.SetAlignment(tablewriter.ALIGN_LEFT)
	t.w.AppendBulk(tab.Rows())
	t.w.Render()
}
