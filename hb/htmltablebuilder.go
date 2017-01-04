package hb

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/DaKine23/webapp/hb/bsgrid"
	"github.com/DaKine23/webapp/hb/bspagination"
	"github.com/DaKine23/webapp/hb/bstable"
	"github.com/DaKine23/webapp/hb/svg"
)

//HTMLTable represents a HTML table
type HTMLTable struct {
	ID       string
	Titles   *HTMLTableRow
	Alligns  []string
	Rows     []*HTMLTableRow
	Scripts  []*Script
	PageSize int
	Page     int
	Rowcount int
}

//HTMLTableRow represents a HTML tables row
type HTMLTableRow struct {
	Row         *[]interface{}
	ParentTable *HTMLTable
	Status      string
}

//NewHTMLTableContainer should be used as an constructor for table containing HTMLParts
func NewHTMLTableContainer(ht *HTMLTable) *HTMLPart {
	return NewHTMLPart("tablecontainer", ht.ID+"container", ht.String()).AddBootstrapClasses(bsgrid.Cell(12, bsgrid.Large))
}

//NewTableButtonScript adds current page of the table in the end of the apicall
func NewTableButtonScript(source, action, target, table, restType, apicall string, newContent string) *Script {

	return script(source, action, target, restType, apicall+`/"+document.getElementById("`+table+`currentpage").getAttribute("currentpage")`, newContent)
}

//NewHTMLTable should be used as an constructor for *HTMLTable objects
func NewHTMLTable(id string, titles []string, rows []*HTMLTableRow, alligns []string, pagesize, page int) *HTMLTable {

	ht := HTMLTable{}
	ht.PageSize = pagesize
	ht.Rowcount = len(rows)
	lastpage := ht.Rowcount / ht.PageSize
	if ht.Rowcount%ht.PageSize != 0 {
		lastpage++
	}
	if page > lastpage {
		ht.Page = lastpage
	} else {
		ht.Page = page
	}
	ht.ID = id
	ht.Scripts = []*Script{}
	al := []string{}
	t := []interface{}{}
	for _, v := range titles {
		t = append(t, v)
	}
	ht.Titles = NewHTMLTableRow(t)
	if alligns == nil || len(alligns) < len(titles) {
		count := len(titles) - len(alligns)
		for i := 0; i < count; i++ {
			al = append(al, "left")
		}
		ht.Alligns = al
	} else {
		ht.Alligns = alligns
	}
	if pagesize > 0 && ht.Page > 0 {

		end := pagesize * ht.Page
		start := end - pagesize
		if len(rows) < end {
			end = len(rows)
		}
		for i := start; i < end; i++ {
			rows[i].ParentTable = &ht
		}
		ht.Rows = rows[start:end]
	} else {
		for _, v := range rows {
			v.ParentTable = &ht
		}
		ht.Rows = rows
	}

	ht.Titles.ParentTable = &ht

	for _, v := range titles {
		reducedTitle := strings.Replace(fmt.Sprint(v), " ", "", -1)
		reducedTitle = strings.ToLower(reducedTitle)
		script := NewTableButtonScript(id+"_"+reducedTitle, "dblclick", id+"container", id, "GET", "/table/"+id+"/sort/"+reducedTitle, JSONResultValue("table"))
		ht.Scripts = append(ht.Scripts, script)
	}

	return &ht
}

func (htr HTMLTableRow) asTableHeader() string {

	return htr.string("th")
}

//String returns the HTML String for the HTMLTableRow struct
func (htr HTMLTableRow) String() string {
	return htr.string("td")
}

//NewHTMLTableRow should be used as an constructor for *HTMLTableRow objects
func NewHTMLTableRow(data []interface{}, status ...string) *HTMLTableRow {

	htr := HTMLTableRow{}
	row := []interface{}{}
	if len(status) > 0 {
		htr.Status = status[0]
	}
	htr.Row = &row
	for _, v := range data {
		*htr.Row = append(*htr.Row, v)
	}
	return &htr
}

func (htr HTMLTableRow) string(rowType string) string {
	tr := NewHTMLPart("tr", "", "")

	if rowType != "th" && len(htr.Status) > 0 {
		tr.AddBootstrapClasses(htr.Status)
	}

	for i, v := range *htr.Row {
		var th *HTMLPart
		if rowType == "th" {
			th = NewHTMLPart(rowType, htr.ParentTable.ID+"_"+strings.Replace(fmt.Sprint(v), " ", "", -1), fmt.Sprint(v))
		} else {
			th = NewHTMLPart(rowType, "", fmt.Sprint(v))
		}

		if htr.ParentTable != nil && i < len(htr.ParentTable.Alligns) {
			th.AddOption(&HTMLOption{
				Name:  "align",
				Value: htr.ParentTable.Alligns[i],
			})
		}
		tr.addSubPart(th)
	}

	return tr.String(false)
}

