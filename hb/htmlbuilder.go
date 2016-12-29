package hb

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/DaKine23/webapp/hb/bsgrid"
	"github.com/DaKine23/webapp/hb/bspagination"
	"github.com/DaKine23/webapp/hb/bstable"
)

//HTMLOption represents an HTMLPart Option
type HTMLOption struct {
	Name  string
	Value string
}

//HTMLPart represents a general HTML Tag and its contents
type HTMLPart struct {
	ID       string
	Class    string
	Options  *[]HTMLOption
	SubParts *[]HTMLPart
	Content  string
}

//Script represents an HTMLPart containing a Script (mostly to do an Ajax call)
type Script struct {
	*HTMLPart
	Result JSONResult
}

//JSONResult represents the JSON schema in a (js) Script
type JSONResult struct {
	Names []JSONResultName
}

//JSONResultName represents the JSON variable in a (js) Script
type JSONResultName struct {
	Name string
}

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

var tableResultType = JSONResult{
	Names: []JSONResultName{{"table"}},
}

//NewHTMLPart should be used as an constructor for *HTMLPart objects
func NewHTMLPart(class, id, content string) *HTMLPart {
	subParts := []HTMLPart{}
	options := []HTMLOption{}
	htmlp := HTMLPart{
		Class:    class,
		Content:  content,
		Options:  &options,
		SubParts: &subParts,
		ID:       id,
	}
	if len(id) != 0 {
		htmlp.AddOption(&HTMLOption{
			Name:  "id",
			Value: id,
		})
	}
	return &htmlp
}

//NewHTMLTableContainer should be used as an constructor for table containing HTMLParts
func NewHTMLTableContainer(ht *HTMLTable) *HTMLPart {
	return NewHTMLPart("tablecontainer", ht.ID+"container", ht.String()).AddBootstrapClasses(bsgrid.Cell(12, bsgrid.Large))
}

//NewScript should be used as an constructor for *Script objects
func NewScript(source, action, target, restType, apicall string, result JSONResult, newContent string) *Script {

	return script(source, action, target, restType, apicall+`"`, result, newContent)
}

func script(source, action, target, restType, apicall string, result JSONResult, newContent string) *Script {

	part := NewHTMLPart("script", "", fmt.Sprintf(`$(document).ready(function(){
    $("#%s").%s(function(){
        $.ajax({type: "%s", url: "%s, async: true, success: function(result){
            $("#%s").html(%s);
        }});
    });
    });`, source, action, restType, apicall, target, newContent))

	script := Script{
		HTMLPart: part,
		Result:   result,
	}
	return &script
}

//NewTableButtonScript adds current page of the table in the end of the apicall
func NewTableButtonScript(source, action, target, table, restType, apicall string, result JSONResult, newContent string) *Script {

	return script(source, action, target, restType, apicall+`/"+document.getElementById("`+table+`currentpage").getAttribute("currentpage")`, result, newContent)
}

//NewHTMLTable should be used as an constructor for *HTMLTable objects
func NewHTMLTable(id string, titles []string, rows []*HTMLTableRow, alligns []string, pagesize, page int) *HTMLTable {

	ht := HTMLTable{}
	ht.Page = page
	ht.PageSize = pagesize
	ht.Rowcount = len(rows)
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
	if pagesize > 0 && page > 0 {

		end := pagesize * page
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
		script := NewTableButtonScript(id+"_"+reducedTitle, "click", id+"container", id, "GET", "/table/"+id+"/sort/"+reducedTitle, tableResultType, tableResultType.Names[0].Value())
		ht.Scripts = append(ht.Scripts, script)
	}

	return &ht
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

//String returns the HTML String for the HTMLPart struct includes all subparts subsubparts ...
func (hp HTMLPart) String() string {
	result := fmt.Sprintf("<%s", hp.Class)

	for _, v := range *hp.Options {
		result += fmt.Sprintf(" %s=\"%s\"", v.Name, v.Value)
	}

	result += ">"

	result += hp.Content

	for _, v := range *hp.SubParts {
		result += v.String()
	}

	result += fmt.Sprintf("</%s>", hp.Class)

	return result
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

		list := []*HTMLPart{}
		for i := 0; i < 5; i++ {
			list = append(list, NewHTMLPart("li", "", ""))
		}
		buttonfirst := NewHTMLPart("a", ht.ID+"buttonfirst", "<<")
		buttonbefore := NewHTMLPart("a", ht.ID+"buttonbefore", "<")
		if ht.Page > 1 {
			buttonfirst.addSubPart(pagingButtonScript(buttonfirst.ID, ht.ID, 1).HTMLPart)
			buttonbefore.addSubPart(pagingButtonScript(buttonbefore.ID, ht.ID, ht.Page-1).HTMLPart)
		} else {
			list[0].AddBootstrapClasses(bspagination.Disabled)
			list[1].AddBootstrapClasses(bspagination.Disabled)
		}

		buttonlast := NewHTMLPart("a", ht.ID+"buttonlast", ">>")
		buttonnext := NewHTMLPart("a", ht.ID+"buttonnext", ">")
		if ht.Page*ht.PageSize < ht.Rowcount {

			buttonlast.addSubPart(pagingButtonScript(buttonlast.ID, ht.ID, lastpage).HTMLPart)
			buttonnext.addSubPart(pagingButtonScript(buttonnext.ID, ht.ID, ht.Page+1).HTMLPart)
		} else {
			list[3].AddBootstrapClasses(bspagination.Disabled)
			list[4].AddBootstrapClasses(bspagination.Disabled)
		}

		list[0].addSubPart(buttonfirst).AddBootstrapClasses(bspagination.Previous)
		list[1].addSubPart(buttonbefore)
		list[2].addSubPart(
			NewHTMLPart("li", "", "").
				AddBootstrapClasses(bspagination.Disabled).
				addSubPart(
					NewHTMLPart("a", "", fmt.Sprintf("%d / %d (%d)", ht.Page, lastpage, ht.Rowcount)),
				),
		)
		list[3].addSubPart(buttonnext)
		list[4].addSubPart(buttonlast).AddBootstrapClasses(bspagination.Next)
		container.AddSubParts(list...).AddBootstrapClasses(bspagination.Pager, bspagination.Small)

	}
	result := table.String()

	if len(*container.SubParts) > 0 || len(container.Content) > 0 {
		result += container.String()
	}

	return result
}

