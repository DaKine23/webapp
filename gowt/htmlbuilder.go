package gowt

import (
	"fmt"
	"strings"

	"github.bus.zalan.do/ale/gowt/bscontainer"
	"github.bus.zalan.do/ale/gowt/bsglyphicons"
	"github.bus.zalan.do/ale/gowt/bsgrid"
	"github.bus.zalan.do/ale/gowt/faicons"
	"github.bus.zalan.do/ale/gowt/svg"
)

//HTMLOption represents an HTMLPart Option
type HTMLOption struct {
	Name  string
	Value string
}

type BsGrid struct {
	Grid *[][]BsCell
	ID   string
}

type BsCell struct {
	Content *HTMLPart
	Colspan int
}

//NewDefaultPage creates a default page to start with. Expects bootstrap.min.css, font-awesome.min.css, jquery.min.js, bootstrap.min.js and all the additional css and js files in the /static/ path
func NewDefaultPage(title string, jsFiles, cssFiles []string) (html, head, body *HTMLPart) {

	//define <html>
	html = NewHTMLPart("html", "", "")

	//define <head> with title and include bootstrap

	headcontent := `<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	<link rel="stylesheet" href="/static/bootstrap.min.css">
	<link rel="stylesheet" href="/static/font-awesome.min.css">`
	if len(cssFiles) > 0 {
		for _, v := range cssFiles {
			headcontent += `
			<link rel="stylesheet" href="` + v + `">`
		}
	}
	head = NewHTMLPart("head", "", headcontent).
		AddSubParts(NewHTMLPart("title", "", title))
		//svg icon size
	head.AddSubParts(NewHTMLPart("style", "", `.icon {
  width: 14px;
  height: 14px;
}
	`))
	//define js libraries you want to import
	jsLibraries := []string{
		"/static/jquery.min.js",
		"/static/bootstrap.min.js",
	}

	getCookies := NewHTMLPart("script", "", `function getCookieValue(a) {
    var b = document.cookie.match('(^|;)\\s*' + a + '\\s*=\\s*([^;]+)');
    return b ? b.pop() : '';
	}`)
	head.addSubPart(getCookies)
	if len(jsFiles) > 0 {
		jsLibraries = append(jsLibraries, jsFiles...)
	}
	//add js libraries you want to import to the <head>
	for _, v := range jsLibraries {
		jsLibrariesPart := NewHTMLPart("script", "").AddOption(&HTMLOption{
			Name:  "src",
			Value: v,
		})
		head.AddSubParts(jsLibrariesPart)
	}

	//define <body>
	body = NewHTMLPart("body", "", svg.Iconset)

	// add <head> and <body> to <html>
	html.AddSubParts(head, body)

	return html, head, body

}

//HTMLPart returns the coresponding HTMLPart for the Grid
func (bsg BsGrid) HTMLPart() *HTMLPart {

	container := NewHTMLPart("div", bsg.ID).AddBootstrapClasses(bscontainer.ContainerFluid)

	for _, v := range *bsg.Grid {
		row := NewHTMLPart(bsgrid.Row, "").AddBootstrapClasses(bsgrid.Row)
		autocolspan := 12

		colspancounter := len(v)

		for _, v2 := range v {
			if v2.Colspan > 0 {
				autocolspan -= v2.Colspan
				colspancounter--
			}
		}
		if colspancounter > 0 {
			autocolspan = autocolspan / colspancounter
		}
		for _, v2 := range v {
			cell := NewHTMLPart("cell", "")

			switch {
			case v2.Content == nil && v2.Colspan > 0:
				cell.AddBootstrapClasses(bsgrid.Cell(v2.Colspan, bsgrid.Large))
			case v2.Content == nil && v2.Colspan <= 0:
				cell.AddBootstrapClasses(bsgrid.Cell(autocolspan, bsgrid.Large))
			case v2.Colspan > 0:
				cell.AddBootstrapClasses(bsgrid.Cell(v2.Colspan, bsgrid.Large)).
					addSubPart(v2.Content)
			case v2.Colspan <= 0:
				cell.AddBootstrapClasses(bsgrid.Cell(autocolspan, bsgrid.Large)).
					addSubPart(v2.Content)
			}

			row.addSubPart(cell)
		}
		container.addSubPart(row)
	}
	return container

}

//HTMLPart represents a general HTML Tag and its contents
type HTMLPart struct {
	ID       string
	Class    string
	Options  *[]HTMLOption
	SubParts *[]HTMLPart
	Scripts  *[]HTMLPart
	Content  string
}

