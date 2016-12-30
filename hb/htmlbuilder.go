package hb

import (
	"fmt"
	"strings"

	"github.com/DaKine23/webapp/hb/bsglyphicons"
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

//NewGlyphicon returns the HTMLPart needed to display the Icon
func NewGlyphicon(icon string) *HTMLPart {
	return NewHTMLPart("i", "", "").AddBootstrapClasses(bsglyphicons.Glyphicon, icon)

}

//String returns the HTML String for the HTMLPart struct includes all subparts subsubparts ...
func (hp HTMLPart) String() string {

	return string(hp.bytes())

}

func (hp HTMLPart) bytes() []byte {
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
	bb = append(bb, hp.Content...)

	for _, v := range *hp.SubParts {
		bb = append(bb, v.bytes()...)
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
