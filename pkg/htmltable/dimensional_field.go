package htmltable

import "github.com/PuerkitoBio/goquery"

type DimensionalField struct {
	Element       *goquery.Selection
	Row, Column   int
	Height, Width int
}

func getFieldDimensions(t *HTMLTable, row, column int) (df DimensionalField) {
	df.Element = t.Fields[row][column]
	df.Row = row
	for ; t.Fields[df.Row-1][df.Column] == df.Element; df.Row-- {
	}
	df.Column = column
	for ; t.Fields[df.Row][df.Column-1] == df.Element; df.Column-- {
	}
	df.Height = 1
	for ; t.Fields[df.Row+df.Height][df.Column] == df.Element; df.Height++ {
	}
	df.Width = 1
	for ; t.Fields[df.Row][df.Column+df.Width] == df.Element; df.Width++ {
	}

	return
}
