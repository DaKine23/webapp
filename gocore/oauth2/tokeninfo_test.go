package oauth2

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type retrieveScopes_args struct {
	res *http.Response
	err error
}
type retrieveScopes_case struct {
	name    string
	args    retrieveScopes_args
	want    *TokenInfo
	wantErr bool
}

const (
	retrieveScopes_sample              = `{"access_token":"eyJraWQiOiJwbGF0Zm9ybS1pYW0tdmNlaHloajYiLCJhbGciOiJFUzI1NiJ9.eyJzdWIiOiJzdHVwc19tZXRyaXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3JlYWxtIjoic2VydmljZXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3Rva2VuIjoiQmVhcmVyIiwiYXpwIjoic3R1cHNfbWV0cmlzXzA4MjYyMzRiLTZhNGMtNDYxNC04MTIwLWRlNTMzM2RiZTJhOCIsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vYnAiOiI4MTBkMWQwMC00MzEyLTQzZTUtYmQzMS1kODM3M2ZkZDI0YzciLCJpc3MiOiJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tIiwiZXhwIjoxNTA1MjAxNzUxLCJpYXQiOjE1MDUxOTgxNDEsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vcHJpdmlsZWdlcyI6WyJjb20uemFsYW5kbzo6c3BwLWJyYW5kLXNlcnZpY2UuYnJhbmRzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLXByb2R1Y3RzLnByb2R1Y3RzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1lZGlhLm1lZGlhLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1hc3Rlci1hdHRyaWJ1dGVzLmF0dHJpYnV0ZXMucmVhZCIsImNvbS56YWxhbmRvOjpjYXRhbG9nX3NlcnZpY2UucmVhZF9hbGwiLCJjb20uemFsYW5kbzo6dHJhbnNsYXRpb24ud3JpdGVfYWxsIl19._X4zlW2n8EHuXz3RnLa4k-sguVoVKPlbTymnBODYqRQnwPSpVJHfy6nWHJaOaS_HaCamFYdPH_d8TYruRUNEtg","catalog_service.read_all":true,"client_id":"stups_metris_0826234b-6a4c-4614-8120-de5333dbe2a8","expires_in":3575,"grant_type":"password","realm":"/services","scope":["spp-brand-service.brands.read","spp-products.products.read","spp-media.media.read","spp-master-attributes.attributes.read","catalog_service.read_all","translation.write_all","uid"],"spp-brand-service.brands.read":true,"spp-master-attributes.attributes.read":true,"spp-media.media.read":true,"spp-products.products.read":true,"token_type":"Bearer","translation.write_all":true,"uid":"stups_metris"}`
	retrieveScopes_just_expired_sample = `{"access_token":"ayJraWQiOiJwbGF0Zm9ybS1pYW0tdmNlaHloajYiLCJhbGciOiJFUzI1NiJ9.eyJzdWIiOiJzdHVwc19tZXRyaXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3JlYWxtIjoic2VydmljZXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3Rva2VuIjoiQmVhcmVyIiwiYXpwIjoic3R1cHNfbWV0cmlzXzA4MjYyMzRiLTZhNGMtNDYxNC04MTIwLWRlNTMzM2RiZTJhOCIsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vYnAiOiI4MTBkMWQwMC00MzEyLTQzZTUtYmQzMS1kODM3M2ZkZDI0YzciLCJpc3MiOiJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tIiwiZXhwIjoxNTA1MjAxNzUxLCJpYXQiOjE1MDUxOTgxNDEsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vcHJpdmlsZWdlcyI6WyJjb20uemFsYW5kbzo6c3BwLWJyYW5kLXNlcnZpY2UuYnJhbmRzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLXByb2R1Y3RzLnByb2R1Y3RzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1lZGlhLm1lZGlhLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1hc3Rlci1hdHRyaWJ1dGVzLmF0dHJpYnV0ZXMucmVhZCIsImNvbS56YWxhbmRvOjpjYXRhbG9nX3NlcnZpY2UucmVhZF9hbGwiLCJjb20uemFsYW5kbzo6dHJhbnNsYXRpb24ud3JpdGVfYWxsIl19._X4zlW2n8EHuXz3RnLa4k-sguVoVKPlbTymnBODYqRQnwPSpVJHfy6nWHJaOaS_HaCamFYdPH_d8TYruRUNEtg","catalog_service.read_all":true,"client_id":"stups_metris_0826234b-6a4c-4614-8120-de5333dbe2a8","expires_in":0,"grant_type":"password","realm":"/services","scope":["spp-brand-service.brands.read","spp-products.products.read","spp-media.media.read","spp-master-attributes.attributes.read","catalog_service.read_all","translation.write_all","uid"],"spp-brand-service.brands.read":true,"spp-master-attributes.attributes.read":true,"spp-media.media.read":true,"spp-products.products.read":true,"token_type":"Bearer","translation.write_all":true,"uid":"stups_metris"}`
)

