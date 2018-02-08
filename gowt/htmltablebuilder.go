package gowt

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.bus.zalan.do/ale/gowt/bspagination"
	"github.bus.zalan.do/ale/gowt/bstable"
	"github.bus.zalan.do/ale/gowt/faicons"
	"github.bus.zalan.do/ale/gowt/jqaction"
	"github.bus.zalan.do/ale/gowt/svg"
)

//HTMLTable represents a HTML table
type HTMLTable struct {
	ID       string
	Titles   *HTMLTableRow
	Alligns  []string
	Rows     []HTMLTableRow
	Scripts  []HTMLPart
	PageSize int
	Page     int
	Rowcount int
	Sortable bool
	SortedBy string
	Asc      *bool
}

//HTMLTableRow represents a HTML tables row
type HTMLTableRow struct {
	Row         []interface{}
	ParentTable *HTMLTable
	Status      string
}

//NewHTMLTableContainer should be used as an constructor for table containing HTMLParts
func NewHTMLTableContainer(ht *HTMLTable) *HTMLPart {
	return NewHTMLPart("div", ht.ID+"container", ht.String())

}

//NewTableButtonScript adds current page of the table in the end of the apicall
func NewTableButtonScript(source, action, target, table, restType, apicall string, newContent string) *HTMLPart {

	uriFormat := `/"+document.getElementById("%scurrentpage").getAttribute("currentpage")`
	return script(source, action, target, restType, apicall+fmt.Sprintf(uriFormat, table), newContent)
}

//ToSortableTable shows table as sortable (adds icons to the header and scripts that call the tables sort API)
func (ht *HTMLTable) ToSortableTable(sortedby string, asc *bool) *HTMLTable {

	ht.SortedBy = sortedby
	ht.Asc = asc
	ht.Sortable = true

	for _, v := range ht.Titles.Row {
		reducedTitle := fmt.Sprint(v)
		reducedTitle = strings.Replace(reducedTitle, " ", "", -1)
		reducedTitle = strings.ToLower(reducedTitle)
		uriFormat := "/table/%s/sort/%s"
		script := NewTableButtonScript(ht.ID+"_"+reducedTitle, jqaction.Click, ht.ID+"container", ht.ID, "GET", fmt.Sprintf(uriFormat, ht.ID, reducedTitle), JSONResultValue("table"))
		ht.Scripts = append(ht.Scripts, *script)
	}

	return ht
}

