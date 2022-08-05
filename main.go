package main

import (
	"context"
	"fmt"
	"github.com/bytedance/go-dyclog"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func hello(w http.ResponseWriter, req *http.Request) {
	os.Stdout.WriteString("Msg to hello\n")
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	os.Stdout.WriteString("Msg to headers\n")
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func err1(w http.ResponseWriter, req *http.Request) {
	os.Stderr.WriteString("Msg to err1\n")
	http.Error(w, "this is an err interface", 500)
}

func err2(w http.ResponseWriter, req *http.Request) {
	os.Stderr.WriteString("Msg to err2\n")
	http.Error(w, "this is an err interface", 404)
}

func ping(w http.ResponseWriter, req *http.Request) {
	os.Stdout.WriteString("Msg to ping\n")
	fmt.Fprintf(w, "pong!\n")
}

func body(w http.ResponseWriter, req *http.Request) {
	b, e := ioutil.ReadAll(req.Body)
	fmt.Fprintf(w, "body: %v, err: %v", string(b), e)
}

func testPanic(w http.ResponseWriter, req *http.Request) {
	os.Stderr.WriteString("Msg to testPanic\n")
	panic(req)
}

func log(w http.ResponseWriter, req *http.Request) {
	var record string
	v, e := url.ParseQuery(req.URL.RawQuery)
	if e != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", e)
		return
	}
	if v["type"] == nil {
		fmt.Fprintln(os.Stderr, "err: type param is not valid")
		return
	}
	switch v["type"][0] {
	case "D":
		record = "DEBUG 1658217911838250001 example.go:66 10.79.163.90 debug level test!\n"
	case "I":
		record = "INFO 1658217911838250002 example.go:66 10.79.163.90 info level test!\n"
	case "N":
		record = "NOTICE 1658217911838250003 example.go:66 10.79.163.90 notice level test!\n"
	case "E":
		record = "ERROR 1658217911838250004 example.go:66 10.79.163.90 error level test!\n"
	case "W":
		record = "WARN 1658217911838250005 example.go:66 10.79.163.90 warn level test!\n"
	case "F":
		record = "FATAL 1658217911838250006 example.go:66 10.79.163.90 fatal level test!\n"
	default:
	}
	fmt.Fprintln(os.Stdout, record)
	fmt.Fprintf(w, "body: %v, err: %v\n", record, e)
}

func sdkLog(w http.ResponseWriter, req *http.Request) {
	var record string
	v, e := url.ParseQuery(req.URL.RawQuery)
	if e != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", e)
		return
	}
	if v["type"] == nil {
		fmt.Fprintln(os.Stderr, "err: type param is not valid")
		return
	}
	switch v["type"][0] {
	case "D":
		dyclog.Debug("This is Debug log")
	case "I":
		dyclog.Info("This is Info log")
	case "N":
		dyclog.Notice("This is Notice log")
	case "E":
		dyclog.Error("This is Error log")
	case "W":
		dyclog.Warn("This is Warn log")
	case "F":
		dyclog.Fatal("This is Fatal log")
	default:
	}
	fmt.Fprintln(os.Stdout, record)
	fmt.Fprintf(w, "body: %v, err: %v\n", record, e)
}

func sdkCtxLog(w http.ResponseWriter, req *http.Request) {
	ctx := dyclog.InjectLogIDToCtx(context.Background(), req.Header.Get("X-Tt-Logid"))
	var record string
	v, e := url.ParseQuery(req.URL.RawQuery)
	if e != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", e)
		return
	}
	if v["type"] == nil {
		fmt.Fprintln(os.Stderr, "err: type param is not valid")
		return
	}
	switch v["type"][0] {
	case "D":
		dyclog.CtxDebug(ctx, "This is Debug log")
	case "I":
		dyclog.CtxInfo(ctx, "This is Info log")
	case "N":
		dyclog.CtxNotice(ctx, "This is Notice log")
	case "E":
		dyclog.CtxError(ctx, "This is Error log")
	case "W":
		dyclog.CtxWarn(ctx, "This is Warn log")
	case "F":
		dyclog.CtxFatal(ctx, "This is Fatal log")
	default:
	}
	fmt.Fprintln(os.Stdout, record)
	fmt.Fprintf(w, "body: %v, err: %v\n", record, e)
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/v1/ping", ping)
	http.HandleFunc("/err/500", err1)
	http.HandleFunc("/err/404", err2)
	http.HandleFunc("/vi/body", body)
	http.HandleFunc("/panic", testPanic)
	http.HandleFunc("/log", log)
	http.HandleFunc("/sdk/log", sdkLog)
	http.HandleFunc("/sdk/ctx/log", sdkCtxLog)

	http.ListenAndServe(":8000", nil)
}