func Test_retrieveScopes(t *testing.T) {

	tests := []retrieveScopes_case{
		retrieveScopes_case{
			name:    "happy case",
			want:    &TokenInfo{"eyJraWQiOiJwbGF0Zm9ybS1pYW0tdmNlaHloajYiLCJhbGciOiJFUzI1NiJ9.eyJzdWIiOiJzdHVwc19tZXRyaXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3JlYWxtIjoic2VydmljZXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3Rva2VuIjoiQmVhcmVyIiwiYXpwIjoic3R1cHNfbWV0cmlzXzA4MjYyMzRiLTZhNGMtNDYxNC04MTIwLWRlNTMzM2RiZTJhOCIsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vYnAiOiI4MTBkMWQwMC00MzEyLTQzZTUtYmQzMS1kODM3M2ZkZDI0YzciLCJpc3MiOiJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tIiwiZXhwIjoxNTA1MjAxNzUxLCJpYXQiOjE1MDUxOTgxNDEsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vcHJpdmlsZWdlcyI6WyJjb20uemFsYW5kbzo6c3BwLWJyYW5kLXNlcnZpY2UuYnJhbmRzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLXByb2R1Y3RzLnByb2R1Y3RzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1lZGlhLm1lZGlhLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1hc3Rlci1hdHRyaWJ1dGVzLmF0dHJpYnV0ZXMucmVhZCIsImNvbS56YWxhbmRvOjpjYXRhbG9nX3NlcnZpY2UucmVhZF9hbGwiLCJjb20uemFsYW5kbzo6dHJhbnNsYXRpb24ud3JpdGVfYWxsIl19._X4zlW2n8EHuXz3RnLa4k-sguVoVKPlbTymnBODYqRQnwPSpVJHfy6nWHJaOaS_HaCamFYdPH_d8TYruRUNEtg", "stups_metris_0826234b-6a4c-4614-8120-de5333dbe2a8", 3575, "password", "/services", []string{"spp-brand-service.brands.read", "spp-products.products.read", "spp-media.media.read", "spp-master-attributes.attributes.read", "catalog_service.read_all", "translation.write_all", "uid"}, "Bearer", "stups_metris"},
			wantErr: false,
			args: retrieveScopes_args{
				res: &http.Response{
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(retrieveScopes_sample))),
					StatusCode: http.StatusOK,
				},
				err: nil,
			},
		},
		retrieveScopes_case{
			name:    "error before response",
			want:    nil,
			wantErr: true,
			args: retrieveScopes_args{
				res: nil,
				err: errors.New("random error"),
			},
		},
		retrieveScopes_case{
			name:    "statuscode != 2XX",
			want:    nil,
			wantErr: true,
			args: retrieveScopes_args{
				res: &http.Response{
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(retrieveScopes_sample))),
					StatusCode: http.StatusNotFound,
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := retrieveScopes(tt.args.res, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("retrieveScopes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("retrieveScopes() = %v, want %v", got, tt.want)
			}
		})
	}
}

type RetrieveTokenInfo_args struct {
	tokenInfoService string
	token            string
	client           Client
}

type RetrieveTokenInfo_case struct {
	name    string
	args    RetrieveTokenInfo_args
	want    *TokenInfo
	wantErr bool
}

type retrieveToken_client struct {
	resp *http.Response
	err  error
}

func (rtc retrieveToken_client) Do(req *http.Request) (*http.Response, error) {
	return rtc.resp, rtc.err
}

