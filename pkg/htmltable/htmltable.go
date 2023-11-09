package htmltable

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"slices"
	"strconv"
)

type HTMLTable struct {
	dom           *goquery.Selection
	colgroups     int
	width, height int
	Fields        [][]*goquery.Selection
}

func (t *HTMLTable) extractColgroups() bool {
	if t.dom.Find("colgroup").Length() != 1 {
		return false
	}
	t.colgroups = t.dom.Find("colgroup>col").Length()
	return true
}

func (t *HTMLTable) validateColgroups() (err error) {
	t.dom.Find("tbody>tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		rowwidth := 0
		s.Find("td").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			width, err := strconv.Atoi(s.AttrOr("colspan", "1"))
			if err != nil {
				return false
			}
			rowwidth += width
			return true
		})
		if rowwidth > t.colgroups {
			err = errors.New(fmt.Sprintf("rowwidth>tt.colgroups, %d, %d", rowwidth, t.colgroups))
		}
		return err == nil
	})
	return
}

func (t *HTMLTable) extractWidth() (err error) {
	if t.extractColgroups() {
		if err = t.validateColgroups(); err != nil {
			return
		}
		t.width = t.colgroups
	} else {
		widths := make([]int, t.height)
		t.dom.Find("tbody>tr").Each(func(i int, s *goquery.Selection) {
			widths[i] = s.Find("td").Length()
		})
		t.width = slices.Max(widths)
	}
	return
}

func (t *HTMLTable) extractFields() (err error) {
	for i := 0; i < t.height; i++ {
		if err = t.extractRow(i); err != nil {
			return
		}
	}
	return
}

func (t *HTMLTable) extractRow(row int) (err error) {
	col := 0
	t.dom.Find(
		fmt.Sprintf("tbody>tr:nth-child(%d)>td", row+1)).
		EachWithBreak(
			func(_ int, rf *goquery.Selection) bool {
				rowspan, err := strconv.Atoi(rf.AttrOr("rowspan", "1"))
				if err != nil {
					return false
				}
				colspan, err := strconv.Atoi(rf.AttrOr("colspan", "1"))
				if err != nil {
					return false
				}
				t.insertField(row, col, rowspan, colspan, rf)
				col += colspan
				return true
			})
	return
}

func (t *HTMLTable) insertField(row, mincol, rowspan, colspan int, data *goquery.Selection) (err error) {
	if mincol+colspan > t.width {
		err = errors.New(fmt.Sprintf("mincol+colspan>t.width, %d, %d", mincol+colspan, t.width))
		return
	}
	if row+rowspan > t.height {
		err = errors.New(fmt.Sprintf("row+rowspan>t.height, %d, %d", row+rowspan, t.height))
		return
	}
	for ci := mincol; ci+colspan <= t.width; ci++ {
		free := true
		for cii := ci; cii < ci+colspan && free; cii++ {
			for ri := row; ri+rowspan <= t.height && free; ri++ {
				free = t.Fields[ri][cii] == nil
			}
		}
		if free {
			for cii := ci; cii < ci+colspan; cii++ {
				for ri := row; ri < row+rowspan; ri++ {
					t.Fields[ri][cii] = data
				}
			}
			break
		}
	}
	return
}

func (t *HTMLTable) FindMatchingFields(matcher func(*goquery.Selection) bool) (fields []DimensionalField) {
	for ri, row := range t.Fields {
	fieldloop:
		for ci, field := range row {
			for _, kf := range fields {
				if kf.Element == field {
					continue fieldloop
				}
			}
			if !matcher(field) {
				continue fieldloop
			}
			fields = append(fields, getFieldDimensions(t, ri, ci))
		}
	}
	return
}

func ParseHTMLTable(dom *goquery.Selection) (t *HTMLTable, err error) {
	t = &HTMLTable{
		dom: dom,
	}
	t.height = t.dom.Find("tbody>tr").Length()
	if err = t.extractWidth(); err != nil {
		return
	}
	t.Fields = make([][]*goquery.Selection, t.height)
	for i := 0; i < t.height; i++ {
		t.Fields[i] = make([]*goquery.Selection, t.width)
	}
	if err = t.extractFields(); err != nil {
		return
	}
	return
}
