package hb

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type HtmlOption struct {
	Name  string
	Value string
}

type HtmlPart struct {
	Name     string
	Options  *[]HtmlOption
	SubParts *[]HtmlPart
	Content  string
}

type Script struct {
	*HtmlPart
	Result JsResult
}

type JsResult struct {
	Names []JsResultName
}

type JsResultName struct {
	Name string
}

type HtmlTable struct {
	Id      string
	Titles  *HtmlTableRow
	Alligns []string
	Rows    []*HtmlTableRow
	Scripts []*Script
}

func NewHtmlTable(id, parentid string, titles []string, rows []*HtmlTableRow, alligns []string) *HtmlTable {

	ht := HtmlTable{}
	ht.Id = id
	ht.Scripts = []*Script{}
	al := []string{}
	t := []interface{}{}
	for _, v := range titles {
		t = append(t, v)
	}
	ht.Titles = NewTableRow(t)
	if alligns == nil || len(alligns) < len(titles) {
		count := len(titles) - len(alligns)
		for i := 0; i < count; i++ {
			al = append(al, "center")
		}
		ht.Alligns = al
	} else {
		ht.Alligns = alligns
	}
	for _, v := range rows {
		v.ParentTable = &ht
	}
	ht.Rows = rows
	ht.Titles.ParentTable = &ht
	tableResultType := JsResult{
		Names: []JsResultName{{"table"}},
	}

	for _, v := range titles {
		reducedTitle := strings.Replace(fmt.Sprint(v), " ", "", -1)
		reducedTitle = strings.ToLower(reducedTitle)
		script := NewScript(id+"_"+reducedTitle, "click", parentid, "GET", "/table/"+id+"/"+reducedTitle+"/sort", tableResultType, tableResultType.Names[0].Value())
		ht.Scripts = append(ht.Scripts, script)
	}

	return &ht
}

func (ht HtmlTable) String() string {
	result := "<table>\n"
	result += ht.Titles.asTableHeader()
	for _, v := range ht.Rows {
		result += v.String()
	}
	result += "</table>\n"
	for _, v := range ht.Scripts {
		result += v.String()
	}
	return result
}

type HtmlTableRow struct {
	Row         *[]interface{}
	ParentTable *HtmlTable
}

func NewTableRow(data []interface{}) *HtmlTableRow {

	htr := HtmlTableRow{}
	row := []interface{}{}
	//par := HtmlTable{}
	//htr.ParentTable = &par
	htr.Row = &row
	for _, v := range data {
		*htr.Row = append(*htr.Row, v)
	}
	return &htr
}

func (htr HtmlTableRow) string(rowType string) string {
	tr := NewPart("tr", "", "")

	for i, v := range *htr.Row {
		var th *HtmlPart
		if rowType == "th" {
			th = NewPart(rowType, htr.ParentTable.Id+"_"+strings.Replace(fmt.Sprint(v), " ", "", -1), fmt.Sprint(v))
		} else {
			th = NewPart(rowType, "", fmt.Sprint(v))
		}
		if htr.ParentTable != nil && i < len(htr.ParentTable.Alligns) {
			th.AddOption(&HtmlOption{
				Name:  "align",
				Value: htr.ParentTable.Alligns[i],
			})
		}
		tr.AddSubPart(th)
	}

	return tr.String()
}

func (htr HtmlTableRow) asTableHeader() string {

	return htr.string("th")
}

func (htr HtmlTableRow) String() string {
	return htr.string("td")
}

func (jsResultName JsResultName) Value() string {
	return "result." + jsResultName.Name
}

func NewPart(name, id, content string) *HtmlPart {
	subParts := []HtmlPart{}
	options := []HtmlOption{}
	htmlp := HtmlPart{
		Name:     name,
		Content:  content,
		Options:  &options,
		SubParts: &subParts,
	}
	if len(id) != 0 {
		htmlp.AddOption(&HtmlOption{
			Name:  "id",
			Value: id,
		})
	}
	return &htmlp
}

func NewScript(source, action, target, restType, apicall string, result JsResult, newContent string) *Script {

	part := NewPart("script", "", fmt.Sprintf(`$(document).ready(function(){
    $("#%s").%s(function(){
        $.ajax({type: "%s", url: "%s", async: true, success: function(result){
            $("#%s").html(%s);
        }});
    });
    });`, source, action, restType, apicall, target, newContent))

	script := Script{
		HtmlPart: part,
		Result:   result,
	}
	return &script
}

func (hp *HtmlPart) AddSubPart(subpart *HtmlPart) {

	*hp.SubParts = append(*hp.SubParts, *subpart)
}

func (hp *HtmlPart) AddOptions(options *[]HtmlOption) {

	*hp.Options = append(*hp.Options, *options...)
}
func (hp *HtmlPart) AddOption(option *HtmlOption) {

	*hp.Options = append(*hp.Options, *option)
}

func (hp HtmlPart) String() string {
	result := fmt.Sprintf("<%s", hp.Name)

	for _, v := range *hp.Options {
		result += fmt.Sprintf(" %s=\"%s\"", v.Name, v.Value)
	}

	result += ">"

	result += hp.Content

	for _, v := range *hp.SubParts {
		result += v.String()
	}

	result += fmt.Sprintf("</%s>", hp.Name)

	return result
}

func NewCSSStyle(css string) *HtmlPart {
	part := NewPart("style", "", css)
	return part
}

type Sorter struct {
	sortby []*interface{}
	Data   []*HtmlTableRow
}

func (s *Sorter) Sort(index int) {
	s.sortby = []*interface{}{}
	for _, v := range s.Data {

		s.sortby = append(s.sortby, &(*v.Row)[index])
	}

	if !sort.IsSorted(SortByName(*s)) {
		sort.Sort(SortByName(*s))
	} else {
		sort.Sort(sort.Reverse(SortByName(*s)))
	}

}

type SortByName Sorter

func (a SortByName) Len() int { return len(a.Data) }
func (a SortByName) Swap(i, j int) {
	a.sortby[i], a.sortby[j] = a.sortby[j], a.sortby[i]
	a.Data[i], a.Data[j] = a.Data[j], a.Data[i]
}
func (a SortByName) Less(i, j int) bool {
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