func TestRetrieveTokenInfo(t *testing.T) {

	errorClient := retrieveToken_client{nil, errors.New("random error")}
	encodingErrorClient := retrieveToken_client{&http.Response{
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(`1243"!§$!412412`))),
		StatusCode: http.StatusOK,
	}, nil}
	happyClient := retrieveToken_client{&http.Response{
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(retrieveScopes_sample))),
		StatusCode: http.StatusOK,
	}, nil}

	tests := []RetrieveTokenInfo_case{
		RetrieveTokenInfo_case{
			name:    "happy case",
			want:    &TokenInfo{"eyJraWQiOiJwbGF0Zm9ybS1pYW0tdmNlaHloajYiLCJhbGciOiJFUzI1NiJ9.eyJzdWIiOiJzdHVwc19tZXRyaXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3JlYWxtIjoic2VydmljZXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3Rva2VuIjoiQmVhcmVyIiwiYXpwIjoic3R1cHNfbWV0cmlzXzA4MjYyMzRiLTZhNGMtNDYxNC04MTIwLWRlNTMzM2RiZTJhOCIsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vYnAiOiI4MTBkMWQwMC00MzEyLTQzZTUtYmQzMS1kODM3M2ZkZDI0YzciLCJpc3MiOiJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tIiwiZXhwIjoxNTA1MjAxNzUxLCJpYXQiOjE1MDUxOTgxNDEsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vcHJpdmlsZWdlcyI6WyJjb20uemFsYW5kbzo6c3BwLWJyYW5kLXNlcnZpY2UuYnJhbmRzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLXByb2R1Y3RzLnByb2R1Y3RzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1lZGlhLm1lZGlhLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1hc3Rlci1hdHRyaWJ1dGVzLmF0dHJpYnV0ZXMucmVhZCIsImNvbS56YWxhbmRvOjpjYXRhbG9nX3NlcnZpY2UucmVhZF9hbGwiLCJjb20uemFsYW5kbzo6dHJhbnNsYXRpb24ud3JpdGVfYWxsIl19._X4zlW2n8EHuXz3RnLa4k-sguVoVKPlbTymnBODYqRQnwPSpVJHfy6nWHJaOaS_HaCamFYdPH_d8TYruRUNEtg", "stups_metris_0826234b-6a4c-4614-8120-de5333dbe2a8", 3575, "password", "/services", []string{"spp-brand-service.brands.read", "spp-products.products.read", "spp-media.media.read", "spp-master-attributes.attributes.read", "catalog_service.read_all", "translation.write_all", "uid"}, "Bearer", "stups_metris"},
			wantErr: false,
			args: RetrieveTokenInfo_args{
				token:            "not empty", //happyclient has a fixed answer
				tokenInfoService: "not.empty.host.com",
				client:           happyClient,
			},
		},
		RetrieveTokenInfo_case{
			name:    "error",
			want:    nil,
			wantErr: true,
			args: RetrieveTokenInfo_args{
				token:            "not empty",
				tokenInfoService: "not.empty.host.com",
				client:           errorClient,
			},
		},
		RetrieveTokenInfo_case{
			name:    "empty param1",
			want:    nil,
			wantErr: true,
			args: RetrieveTokenInfo_args{
				token:            "",
				tokenInfoService: "not.empty.host.com",
				client:           happyClient,
			},
		},
		RetrieveTokenInfo_case{
			name:    "empty param2",
			want:    nil,
			wantErr: true,
			args: RetrieveTokenInfo_args{
				token:            "not empty",
				tokenInfoService: "",
				client:           happyClient,
			},
		},
		RetrieveTokenInfo_case{
			name:    "empty param3",
			want:    nil,
			wantErr: true,
			args: RetrieveTokenInfo_args{
				token:            "",
				tokenInfoService: "",
				client:           happyClient,
			},
		},
		RetrieveTokenInfo_case{
			name:    "encodingError",
			want:    nil,
			wantErr: true,
			args: RetrieveTokenInfo_args{
				token:            "not empty",
				tokenInfoService: "not.empty.host.com",
				client:           encodingErrorClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.args.client)
			got, err := RetrieveTokenInfo(tt.args.tokenInfoService, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveTokenInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveTokenInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetrieveTokenInfoCache(t *testing.T) {
	//cache test

	errorClient := retrieveToken_client{nil, errors.New("random error")}
	happyClient := retrieveToken_client{&http.Response{
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(retrieveScopes_sample))),
		StatusCode: http.StatusOK,
	}, nil}

	preCacheCase := RetrieveTokenInfo_case{
		name:    "caching",
		want:    &TokenInfo{"eyJraWQiOiJwbGF0Zm9ybS1pYW0tdmNlaHloajYiLCJhbGciOiJFUzI1NiJ9.eyJzdWIiOiJzdHVwc19tZXRyaXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3JlYWxtIjoic2VydmljZXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3Rva2VuIjoiQmVhcmVyIiwiYXpwIjoic3R1cHNfbWV0cmlzXzA4MjYyMzRiLTZhNGMtNDYxNC04MTIwLWRlNTMzM2RiZTJhOCIsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vYnAiOiI4MTBkMWQwMC00MzEyLTQzZTUtYmQzMS1kODM3M2ZkZDI0YzciLCJpc3MiOiJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tIiwiZXhwIjoxNTA1MjAxNzUxLCJpYXQiOjE1MDUxOTgxNDEsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vcHJpdmlsZWdlcyI6WyJjb20uemFsYW5kbzo6c3BwLWJyYW5kLXNlcnZpY2UuYnJhbmRzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLXByb2R1Y3RzLnByb2R1Y3RzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1lZGlhLm1lZGlhLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1hc3Rlci1hdHRyaWJ1dGVzLmF0dHJpYnV0ZXMucmVhZCIsImNvbS56YWxhbmRvOjpjYXRhbG9nX3NlcnZpY2UucmVhZF9hbGwiLCJjb20uemFsYW5kbzo6dHJhbnNsYXRpb24ud3JpdGVfYWxsIl19._X4zlW2n8EHuXz3RnLa4k-sguVoVKPlbTymnBODYqRQnwPSpVJHfy6nWHJaOaS_HaCamFYdPH_d8TYruRUNEtg", "stups_metris_0826234b-6a4c-4614-8120-de5333dbe2a8", 3575, "password", "/services", []string{"spp-brand-service.brands.read", "spp-products.products.read", "spp-media.media.read", "spp-master-attributes.attributes.read", "catalog_service.read_all", "translation.write_all", "uid"}, "Bearer", "stups_metris"},
		wantErr: false,
		args: RetrieveTokenInfo_args{
			token:            "eyJraWQiOiJwbGF0Zm9ybS1pYW0tdmNlaHloajYiLCJhbGciOiJFUzI1NiJ9.eyJzdWIiOiJzdHVwc19tZXRyaXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3JlYWxtIjoic2VydmljZXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3Rva2VuIjoiQmVhcmVyIiwiYXpwIjoic3R1cHNfbWV0cmlzXzA4MjYyMzRiLTZhNGMtNDYxNC04MTIwLWRlNTMzM2RiZTJhOCIsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vYnAiOiI4MTBkMWQwMC00MzEyLTQzZTUtYmQzMS1kODM3M2ZkZDI0YzciLCJpc3MiOiJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tIiwiZXhwIjoxNTA1MjAxNzUxLCJpYXQiOjE1MDUxOTgxNDEsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vcHJpdmlsZWdlcyI6WyJjb20uemFsYW5kbzo6c3BwLWJyYW5kLXNlcnZpY2UuYnJhbmRzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLXByb2R1Y3RzLnByb2R1Y3RzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1lZGlhLm1lZGlhLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1hc3Rlci1hdHRyaWJ1dGVzLmF0dHJpYnV0ZXMucmVhZCIsImNvbS56YWxhbmRvOjpjYXRhbG9nX3NlcnZpY2UucmVhZF9hbGwiLCJjb20uemFsYW5kbzo6dHJhbnNsYXRpb24ud3JpdGVfYWxsIl19._X4zlW2n8EHuXz3RnLa4k-sguVoVKPlbTymnBODYqRQnwPSpVJHfy6nWHJaOaS_HaCamFYdPH_d8TYruRUNEtg",
			tokenInfoService: "not.empty.host.com",
			client:           happyClient,
		},
	}

	postCacheCase := RetrieveTokenInfo_case{
		name:    "caching",
		want:    &TokenInfo{"eyJraWQiOiJwbGF0Zm9ybS1pYW0tdmNlaHloajYiLCJhbGciOiJFUzI1NiJ9.eyJzdWIiOiJzdHVwc19tZXRyaXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3JlYWxtIjoic2VydmljZXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3Rva2VuIjoiQmVhcmVyIiwiYXpwIjoic3R1cHNfbWV0cmlzXzA4MjYyMzRiLTZhNGMtNDYxNC04MTIwLWRlNTMzM2RiZTJhOCIsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vYnAiOiI4MTBkMWQwMC00MzEyLTQzZTUtYmQzMS1kODM3M2ZkZDI0YzciLCJpc3MiOiJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tIiwiZXhwIjoxNTA1MjAxNzUxLCJpYXQiOjE1MDUxOTgxNDEsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vcHJpdmlsZWdlcyI6WyJjb20uemFsYW5kbzo6c3BwLWJyYW5kLXNlcnZpY2UuYnJhbmRzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLXByb2R1Y3RzLnByb2R1Y3RzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1lZGlhLm1lZGlhLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1hc3Rlci1hdHRyaWJ1dGVzLmF0dHJpYnV0ZXMucmVhZCIsImNvbS56YWxhbmRvOjpjYXRhbG9nX3NlcnZpY2UucmVhZF9hbGwiLCJjb20uemFsYW5kbzo6dHJhbnNsYXRpb24ud3JpdGVfYWxsIl19._X4zlW2n8EHuXz3RnLa4k-sguVoVKPlbTymnBODYqRQnwPSpVJHfy6nWHJaOaS_HaCamFYdPH_d8TYruRUNEtg", "stups_metris_0826234b-6a4c-4614-8120-de5333dbe2a8", 3575, "password", "/services", []string{"spp-brand-service.brands.read", "spp-products.products.read", "spp-media.media.read", "spp-master-attributes.attributes.read", "catalog_service.read_all", "translation.write_all", "uid"}, "Bearer", "stups_metris"},
		wantErr: false,
		args: RetrieveTokenInfo_args{
			token:            "eyJraWQiOiJwbGF0Zm9ybS1pYW0tdmNlaHloajYiLCJhbGciOiJFUzI1NiJ9.eyJzdWIiOiJzdHVwc19tZXRyaXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3JlYWxtIjoic2VydmljZXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3Rva2VuIjoiQmVhcmVyIiwiYXpwIjoic3R1cHNfbWV0cmlzXzA4MjYyMzRiLTZhNGMtNDYxNC04MTIwLWRlNTMzM2RiZTJhOCIsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vYnAiOiI4MTBkMWQwMC00MzEyLTQzZTUtYmQzMS1kODM3M2ZkZDI0YzciLCJpc3MiOiJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tIiwiZXhwIjoxNTA1MjAxNzUxLCJpYXQiOjE1MDUxOTgxNDEsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vcHJpdmlsZWdlcyI6WyJjb20uemFsYW5kbzo6c3BwLWJyYW5kLXNlcnZpY2UuYnJhbmRzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLXByb2R1Y3RzLnByb2R1Y3RzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1lZGlhLm1lZGlhLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1hc3Rlci1hdHRyaWJ1dGVzLmF0dHJpYnV0ZXMucmVhZCIsImNvbS56YWxhbmRvOjpjYXRhbG9nX3NlcnZpY2UucmVhZF9hbGwiLCJjb20uemFsYW5kbzo6dHJhbnNsYXRpb24ud3JpdGVfYWxsIl19._X4zlW2n8EHuXz3RnLa4k-sguVoVKPlbTymnBODYqRQnwPSpVJHfy6nWHJaOaS_HaCamFYdPH_d8TYruRUNEtg",
			tokenInfoService: "not.empty.host.com",
			client:           errorClient,
		},
	}

	t.Run(preCacheCase.name, func(t *testing.T) {
		Init(preCacheCase.args.client)
		got, err := RetrieveTokenInfo(preCacheCase.args.tokenInfoService, preCacheCase.args.token)
		if (err != nil) != preCacheCase.wantErr {
			t.Errorf("RetrieveTokenInfo() error = %v, wantErr %v", err, preCacheCase.wantErr)
			return
		}
		if !reflect.DeepEqual(got, preCacheCase.want) {
			t.Errorf("RetrieveTokenInfo() = %v, want %v", got, preCacheCase.want)
		}
		Init(postCacheCase.args.client)
		got, err = RetrieveTokenInfo(postCacheCase.args.tokenInfoService, postCacheCase.args.token)
		if (err != nil) != postCacheCase.wantErr {
			t.Errorf("RetrieveTokenInfo() error = %v, wantErr %v", err, postCacheCase.wantErr)
			return
		}
		if !reflect.DeepEqual(got, postCacheCase.want) {
			t.Errorf("RetrieveTokenInfo() = %v, want %v", got, postCacheCase.want)
		}

	})

}

func TestRetrieveTokenInfoCacheExpired(t *testing.T) {
	//cache test

	errorClient := retrieveToken_client{nil, errors.New("random error")}
	happyClientWithJustExpiredToken := retrieveToken_client{&http.Response{
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(retrieveScopes_just_expired_sample))),
		StatusCode: http.StatusOK,
	}, nil}

	preCacheCase := RetrieveTokenInfo_case{
		name:    "caching expired",
		want:    &TokenInfo{"ayJraWQiOiJwbGF0Zm9ybS1pYW0tdmNlaHloajYiLCJhbGciOiJFUzI1NiJ9.eyJzdWIiOiJzdHVwc19tZXRyaXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3JlYWxtIjoic2VydmljZXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3Rva2VuIjoiQmVhcmVyIiwiYXpwIjoic3R1cHNfbWV0cmlzXzA4MjYyMzRiLTZhNGMtNDYxNC04MTIwLWRlNTMzM2RiZTJhOCIsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vYnAiOiI4MTBkMWQwMC00MzEyLTQzZTUtYmQzMS1kODM3M2ZkZDI0YzciLCJpc3MiOiJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tIiwiZXhwIjoxNTA1MjAxNzUxLCJpYXQiOjE1MDUxOTgxNDEsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vcHJpdmlsZWdlcyI6WyJjb20uemFsYW5kbzo6c3BwLWJyYW5kLXNlcnZpY2UuYnJhbmRzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLXByb2R1Y3RzLnByb2R1Y3RzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1lZGlhLm1lZGlhLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1hc3Rlci1hdHRyaWJ1dGVzLmF0dHJpYnV0ZXMucmVhZCIsImNvbS56YWxhbmRvOjpjYXRhbG9nX3NlcnZpY2UucmVhZF9hbGwiLCJjb20uemFsYW5kbzo6dHJhbnNsYXRpb24ud3JpdGVfYWxsIl19._X4zlW2n8EHuXz3RnLa4k-sguVoVKPlbTymnBODYqRQnwPSpVJHfy6nWHJaOaS_HaCamFYdPH_d8TYruRUNEtg", "stups_metris_0826234b-6a4c-4614-8120-de5333dbe2a8", 0, "password", "/services", []string{"spp-brand-service.brands.read", "spp-products.products.read", "spp-media.media.read", "spp-master-attributes.attributes.read", "catalog_service.read_all", "translation.write_all", "uid"}, "Bearer", "stups_metris"},
		wantErr: false,
		args: RetrieveTokenInfo_args{
			token:            "ayJraWQiOiJwbGF0Zm9ybS1pYW0tdmNlaHloajYiLCJhbGciOiJFUzI1NiJ9.eyJzdWIiOiJzdHVwc19tZXRyaXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3JlYWxtIjoic2VydmljZXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3Rva2VuIjoiQmVhcmVyIiwiYXpwIjoic3R1cHNfbWV0cmlzXzA4MjYyMzRiLTZhNGMtNDYxNC04MTIwLWRlNTMzM2RiZTJhOCIsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vYnAiOiI4MTBkMWQwMC00MzEyLTQzZTUtYmQzMS1kODM3M2ZkZDI0YzciLCJpc3MiOiJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tIiwiZXhwIjoxNTA1MjAxNzUxLCJpYXQiOjE1MDUxOTgxNDEsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vcHJpdmlsZWdlcyI6WyJjb20uemFsYW5kbzo6c3BwLWJyYW5kLXNlcnZpY2UuYnJhbmRzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLXByb2R1Y3RzLnByb2R1Y3RzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1lZGlhLm1lZGlhLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1hc3Rlci1hdHRyaWJ1dGVzLmF0dHJpYnV0ZXMucmVhZCIsImNvbS56YWxhbmRvOjpjYXRhbG9nX3NlcnZpY2UucmVhZF9hbGwiLCJjb20uemFsYW5kbzo6dHJhbnNsYXRpb24ud3JpdGVfYWxsIl19._X4zlW2n8EHuXz3RnLa4k-sguVoVKPlbTymnBODYqRQnwPSpVJHfy6nWHJaOaS_HaCamFYdPH_d8TYruRUNEtg",
			tokenInfoService: "not.empty.host.com",
			client:           happyClientWithJustExpiredToken,
		},
	}

	postCacheCase := RetrieveTokenInfo_case{
		name:    "caching",
		want:    nil,
		wantErr: true,
		args: RetrieveTokenInfo_args{
			token:            "ayJraWQiOiJwbGF0Zm9ybS1pYW0tdmNlaHloajYiLCJhbGciOiJFUzI1NiJ9.eyJzdWIiOiJzdHVwc19tZXRyaXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3JlYWxtIjoic2VydmljZXMiLCJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tL3Rva2VuIjoiQmVhcmVyIiwiYXpwIjoic3R1cHNfbWV0cmlzXzA4MjYyMzRiLTZhNGMtNDYxNC04MTIwLWRlNTMzM2RiZTJhOCIsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vYnAiOiI4MTBkMWQwMC00MzEyLTQzZTUtYmQzMS1kODM3M2ZkZDI0YzciLCJpc3MiOiJodHRwczovL2lkZW50aXR5LnphbGFuZG8uY29tIiwiZXhwIjoxNTA1MjAxNzUxLCJpYXQiOjE1MDUxOTgxNDEsImh0dHBzOi8vaWRlbnRpdHkuemFsYW5kby5jb20vcHJpdmlsZWdlcyI6WyJjb20uemFsYW5kbzo6c3BwLWJyYW5kLXNlcnZpY2UuYnJhbmRzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLXByb2R1Y3RzLnByb2R1Y3RzLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1lZGlhLm1lZGlhLnJlYWQiLCJjb20uemFsYW5kbzo6c3BwLW1hc3Rlci1hdHRyaWJ1dGVzLmF0dHJpYnV0ZXMucmVhZCIsImNvbS56YWxhbmRvOjpjYXRhbG9nX3NlcnZpY2UucmVhZF9hbGwiLCJjb20uemFsYW5kbzo6dHJhbnNsYXRpb24ud3JpdGVfYWxsIl19._X4zlW2n8EHuXz3RnLa4k-sguVoVKPlbTymnBODYqRQnwPSpVJHfy6nWHJaOaS_HaCamFYdPH_d8TYruRUNEtg",
			tokenInfoService: "not.empty.host.com",
			client:           errorClient,
		},
	}

	t.Run(preCacheCase.name, func(t *testing.T) {
		Init(preCacheCase.args.client)
		got, err := RetrieveTokenInfo(preCacheCase.args.tokenInfoService, preCacheCase.args.token)
		if (err != nil) != preCacheCase.wantErr {
			t.Errorf("RetrieveTokenInfo() error = %v, wantErr %v", err, preCacheCase.wantErr)
			return
		}
		if !reflect.DeepEqual(got, preCacheCase.want) {
			t.Errorf("RetrieveTokenInfo() = %v, want %v", got, preCacheCase.want)
		}
		time.Sleep(time.Millisecond * 50) // wait for cache to die

		Init(postCacheCase.args.client)
		got, err = RetrieveTokenInfo(postCacheCase.args.tokenInfoService, postCacheCase.args.token)
		if (err != nil) != postCacheCase.wantErr {
			t.Errorf("RetrieveTokenInfo() error = %v, wantErr %v", err, postCacheCase.wantErr)
			return
		}
		if !reflect.DeepEqual(got, postCacheCase.want) {
			t.Errorf("RetrieveTokenInfo() = %v, want %v", got, postCacheCase.want)
		}

	})

}