//String returns the HTML String for the HTMLTable struct
func (ht HTMLTable) String() string {

	content := ht.Titles.asTableHeader()
	for _, v := range ht.Rows {
		content += v.String()
	}

	table := NewHTMLPart("table", ht.ID, content).AddBootstrapClasses(bstable.Table, bstable.TableHoverRows, bstable.TableStripedRows)

	tabledata := NewHTMLPart("div", ht.ID+"currentpage", "").AddOption(&HTMLOption{
		Name:  "style",
		Value: "display: none;",
	}).AddOption(&HTMLOption{
		Name:  "currentpage",
		Value: fmt.Sprint(ht.Page),
	})

	table.AddSubParts(tabledata)
	for _, v := range ht.Scripts {
		table.AddSubParts(v.HTMLPart)
	}

	container := NewHTMLPart("ul", ht.ID+"buttoncontainer", "")

	if ht.Page > 0 && ht.PageSize > 0 && ht.Rowcount > ht.PageSize {

		lastpage := ht.Rowcount / ht.PageSize
		if ht.Rowcount%ht.PageSize != 0 {
			lastpage++
		}

		end := ht.PageSize * ht.Page
		start := end - ht.PageSize
		if ht.Rowcount < end {
			end = ht.Rowcount
		}
		start++

		list := []*HTMLPart{}
		for i := 0; i < 5; i++ {
			list = append(list, NewHTMLPart("li", "", ""))
		}

		activatedarrowleft := NewSVGIcon(svg.CaretLeft, "#F08532").String(false)
		activatedarrowright := NewSVGIcon(svg.CaretRight, "#F08532").String(false)
		deactivatedarrowleft := NewSVGIcon(svg.CaretLeft, "#707070").String(false)
		deactivatedarrowright := NewSVGIcon(svg.CaretRight, "#707070").String(false)
		var buttonfirst, buttonbefore, buttonlast, buttonnext *HTMLPart
		if ht.Page > 1 {
			buttonfirst = NewHTMLPart("a", ht.ID+"buttonfirst", activatedarrowleft+activatedarrowleft)
			buttonbefore = NewHTMLPart("a", ht.ID+"buttonbefore", activatedarrowleft)

			buttonfirst.addSubPart(pagingScript(buttonfirst.ID, ht.ID, 1).HTMLPart)
			buttonbefore.addSubPart(pagingScript(buttonbefore.ID, ht.ID, ht.Page-1).HTMLPart)
		} else {
			buttonfirst = NewHTMLPart("a", ht.ID+"buttonfirst", deactivatedarrowleft+deactivatedarrowleft)
			buttonbefore = NewHTMLPart("a", ht.ID+"buttonbefore", deactivatedarrowleft)

			list[0].AddBootstrapClasses(bspagination.Disabled)
			list[1].AddBootstrapClasses(bspagination.Disabled)
		}

		if ht.Page*ht.PageSize < ht.Rowcount {
			buttonlast = NewHTMLPart("a", ht.ID+"buttonlast", activatedarrowright+activatedarrowright)
			buttonnext = NewHTMLPart("a", ht.ID+"buttonnext", activatedarrowright)
			buttonlast.addSubPart(pagingScript(buttonlast.ID, ht.ID, lastpage).HTMLPart)
			buttonnext.addSubPart(pagingScript(buttonnext.ID, ht.ID, ht.Page+1).HTMLPart)
		} else {
			buttonlast = NewHTMLPart("a", ht.ID+"buttonlast", deactivatedarrowright+deactivatedarrowright)
			buttonnext = NewHTMLPart("a", ht.ID+"buttonnext", deactivatedarrowright)
			list[3].AddBootstrapClasses(bspagination.Disabled)
			list[4].AddBootstrapClasses(bspagination.Disabled)
		}

		list[0].addSubPart(buttonfirst).AddBootstrapClasses(bspagination.Previous)
		list[1].addSubPart(buttonbefore)
		list[2].addSubPart(
			NewHTMLPart("li", "", "").
				AddBootstrapClasses(bspagination.Disabled).
				AddSubParts(
					NewHTMLPart("a", "", fmt.Sprintf("%d / %d<br>%d - %d / %d ", ht.Page, lastpage, start, end, ht.Rowcount)),
				),
		)
		list[3].addSubPart(buttonnext)
		list[4].addSubPart(buttonlast).AddBootstrapClasses(bspagination.Next)
		container.AddSubParts(list...).AddBootstrapClasses(bspagination.Pager, bspagination.Small)

	}
	result := table.String(false)

	if len(*container.SubParts) > 0 || len(container.Content) > 0 {
		result += container.String(false)
	}

	return result
}

func pagingScript(buttonID, tableID string, page int) *Script {

	return NewScript(
		buttonID,
		"click",
		tableID+"container",
		"GET",
		"/table/"+tableID+"/show/"+fmt.Sprint(page),
		JSONResultValue("table"),
	)
}

//Sort sorts data for your tables. First execution sorts ascending second sorts descending
func Sort(index int, data []*HTMLTableRow) []*HTMLTableRow {

	s := sorter{}
	s.sortby = []*interface{}{}
	s.Data = data
	for _, v := range s.Data {

		s.sortby = append(s.sortby, &(*v.Row)[index])
	}

	if !sort.IsSorted(tablesort(s)) {
		sort.Sort(tablesort(s))
	} else {
		sort.Sort(sort.Reverse(tablesort(s)))
	}
	return s.Data
}

type sorter struct {
	sortby []*interface{}
	Data   []*HTMLTableRow
}
type tablesort sorter

func (a tablesort) Len() int { return len(a.Data) }
func (a tablesort) Swap(i, j int) {
	a.sortby[i], a.sortby[j] = a.sortby[j], a.sortby[i]
	a.Data[i], a.Data[j] = a.Data[j], a.Data[i]
}
func (a tablesort) Less(i, j int) bool {
	result := false

	theint1, isint1 := (*a.sortby[i]).(int)
	theint2, isint2 := (*a.sortby[j]).(int)
	thefloat1, isfloat1 := (*a.sortby[i]).(float64)
	thefloat2, isfloat2 := (*a.sortby[j]).(float64)
	thetime1, istime1 := (*a.sortby[i]).(time.Time)
	thetime2, istime2 := (*a.sortby[j]).(time.Time)
	switch {
	case isint1 && isint2:
		result = theint1 < theint2
	case isfloat1 && isfloat2:
		result = thefloat1 < thefloat2
	case istime1 && istime2:
		result = thetime1.Before(thetime2)
	default:
		result = fmt.Sprint(*a.sortby[i]) < fmt.Sprint((*a.sortby[j]))
	}
	return result
}