//NewHTMLPart should be used as an constructor for *HTMLPart objects
func NewHTMLPart(class, id string, content ...string) *HTMLPart {
	subParts := []HTMLPart{}
	scripts := []HTMLPart{}
	options := []HTMLOption{}
	concatcontent := []byte{}
	for _, v := range content {
		concatcontent = append(concatcontent, v...)

	}
	htmlp := HTMLPart{
		Class:    class,
		Content:  string(concatcontent),
		Options:  &options,
		SubParts: &subParts,
		Scripts:  &scripts,
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

//NewScript should be used as an constructor for *Script objects
func NewScript(source, action, target, restType, apicall, newContent string) *HTMLPart {

	return script(source, action, target, restType, apicall+`"`, newContent)
}

func NewModalScript(source, action, target, restType, apicall, newContent, dialogID, message, postitiveRequest, negativeRequest string) *HTMLPart {

	return modalRequestScript(source, action, target, restType, apicall+`"`, newContent, dialogID, message, postitiveRequest, negativeRequest)

}

func modalRequestScript(source, action, target, restType, apicall, newContent, dialogID, message, postitiveRequest, negativeRequest string) *HTMLPart {

	call := fmt.Sprintf(`  $.ajax({type: "%s", url: "%s, async: false, success: function(result){
            $("#%s").html(%s);
        }});`, restType, apicall, target, newContent)

	modal := fmt.Sprintf(`$("#%s").dialog({
				modal: true,
				resizable: false,
				buttons: {
					"%s": function() {
						%s
						$(this).dialog("close");
					},
					"%s": function() {
						$(this).dialog("close");
					}
				}
			});`, dialogID, postitiveRequest, call, negativeRequest)

	return NewHTMLPart("script", "", fmt.Sprintf(`$(document).ready(function(){

	$("#%s").%s(function(){
		$("#%s").text("%s");
		%s
    });
	});`, source, action, dialogID, message, modal))

}

func script(source, action, target, restType, apicall, newContent string) *HTMLPart {

	part := NewHTMLPart("script", "", fmt.Sprintf(`	$(document).ready(function(){
    $("#%s").%s(function(){
        $.ajax({type: "%s", url: "%s, async: false, success: function(result){
            $("#%s").html(%s);
        }});
    });
    });`, source, action, restType, apicall, target, newContent))

	return part
}

//AddTooltip adds a Tooltip you need to add TooltipScript() once to your page as a script
func (hp *HTMLPart) AddTooltip(tip, placement string) *HTMLPart {

	hp.AddOption(&HTMLOption{"data-toggle", "tooltip"})
	if len(placement) > 0 {
		hp.AddOption(&HTMLOption{"data-placement", placement})
	}
	hp.AddOption(&HTMLOption{"title", tip})

	return hp

}

//TooltipScript append the returned Script once as a script to you result if you want to see Tooltips
func TooltipScript() *HTMLPart {

	return NewHTMLPart("script", "", `$(document).ready(function(){$('[data-toggle="tooltip"]').tooltip();});`).AddOption(&HTMLOption{"type", "text/javascript"})

}

type FontAwesomeIconDefinition struct {
	Classes []string
}

func NewFontAwesomeIconDefinition(definition ...string) FontAwesomeIconDefinition {

	definition = append([]string{faicons.BaseFontAwesome}, definition...)
	return FontAwesomeIconDefinition{definition}

}

func NewSVGIcon(icon, color string) *HTMLPart {

	format := `<use xlink:href="%s">`
	return NewHTMLPart("svg", "", fmt.Sprintf(format, icon)).
		AddBootstrapClasses("icon").
		AddOption(&HTMLOption{"viewBox", "0 0 8 8"}).
		AddOption(&HTMLOption{"style", fmt.Sprintf("fill: %s;", color)})

}

func NewFontAwesomeIcon(icons ...FontAwesomeIconDefinition) *HTMLPart {

	result := NewHTMLPart("span", "", "").AddBootstrapClasses(faicons.ContainerStack, faicons.ModifyFixedWidth)
	for _, v := range icons {
		result.addSubPart(NewHTMLPart("i", "", "").AddBootstrapClasses(v.Classes...))
	}

	return result
}

//NewGlyphicon returns the HTMLPart needed to display the Icon
func NewGlyphicon(icon string) *HTMLPart {
	return NewHTMLPart("i", "", "").AddBootstrapClasses(bsglyphicons.Glyphicon, icon)

}

func (hp HTMLPart) String() string {
	return string(hp.bytes(false))
}

//String returns the HTML String for the HTMLPart struct includes all subparts subsubparts ...
func (hp HTMLPart) StringWithScripts() string {

	return string(hp.bytes(true))
}

func (hp HTMLPart) allScripts() *[]HTMLPart {

	scripts := []HTMLPart{}

	scripts = append(scripts, *hp.Scripts...)

	for _, v := range *hp.SubParts {
		scripts = append(scripts, *v.allScripts()...)
	}

	return &scripts

}

func (hp HTMLPart) bytes(withScripts bool) []byte {
	bb := make([]byte, 0, 1024)
	bb = append(bb, '<')
	bb = append(bb, hp.Class...)
	for _, v := range *hp.Options {
		bb = append(bb, ' ')
		bb = append(bb, v.Name...)
		bb = append(bb, '=', '"')
		bb = append(bb, v.Value...)
		bb = append(bb, '"')

	}
	bb = append(bb, '>')

	if hp.Class == "body" || hp.Class == "head" || withScripts {

		allscripts := *hp.allScripts()
		for _, v := range allscripts {
			bb = append(bb, v.bytes(false)...)
		}
	}
	bb = append(bb, hp.Content...)
	for _, v := range *hp.SubParts {
		bb = append(bb, v.bytes(false)...)
	}

	bb = append(bb, '<', '/')
	bb = append(bb, hp.Class...)
	bb = append(bb, '>')

	return bb
}

//JSONResultValue returns the string as js ready string "result.<myvalue>"
func JSONResultValue(myvalue string) string {
	return "result." + myvalue
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

//AddScripts adds one or more HTMLParts (subparts) in your HTMLPart
func (hp *HTMLPart) AddScripts(subparts ...*HTMLPart) *HTMLPart {
	for _, v := range subparts {
		*hp.Scripts = append(*hp.Scripts, *v)
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
