package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.bus.zalan.do/ale/gocore/oauth2"
	"github.com/gin-gonic/gin"
)

const (
	readScope  = "veritas.article_measurement.read"
	writeScope = "veritas.article_measurement.write"
)

func TestMiddleWare(t *testing.T) {

	w := bytes.NewBuffer([]byte{})

	log.SetOutput(w)

	gin.SetMode(gin.ReleaseMode)

	ping := func(c *gin.Context) {

		if c.Keys[FlowIDKey] != "someKey" {
			c.String(500, "NO!")
		}

		c.String(200, "Yeah!")
	}

	go func() {
		ginr := gin.New()
		ginr.Use(gin.Recovery(), FlowID(), Logger())
		ginr.GET("/ping", ping)
		ginr.Run(":8080")
	}()

	var resp []byte

	client := http.DefaultClient
	client.Timeout = time.Second * 1
	time.Sleep(time.Second * 1)
	code := 0

	for len(resp) == 0 {

		req, err := http.NewRequest("GET", "http://127.0.0.1:8080/ping", bytes.NewBufferString("Yeah!"))
		if err != nil {
			fmt.Println(err)
		}

		req.Header.Add(FlowIDHeaderKey, "someKey")

		res, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Millisecond * 200)
		if res != nil {
			resp, _ = ioutil.ReadAll(res.Body)
			code = res.StatusCode
		}

	}

	if code != 200 {

		t.Fail()
	}
	b, _ := ioutil.ReadAll(w)

	reg := regexp.MustCompile(`^(\[[/.:\w]+\]\s){3,3}`)

	if !reg.MatchString(string(b)) {
		t.Log("prefix", string(b))
		t.Fail()
	}

	t.Fail()
}

type extractTokenCase struct {
	name      string
	args      extractTokensArgs
	wantToken string
}
type extractTokensArgs struct {
	authorization string
}

func Test_extractToken(t *testing.T) {

	tests := []extractTokenCase{
		extractTokenCase{
			name:      "happy case",
			wantToken: "1234",
			args: extractTokensArgs{
				authorization: "Bearer 1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotToken := extractToken(tt.args.authorization); gotToken != tt.wantToken {
				t.Errorf("extractToken() = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

type checkScopeListArgs struct {
	scope string
	ti    *oauth2.TokenInfo
}

type checkScopeListCase struct {
	name    string
	args    checkScopeListArgs
	wantErr bool
}

func Test_checkScopeList(t *testing.T) {

	tests := []checkScopeListCase{
		checkScopeListCase{
			name:    "happy case",
			wantErr: false,
			args: checkScopeListArgs{
				scope: readScope,
				ti: &oauth2.TokenInfo{
					Scope: []string{readScope, writeScope},
				},
			},
		},
		checkScopeListCase{
			name:    "error case",
			wantErr: true,
			args: checkScopeListArgs{
				scope: readScope,
				ti: &oauth2.TokenInfo{
					Scope: []string{writeScope},
				},
			},
		},
		checkScopeListCase{
			name:    "error case",
			wantErr: true,
			args: checkScopeListArgs{
				scope: writeScope,
				ti: &oauth2.TokenInfo{
					Scope: []string{},
				},
			},
		},
		checkScopeListCase{
			name:    "empty so error case",
			wantErr: true,
			args: checkScopeListArgs{
				scope: "",
				ti: &oauth2.TokenInfo{
					Scope: []string{readScope, writeScope},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkScopeList(tt.args.scope, tt.args.ti); (err != nil) != tt.wantErr {
				t.Errorf("checkScopeList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
