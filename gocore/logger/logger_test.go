package logger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"testing"
)

var regexPattern = `(\[[/.:\w]*\]) \[%s\] (\[[/.:\w]*\])\s`
var logMsgPattern = `.*%s.*`

func TestLogger(t *testing.T) {
	Init(Debug)
	w := bytes.NewBuffer([]byte{})

	log.SetOutput(w)

	LogDebug("testFlow", "%s", "#")

	b, err := ioutil.ReadAll(w)

	if err != nil {
		t.Log(err)
		t.Fatal()
	}
	if b == nil || len(b) == 0 {
		t.Log(string(b))
		t.Fail()
	}

	reg := regexp.MustCompile(fmt.Sprintf(regexPattern, "DEBUG"))

	matched := reg.MatchString(string(b))
	if !matched {
		t.Fail()
	}

	//Calls LogInfo and checks the output to be logged properly.

	w = bytes.NewBuffer([]byte{})
	log.SetOutput(w)
	LogInfo("testFlow", "%s", "#")

	b, err = ioutil.ReadAll(w)
	if err != nil {
		t.Log(err)
		t.Fatal()
	}
	if b == nil || len(b) == 0 {
		t.Log(string(b))
		t.Fail()
	}

	reg = regexp.MustCompile(fmt.Sprintf(regexPattern, "INFO"))

	matched = reg.MatchString(string(b))
	if !matched {
		t.Fail()
	}

	//Calls LogNotice and checks the output to be logged properly.

	w = bytes.NewBuffer([]byte{})
	log.SetOutput(w)
	LogNotice("testFlow", "%s", "#")

	b, err = ioutil.ReadAll(w)
	if err != nil {
		t.Log(err)
		t.Fatal()
	}
	if b == nil || len(b) == 0 {
		t.Log(string(b))
		t.Fail()
	}

	reg = regexp.MustCompile(fmt.Sprintf(regexPattern, "NOTICE"))

	matched = reg.MatchString(string(b))
	if !matched {
		t.Fail()
	}

	//Calls LogWarning and checks the output to be logged properly.

	w = bytes.NewBuffer([]byte{})
	log.SetOutput(w)
	LogWarning("testFlow", "%s", "#")

	b, err = ioutil.ReadAll(w)
	if err != nil {
		t.Log(err)
		t.Fatal()
	}
	if b == nil || len(b) == 0 {
		t.Log(string(b))
		t.Fail()
	}

	reg = regexp.MustCompile(fmt.Sprintf(regexPattern, "WARNING"))

	matched = reg.MatchString(string(b))
	if !matched {
		t.Fail()
	}

	//Calls LogInfo and checks the output to be logged properly.

	w = bytes.NewBuffer([]byte{})
	log.SetOutput(w)
	LogInfo("testFlow", "%s", "#")

	b, err = ioutil.ReadAll(w)
	if err != nil {
		t.Log(err)
		t.Fatal()
	}
	if b == nil || len(b) == 0 {
		t.Log(string(b))
		t.Fail()
	}
	reg = regexp.MustCompile(fmt.Sprintf(regexPattern, "INFO"))

	matched = reg.MatchString(string(b))
	if !matched {
		t.Fail()
	}

	//Calls LogError and checks the output to be logged properly.

	w = bytes.NewBuffer([]byte{})
	log.SetOutput(w)
	LogError("testFlow", "%s", "#")

	b, err = ioutil.ReadAll(w)
	if err != nil {
		t.Log(err)
		t.Fatal()
	}
	if b == nil || len(b) == 0 {
		t.Log(string(b))
		t.Fail()
	}
	reg = regexp.MustCompile(fmt.Sprintf(regexPattern, "ERROR"))

	matched = reg.MatchString(string(b))
	if !matched {
		t.Fail()
	}

	//Calls LogCritical and checks the output to be logged properly.

	w = bytes.NewBuffer([]byte{})
	log.SetOutput(w)
	LogCritical("testFlow", "%s", "#")

	b, err = ioutil.ReadAll(w)
	if err != nil {
		t.Log(err)
		t.Fatal()
	}
	if b == nil || len(b) == 0 {
		t.Log(string(b))
		t.Fail()
	}

	reg = regexp.MustCompile(fmt.Sprintf(regexPattern, "CRITICAL"))

	matched = reg.MatchString(string(b))
	if !matched {
		t.Fail()
	}

	//Calls LogEmergency and checks the output to be logged properly.

	w = bytes.NewBuffer([]byte{})
	log.SetOutput(w)
	LogEmergency("testFlow", "%s", "#")

	b, err = ioutil.ReadAll(w)
	if err != nil {
		t.Log(err)
		t.Fatal()
	}
	if b == nil || len(b) == 0 {
		t.Log(string(b))
		t.Fail()
	}

	reg = regexp.MustCompile(fmt.Sprintf(regexPattern, "EMERGENCY"))

	matched = reg.MatchString(string(b))
	if !matched {
		t.Log(reg.String())
		t.Log(string(b))
		t.Fail()
	}
	// add th test.v flag to restore the initial state

}
