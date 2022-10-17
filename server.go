package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	PORT    string = ":8010"
	SSLPORT string = ":4431"
	BUCKET  string = "/6ry-poi"
	LAYOUT  string = "2006_01_02"
	LOG_DIR string = "./log/"

	CERT_FILE = "/cert/6ry_cert/6ry.crt"
	KEY_FILE  = "/cert/6ry_cert/6ry.key"
)

type MyHandler func(http.ResponseWriter, *http.Request)

var (
	ROUTER           map[*regexp.Regexp]MyHandler = make(map[*regexp.Regexp]MyHandler)
	myHandlers       map[string]MyHandler         = make(map[string]MyHandler)
	exit             chan int                     = make(chan int)
	log_chan         chan string                  = make(chan string, 100)
	last_handler     *os.File                     = nil
	last_logger      *log.Logger                  = nil
	last_logger_date string                       = ""

	OSS_ACCESS_ID  string = os.Getenv("OSS_ACCESS_ID")
	OSS_ACCESS_KEY string = os.Getenv("OSS_ACCESS_KEY")
)

func logContent(log_str string) {
	current := time.Now().Format(LAYOUT)
	if !strings.EqualFold(current, last_logger_date) {
		if last_handler != nil {
			last_handler.Close()
		}
		last_logger_date = current
		var err error
		last_handler, err = os.OpenFile(fmt.Sprintf("%s%s", LOG_DIR, last_logger_date), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		last_logger = log.New(last_handler, "", log.LstdFlags)
	}

	last_logger.Print(log_str)
}

func startLogger() {
	for {
		select {
		case log_str := <-log_chan:
			logContent(log_str)
		case <-exit:
			return
		}
	}
}

func logger(content interface{}) {
	if strcontent, ok := content.([]byte); ok {
		log_chan <- string(strcontent)
		return
	}

	if strcontent, ok := content.(string); ok {
		log_chan <- strcontent
		return
	}

	log_chan <- fmt.Sprintf("%v", content)
}

func main() {
	go startLogger()
	defer func() {
		if last_handler != nil {
			last_handler.Close()
		}
		exit <- 1
	}()

	r := NewRouter()

	logger("server is startting ...")
	go func() {
		err := http.ListenAndServeTLS(SSLPORT, CERT_FILE, KEY_FILE, r)
		if err != nil {
			logger(err)
		}
	}()
	err := http.ListenAndServe(PORT, r)
	if err != nil {
		logger(err)
	}
}
