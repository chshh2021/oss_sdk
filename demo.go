package main

import (
	"fmt"
	oss "github.com/kyf/oss_sdk/lib"
	"net/http"
)

func init() {
	var group string = "demo"
	myHandlers[fmt.Sprintf("/%s", group)] = func(w http.ResponseWriter, r *http.Request) {
		oss.Init(OSS_ACCESS_ID, OSS_ACCESS_KEY, logger)
		/*
			oss.RmDir(BUCKET, "/sight")
			oss.Remove(BUCKET, "/sight/2015_10_14_14447985047668171981298498081.jpg")
			oss.RemoveBucket(BUCKET)
		*/
		//oss.CreateBucket(BUCKET)
		oss.Remove(BUCKET, "/sight/2015_10_14_14447987323049666151298498081.jpg")
		oss.Remove(BUCKET, "/sight/2015_10_14_14447988769221438841298498081.jpg")
		oss.Remove(BUCKET, "/sight/2015_10_14_14448006768876860831298498081.jpg")
		oss.Remove(BUCKET, "/sight/2015_10_14_14448011257944723231298498081.jpg")
		w.Write([]byte("done"))
	}
}