// GET http://127.0.0.1:3000/index.html
func pagingButtonScript(buttonID, tableID string, page int) *Script {

	return NewScript(
		buttonID,
		"click",
		tableID+"container",
		"GET",
		"/table/"+tableID+"/show/"+fmt.Sprint(page),
		tableResultType,
		tableResultType.Names[0].Value(),
	)
}

func (htr HTMLTableRow) asTableHeader() string {

	return htr.string("th")
}

//String returns the HTML String for the HTMLTableRow struct
func (htr HTMLTableRow) String() string {
	return htr.string("td")
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

	return tr.String()
}

//Value returns the Resultjsons variable as js ready string "result.<myvalue>"
func (jsResultName JSONResultName) Value() string {
	return "result." + jsResultName.Name
}

func (hp *HTMLPart) addSubPart(subpart *HTMLPart) *HTMLPart {

	*hp.SubParts = append(*hp.SubParts, *subpart)
	return hp
}

//AddSubParts adds one or more HTMLParts (subparts) in your HTMLPart
func (hp *HTMLPart) AddSubParts(subparts ...*HTMLPart) *HTMLPart {
	for _, v := range subparts {
		*hp.SubParts = append(*hp.SubParts, *v)
	}
	return hp
}

//AddOptions adds Options to you HTMLParts
func (hp *HTMLPart) AddOptions(options *[]HTMLOption) *HTMLPart {

	*hp.Options = append(*hp.Options, *options...)
	return hp
}

//AddOption adds an option to you HTMLParts if name is already there it concats the Value to the existing option
func (hp *HTMLPart) AddOption(option *HTMLOption) *HTMLPart {

	for i, v := range *hp.Options {
		if v.Name == option.Name {

			newOption := HTMLOption{
				Name:  option.Name,
				Value: v.Value + " " + option.Value,
			}
			*hp.Options = append((*hp.Options)[:i], (*hp.Options)[i+1:]...)
			*hp.Options = append(*hp.Options, newOption)
			return hp

		}
	}

	*hp.Options = append(*hp.Options, *option)
	return hp
}

//SetOption adds an option to you HTMLParts if name is already there it replaces the Value of the existing option
func (hp *HTMLPart) SetOption(option *HTMLOption) *HTMLPart {

	for i, v := range *hp.Options {
		if v.Name == option.Name {

			*hp.Options = append((*hp.Options)[:i], (*hp.Options)[i+1:]...)
			*hp.Options = append(*hp.Options, *option)
			return hp

		}
	}

	*hp.Options = append(*hp.Options, *option)
	return hp
}

//AddBootstrapClasses adds Bootstrap class as and Option to the HTMLPart. Dont forget to add the base class if existent
func (hp *HTMLPart) AddBootstrapClasses(classes ...string) *HTMLPart {
	hp.AddOption(&HTMLOption{
		Name:  "class",
		Value: strings.Join(classes, " "),
	})
	return hp
}

//NewCSSStyle creates a HTMLPart for plain CSS styles
func NewCSSStyle(css string) *HTMLPart {
	part := NewHTMLPart("style", "", css)
	return part
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
