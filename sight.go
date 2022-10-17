package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	oss "github.com/kyf/oss_sdk/lib"
	"net/http"
	"strings"
)

func returnHandler(group string) MyHandler {
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

		file_content = strings.Replace(file_content, "data:image/jpg;base64,", "", -1)

		oss.Init(OSS_ACCESS_ID, OSS_ACCESS_KEY, logger)
		path := generationPath(group, "jpg")
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
		statusCode = oss.Create(BUCKET, path, string(img_content))
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
	var groups []string = []string{"sight", "hotel"}
	for _, group := range groups {
		myHandlers[fmt.Sprintf("/%s", group)] = returnHandler(group)
	}
}
