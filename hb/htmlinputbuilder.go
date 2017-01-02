package hb

import (
	"fmt"

	"github.com/DaKine23/webapp/hb/bsbutton"
	"github.com/DaKine23/webapp/hb/bsglyphicons"
	"github.com/DaKine23/webapp/hb/bsinput"
)

//InputGroup holds the IDs and corresponding json field names of all input elements that fill a scripts json data
type InputGroup struct {
	IDs        []string
	ValueNames []string
}

func asJStoJSON(names, values []string, withQuotes bool) string {

	if len(names) != len(values) {
		return ""
	}
	result := `JSON.stringify({ `

	for i, v := range values {
		result += names[i] + `: `
		if withQuotes {
			result += "'"
		}
		result += v
		if withQuotes {
			result += "'"
		}
		result += `, `
	}

	return result + `})`
}

func getValueByID(ID string) string {
	format := `$("#%s").val()`
	return fmt.Sprintf(format, ID)
}

func setValueByID(ID, value string) string {
	format := `$("#%s").val("%s")`
	return fmt.Sprintf(format, ID, value)
}

func NewInputGroupScript() {

}

func newAjaxCall(source, action, restType, uri, data, onsuccess, onerror string) string {

	format := `$(document).ready(function(){
    $("#%s").%s(function(){
    $.ajax({
        type: '%s',
        url: '%s',
        dataType: 'json',
        data: {
            jsonData: %s
    },
    success: function (result) {
        %s
    },
    error: function () {
        %s
    }
    });
});`

	return fmt.Sprintf(format, source, action, restType, uri, data, onsuccess, onerror)
}

//OnSuccess generates the js function that puts the content in target if call was successfull
func OnSuccess(targetID, content string) string {

	return onResult(targetID, content)

}

//OnError generates the js function that puts the content in target if call was an error
func OnError(targetID, content string) string {
	return onResult(targetID, content)
}

func onResult(targetID, content string) string {
	format := `$("#%s").html(%s);`

	return fmt.Sprintf(format, targetID, content)
}

//Validation contains a regex for validation on client side
type Validation struct {
	RegEx *string
}

//NewLineEdit creates an HTMLPart containing a LineEdit
func NewLineEdit(ID, title, placeholder, content string, valid Validation) *HTMLPart {

	container := NewHTMLPart("div", ID+"inputgroup", "").AddBootstrapClasses(bsinput.InputGroup)
	if valid.RegEx != nil {

		format := `"use strict";
		$(document).ready(function(){
		$("#%s").%s(function(){
		    console.log("action is triggered");
		    var %sregex = new RegExp("%s");
		    if ($("#%s").val().match(%sregex) === null) {
		         %s
		        } else {
		         %s
		        }
		});
        });`
		setclasses := `$("#%s").attr("class", "%s");`
		validJs := fmt.Sprintf(setclasses, ID+"inputgroup", bsinput.InputGroup)
		notvalidJs := fmt.Sprintf(setclasses, ID+"inputgroup", bsinput.InputGroup+" "+bsinput.FormGroupHasFeedback+" "+bsinput.FormGroupHasWarning)

		script := NewHTMLPart("script", "", fmt.Sprintf(format, ID, "keyup", ID, *valid.RegEx, ID, ID, notvalidJs, validJs))

		script.AddOption(&HTMLOption{"type", "text/javascript"})
		container.AddScripts(script)
	}
	if len(title) > 0 {
		label := NewHTMLPart("span", "", title).AddBootstrapClasses(bsinput.InputGroupAddon)
		container.AddSubParts(label)
	}

	input := NewHTMLPart("input", ID, "").AddBootstrapClasses(bsinput.FormControl).
		AddOption(&HTMLOption{
			Name:  "type",
			Value: "text",
		}).
		AddOption(&HTMLOption{
			Name:  "name",
			Value: ID,
		}).
		AddOption(&HTMLOption{
			Name:  "placeholder",
			Value: placeholder,
		})
	if len(content) > 0 {
		input.AddOption(&HTMLOption{
			Name:  "value",
			Value: content,
		})
	}

	return container.AddSubParts(input)

}

type Range struct {
	Min *int
	Max *int
}

//AddLineEditSearch adds a search button in the end of a lineedit button as ID as
func (hp *HTMLPart) AddLineEditSearchButton(ID string) *HTMLPart {

	button := NewHTMLPart("div", "", "").AddBootstrapClasses(bsinput.InputGroupButton).addSubPart(
		NewHTMLPart("button", ID, "").
			AddBootstrapClasses(bsbutton.B, bsbutton.Default).
			AddOption(&HTMLOption{"type", "submit"}).
			addSubPart(NewGlyphicon(bsglyphicons.GlyphiconSearch)),
	)

	hp.addSubPart(button)

	return hp
}
