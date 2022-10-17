package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	oss "github.com/kyf/oss_sdk/lib"
)

func parseFileMime(base64Str string) (mtype string) {
	pos := strings.Index(base64Str, ";")
	prefix := base64Str[0:pos]
	mtype = strings.Replace(prefix, "data:", "", -1)
	return mtype
}

func voucherHandler(group string, mimeMap map[string]string) MyHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		data := r.Form
		file_content := data.Get("file")
		var result map[string]string
		var jsonResult []byte
		if strings.EqualFold(file_content, "") {
			result = map[string]string{
				"status": "error",
				"msg":    "file is empty",
			}

			jsonResult, _ = json.Marshal(result)
			w.Write(jsonResult)
			return
		}

		flagNum := strings.Index(file_content, ";base64,")
		if flagNum < 0 {
			result = map[string]string{
				"status": "error",
				"msg":    "MimeType is empty",
			}

			jsonResult, _ = json.Marshal(result)
			w.Write(jsonResult)
			return
		}

		mime_type := parseFileMime(file_content)
		if _, ok := mimeMap[mime_type]; !ok {
			result = map[string]string{
				"status": "error",
				"msg":    "MimeType error",
			}
			jsonResult, _ = json.Marshal(result)
			w.Write(jsonResult)
			return
		}

		pos_num := strings.Index(file_content, ",")
		file_content = file_content[(pos_num + 1):]
		suffix := mimeMap[mime_type]
		//file_content = strings.Replace(file_content, "data:image/jpg;base64,", "", -1)

		oss.Init(OSS_ACCESS_ID, OSS_ACCESS_KEY, logger)
		path := generationPath(group, suffix)
		img_content, err := base64.StdEncoding.DecodeString(file_content)
		if err != nil {
			logger(err)
			result = map[string]string{
				"status": "error",
				"msg":    fmt.Sprintf("%v", err),
			}

			jsonResult, _ = json.Marshal(result)
			w.Write(jsonResult)
			return

		}
		statusCode := oss.MkDir(BUCKET, fmt.Sprintf("/%s", group))
		statusCode = oss.CreateMore(BUCKET, path, string(img_content), suffix)
		if statusCode != 200 {
			result = map[string]string{
				"status": "error",
				"msg":    "server occour error",
			}

			jsonResult, _ = json.Marshal(result)
			w.Write(jsonResult)
			return
		}

		logger(fmt.Sprintf("file [%s] upload success!", path))
		result = map[string]string{
			"status": "ok",
			"path":   path,
		}

		jsonResult, _ = json.Marshal(result)
		w.Write(jsonResult)
	}

}

func init() {
	var groups []string = []string{"voucher"}
	typeMap := map[string]string{
		"application/pdf": "pdf",
		"image/jpeg":      "jpg",
		"image/png":       "png",
	}
	for _, group := range groups {
		myHandlers[fmt.Sprintf("/%s", group)] = voucherHandler(group, typeMap)
	}
}