//NewHTMLTable should be used as an constructor for *HTMLTable objects
func NewHTMLTable(id string, titles []string, rows []HTMLTableRow, alligns []string, pagesize, page int) *HTMLTable {

	ht := HTMLTable{}
	ht.PageSize = pagesize
	ht.Rowcount = len(rows)
	if ht.PageSize > 0 {

		lastpage := ht.Rowcount / ht.PageSize
		if ht.Rowcount%ht.PageSize != 0 {
			lastpage++
		}
		if page > lastpage {
			ht.Page = lastpage
		} else {
			ht.Page = page
		}
	}
	if ht.Page == 0 {
		ht.Page = 1
	}
	ht.ID = id
	ht.Scripts = []HTMLPart{}
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
	if ht.PageSize > 0 {

		end := ht.PageSize * ht.Page
		start := end - ht.PageSize
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

	return &ht
}

func (htr HTMLTableRow) asTableHeader(sortable bool, sortedBy string, asc *bool) string {

	return htr.string("th", sortable, sortedBy, asc)
}

//String returns the HTML String for the HTMLTableRow struct
func (htr HTMLTableRow) String() string {
	return htr.string("td", false, "", nil)
}

//NewHTMLTableRow should be used as an constructor for *HTMLTableRow objects
func NewHTMLTableRow(data []interface{}, status ...string) *HTMLTableRow {

	htr := HTMLTableRow{}
	row := []interface{}{}
	if len(status) > 0 {
		htr.Status = status[0]
	}
	htr.Row = row
	for _, v := range data {
		htr.Row = append(htr.Row, v)
	}
	return &htr
}

func (htr HTMLTableRow) string(rowType string, sortable bool, sortedby string, asc *bool) string {

	tr := NewHTMLPart("tr", "")

	if rowType != "th" && len(htr.Status) > 0 {
		tr.AddBootstrapClasses(htr.Status)
	}

	for i, v := range htr.Row {
		var th *HTMLPart
		if rowType == "th" {
			title := fmt.Sprint(v)
			if sortable {
				title = addSortIcon(title, sortedby, asc)
			}
			th = NewHTMLPart(rowType, htr.ParentTable.ID+"_"+strings.Replace(strings.ToLower(fmt.Sprint(v)), " ", "", -1), title)
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

	return tr.String()
}
func addSortIcon(title, sortedby string, asc *bool) string {
	reducedTitle := title
	reducedTitle = strings.Replace(title, " ", "", -1)
	reducedTitle = strings.ToLower(reducedTitle)

	top := NewFontAwesomeIconDefinition(faicons.SortAsc, faicons.ModifyRegularStackIcon, "text-muted")
	down := NewFontAwesomeIconDefinition(faicons.SortDesc, faicons.ModifyRegularStackIcon, "text-muted")
	icon := NewFontAwesomeIcon(top, down)
	if reducedTitle == sortedby && asc != nil {
		if *asc {
			top = NewFontAwesomeIconDefinition(faicons.SortAsc, faicons.ModifyRegularStackIcon, "text-warning")
			icon = NewFontAwesomeIcon(top, down)
		} else {
			down = NewFontAwesomeIconDefinition(faicons.SortDesc, faicons.ModifyRegularStackIcon, "text-warning")
			icon = NewFontAwesomeIcon(top, down)
		}
	}

	return title + icon.String()
}

//String returns the HTML String for the HTMLTable struct
func (ht HTMLTable) String() string {

	innercontainer := NewHTMLPart("div", "").AddBootstrapClasses(bstable.TableResponsiveTable)

	content := ht.Titles.asTableHeader(ht.Sortable, ht.SortedBy, ht.Asc)
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
		table.AddScripts(&v)
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

		activatedarrowleft := NewSVGIcon(svg.CaretLeft, "#FF6900").String()
		activatedarrowright := NewSVGIcon(svg.CaretRight, "#FF6900").String()
		deactivatedarrowleft := NewSVGIcon(svg.CaretLeft, "#707070").String()
		deactivatedarrowright := NewSVGIcon(svg.CaretRight, "#707070").String()
		var buttonfirst, buttonbefore, buttonlast, buttonnext *HTMLPart
		if ht.Page > 1 {
			buttonfirst = NewHTMLPart("a", ht.ID+"buttonfirst", activatedarrowleft+activatedarrowleft)
			buttonbefore = NewHTMLPart("a", ht.ID+"buttonbefore", activatedarrowleft)

			buttonfirst.addSubPart(pagingScript(buttonfirst.ID, ht.ID, 1))
			buttonbefore.addSubPart(pagingScript(buttonbefore.ID, ht.ID, ht.Page-1))
		} else {
			buttonfirst = NewHTMLPart("a", ht.ID+"buttonfirst", deactivatedarrowleft+deactivatedarrowleft)
			buttonbefore = NewHTMLPart("a", ht.ID+"buttonbefore", deactivatedarrowleft)

			list[0].AddBootstrapClasses(bspagination.Disabled)
			list[1].AddBootstrapClasses(bspagination.Disabled)
		}

		if ht.Page*ht.PageSize < ht.Rowcount {
			buttonlast = NewHTMLPart("a", ht.ID+"buttonlast", activatedarrowright+activatedarrowright)
			buttonnext = NewHTMLPart("a", ht.ID+"buttonnext", activatedarrowright)
			buttonlast.addSubPart(pagingScript(buttonlast.ID, ht.ID, lastpage))
			buttonnext.addSubPart(pagingScript(buttonnext.ID, ht.ID, ht.Page+1))
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

	innercontainer.AddSubParts(table)
	result := innercontainer.StringWithScripts()

	if len(*container.SubParts) > 0 || len(container.Content) > 0 {
		result += container.String()
	}

	return result
}

func pagingScript(buttonID, tableID string, page int) *HTMLPart {

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
func Sort(index int, data []HTMLTableRow) []HTMLTableRow {

	s := sorter{}
	s.sortby = []interface{}{}
	s.Data = data
	for _, v := range s.Data {

		s.sortby = append(s.sortby, &(v.Row)[index])
	}

	if !sort.IsSorted(tablesort(s)) {
		sort.Sort(tablesort(s))
	} else {
		sort.Sort(sort.Reverse(tablesort(s)))
	}
	return s.Data
}

type sorter struct {
	sortby []interface{}
	Data   []HTMLTableRow
}
type tablesort sorter

func (a tablesort) Len() int { return len(a.Data) }
func (a tablesort) Swap(i, j int) {
	a.sortby[i], a.sortby[j] = a.sortby[j], a.sortby[i]
	a.Data[i], a.Data[j] = a.Data[j], a.Data[i]
}
func (a tablesort) Less(i, j int) bool {
	result := false

	theint1, isint1 := a.sortby[i].(int)
	theint2, isint2 := a.sortby[j].(int)
	thefloat1, isfloat1 := a.sortby[i].(float64)
	thefloat2, isfloat2 := a.sortby[j].(float64)
	thetime1, istime1 := a.sortby[i].(time.Time)
	thetime2, istime2 := a.sortby[j].(time.Time)
	switch {
	case isint1 && isint2:
		result = theint1 < theint2
	case isfloat1 && isfloat2:
		result = thefloat1 < thefloat2
	case istime1 && istime2:
		result = thetime1.Before(thetime2)
	default:
		result = fmt.Sprint(a.sortby[i]) < fmt.Sprint((a.sortby[j]))
	}
	return result
}
