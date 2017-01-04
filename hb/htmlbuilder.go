package hb

import (
	"fmt"
	"strings"

	"github.com/DaKine23/webapp/hb/bscontainer"
	"github.com/DaKine23/webapp/hb/bsglyphicons"
	"github.com/DaKine23/webapp/hb/bsgrid"
	"github.com/DaKine23/webapp/hb/faicons"
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

func (bsg BsGrid) HTMLPart() *HTMLPart {

	container := NewHTMLPart("div", bsg.ID).AddBootstrapClasses(bscontainer.ContainerFluid)

	for _, v := range *bsg.Grid {
		row := NewHTMLPart(bsgrid.Row, "").AddBootstrapClasses(bsgrid.Row)
		autocolspan := 12
		for _, v2 := range v {
			autocolspan -= v2.Colspan
		}
		autocolspan = autocolspan / len(v)
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

	// root := hb.NewHTMLPart("root", "", "").AddBootstrapClasses(bscontainer.Container)
	// 	row1 := hb.NewHTMLPart("row", "", "").AddBootstrapClasses(bsgrid.Row)
	// 	cell11 := hb.NewHTMLPart("cell", "", "").AddBootstrapClasses(bsgrid.Cell(12, bsgrid.Large))
	// 	cell11.AddSubParts(buttongroup)
	// 	row1.AddSubParts(cell11)

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

//Script represents an HTMLPart containing a Script (mostly to do an Ajax call)
type Script struct {
	*HTMLPart
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
func NewScript(source, action, target, restType, apicall string, newContent string) *Script {

	return script(source, action, target, restType, apicall+`"`, newContent)
}

func script(source, action, target, restType, apicall string, newContent string) *Script {

	part := NewHTMLPart("script", "", fmt.Sprintf(`$(document).ready(function(){
    $("#%s").%s(function(){
        $.ajax({type: "%s", url: "%s, async: true, success: function(result){
            $("#%s").html(%s);
        }});
    });
    });`, source, action, restType, apicall, target, newContent))

	script := Script{
		HTMLPart: part,
	}
	return &script
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

//String returns the HTML String for the HTMLPart struct includes all subparts subsubparts ...
func (hp HTMLPart) String(withScripts bool) string {

	return string(hp.bytes(withScripts))

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
