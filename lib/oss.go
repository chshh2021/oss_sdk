package lib

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	//ENDPOINT string = "oss-cn-beijing.aliyuncs.com"
	ENDPOINT   string = "oss-cn-beijing-internal.aliyuncs.com"
	GMT_LAYOUT string = "Mon, 02 Jan 2006 15:04:05 GMT"

	METHOD_GET    string = "GET"
	METHOD_PUT    string = "PUT"
	METHOD_DELETE string = "DELETE"
	METHOD_POST   string = "POST"
)

var (
	logger func(interface{})
)

type oss struct {
	access_id   string
	access_key  string
	urlString   string
	method      string
	bucket      string
	contentType string
	resource    string
	content     io.Reader
	req         *http.Request
	ossHeaders  map[string]string
}

func New(access_id, access_key string) *oss {
	return &oss{access_id: access_id, access_key: access_key, ossHeaders: make(map[string]string)}
}

func (o *oss) SetBucket(bucket string) {
	o.bucket = bucket
}

func (o *oss) SetMethod(method string) {
	o.method = method
}

func (o *oss) SetContentType(contentType string) {
	o.contentType = contentType
}

func (o *oss) SetResource(resource string) {
	o.resource = resource
}

func (o *oss) SetContent(content string) {
	o.content = strings.NewReader(content)
}

func (o *oss) SetOSSHeader(key, value string) {
	o.ossHeaders[key] = value
}

func (o *oss) PrepReq() {
	o.urlString = fmt.Sprintf("http://%s%s%s", ENDPOINT, o.bucket, o.resource)
	req, err := http.NewRequest(o.method, o.urlString, o.content)
	if err != nil {
		logger(err)
	}
	location, err := time.LoadLocation("GMT")
	now := time.Now().In(location).Format(GMT_LAYOUT)
	signature := generationSign(o.method, now, o.bucket, o.resource, o.access_key, o.contentType, o.ossHeaders)
	req.Header.Set("Date", now)
	req.Header.Set("Host", ENDPOINT)
	req.Header.Set("Authorization", fmt.Sprintf("OSS %s:%s", o.access_id, signature))
	req.Header.Set("Content-Type", o.contentType)

	for k, v := range o.ossHeaders {
		req.Header.Set(strings.ToUpper(k), v)
	}
	o.req = req
}

func (o *oss) Do() (int, http.Header, []byte) {
	client := &http.Client{}
	res, err := client.Do(o.req)
	if err != nil {
		logger(err)
		return http.StatusInternalServerError, make(http.Header), nil
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger(err)
	}

	if res.StatusCode != 200 {
		logger(body)
	}

	return res.StatusCode, res.Header, body
}

func generationSign(method, date, bucket, resource, access_key, content_type string, headers map[string]string) string {
	content_md5 := ""
	CanonicalizedOSSHeaders := ""
	for k, v := range headers {
		CanonicalizedOSSHeaders = fmt.Sprintf("%s%s:%s\n", CanonicalizedOSSHeaders, k, v)
	}
	CanonicalizedResource := fmt.Sprintf("%s%s", bucket, resource)

	sign := fmt.Sprintf("%s\n%s\n%s\n%s\n%s%s", method, content_md5, content_type, date, CanonicalizedOSSHeaders, CanonicalizedResource)

	hmacencoder := hmac.New(sha1.New, []byte(access_key))
	hmacencoder.Write([]byte(sign))

	signature := base64.StdEncoding.EncodeToString(hmacencoder.Sum(nil))
	return signature
}

var stdoss *oss

func Init(access_id, access_key string, mlogger func(interface{})) {
	stdoss = New(access_id, access_key)
	logger = mlogger
}

func MkDir(bucket, resource string) int {
	stdoss.SetOSSHeader("x-oss-acl", "public-read")
	stdoss.SetMethod(METHOD_PUT)
	stdoss.SetBucket(bucket)
	stdoss.SetContentType(MIMETYPE["default"])
	stdoss.SetContent("")
	stdoss.SetResource(fmt.Sprintf("%s/", resource))
	stdoss.PrepReq()
	status, _, _ := stdoss.Do()
	return status
}

func RmDir(bucket, resource string) int {
	stdoss.SetMethod(METHOD_DELETE)
	stdoss.SetBucket(bucket)
	stdoss.SetContentType(MIMETYPE["default"])
	stdoss.SetContent("")
	stdoss.SetResource(fmt.Sprintf("%s/", resource))
	stdoss.PrepReq()
	status, _, _ := stdoss.Do()
	return status
}

func Create(bucket, resource, content string) int {
	stdoss.SetOSSHeader("x-oss-acl", "public-read")
	stdoss.SetMethod(METHOD_PUT)
	stdoss.SetBucket(bucket)
	stdoss.SetContentType(MIMETYPE["jpg"])
	stdoss.SetContent(content)
	stdoss.SetResource(resource)
	stdoss.PrepReq()
	status, _, _ := stdoss.Do()
	return status
}

func CreateApk(bucket, resource, content string) int {
	stdoss.SetOSSHeader("x-oss-acl", "public-read")
	stdoss.SetMethod(METHOD_PUT)
	stdoss.SetBucket(bucket)
	stdoss.SetContentType(MIMETYPE["apk"])
	stdoss.SetContent(content)
	stdoss.SetResource(resource)
	stdoss.PrepReq()
	status, _, _ := stdoss.Do()
	return status
}

func Remove(bucket, resource string) int {
	stdoss.SetMethod(METHOD_DELETE)
	stdoss.SetBucket(bucket)
	stdoss.SetContentType(MIMETYPE["default"])
	stdoss.SetContent("")
	stdoss.SetResource(resource)
	stdoss.PrepReq()
	status, _, _ := stdoss.Do()
	return status
}

func CreateBucket(bucket string) int {
	stdoss.SetOSSHeader("x-oss-acl", "public-read")
	stdoss.SetMethod(METHOD_PUT)
	stdoss.SetBucket(bucket)
	stdoss.SetContentType(MIMETYPE["default"])
	stdoss.SetContent("")
	stdoss.SetResource("")
	stdoss.PrepReq()
	status, _, _ := stdoss.Do()
	return status
}

func RemoveBucket(bucket string) int {
	stdoss.SetMethod(METHOD_DELETE)
	stdoss.SetBucket(bucket)
	stdoss.SetContentType(MIMETYPE["default"])
	stdoss.SetContent("")
	stdoss.SetResource("")
	stdoss.PrepReq()
	status, _, _ := stdoss.Do()
	return status
}

func CreateMore(bucket, resource, content string, suffix string) int {
	stdoss.SetOSSHeader("x-oss-acl", "public-read")
	stdoss.SetMethod(METHOD_PUT)
	stdoss.SetBucket(bucket)
	stdoss.SetContentType(MIMETYPE[suffix])
	stdoss.SetContent(content)
	stdoss.SetResource(resource)
	stdoss.PrepReq()
	status, _, _ := stdoss.Do()
	return status
}

// 2025-03-18 Tiger
func GetSuffixByMimeType() map[string]string {
	typeInfo := make(map[string]string)
	for k, v := range MIMETYPE {
		typeInfo[v] = k
	}
	return typeInfo
}
