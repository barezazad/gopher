package middleware

import (
	"bytes"
	"encoding/json"
	"gopher/internal/core"
	"gopher/pkg/generr"
	"io"
	"io/ioutil"
	"math"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// APILogger is used to save requests and response by using logapi
func APILogger(engine *core.Engine) gin.HandlerFunc {
	var reqID uint
	env := engine.Environments.API

	log, err := generr.New(env.LogFormat, env.LogOutput, env.LogLevel, env.LogIndentation, true)
	if err != nil {
		engine.ServerLog.Fatal(err)
	}

	return func(c *gin.Context) {
		start := time.Now()
		buf, _ := ioutil.ReadAll(c.Request.Body)
		reqDataReader := ioutil.NopCloser(bytes.NewBuffer(buf))

		//We have to create a new Buffer, because reqDataReader will be read.
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		reqID++

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		loggingRequest(log, c, reqID, reqDataReader)

		c.Next()

		latency := int(math.Ceil(float64(time.Since(start).Nanoseconds()) / 1000000.0))

		loggingResponse(log, c, latency, blw)

	}
}

// Logging Request
func loggingRequest(log *generr.Log, c *gin.Context, reqID uint, reqDataReader io.Reader) {
	//TODO
	request := getBody(log, reqDataReader)
	// prevent to save the passwords
	if strings.Contains(c.Request.URL.Path, "login") {
		request = nil
	}

	log.Logger.WithFields(logrus.Fields{
		"reqID":      reqID,
		"ip":         c.Request.Header.Get("X-User-IP"), //r.Context.ClientIP(),
		"method":     c.Request.Method,
		"uri":        c.Request.RequestURI,
		"path":       c.Request.URL.Path,
		"request":    request,
		"params":     c.Request.URL.Query(),
		"referer":    c.Request.Referer(),
		"user_agent": c.Request.UserAgent(),
	}).Info("request")
	c.Set("resID", reqID)
}

// Logging Response
func loggingResponse(log *generr.Log, c *gin.Context, latency int, blw *bodyLogWriter) {
	resID, ok := c.Get("resID")
	if !ok {
		log.Debug("there is no resIndex for element", getBody(log, blw.body))
	}
	log.Logger.WithFields(logrus.Fields{
		"resID":       resID,
		"status":      c.Writer.Status(),
		"latency":     latency, // time to process
		"data_length": c.Writer.Size(),
		"response":    getBody(log, blw.body),
	}).Info("response")
}

// ActivityRead body
func getBody(log *generr.Log, reader io.Reader) interface{} {

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		log.Debug(err)
	}

	var obj interface{}

	if err := json.NewDecoder(buf).Decode(&obj); err != nil {
		if err.Error() != "EOF" {
			log.Info(err, obj, err.Error())
		}
	}

	return obj
}
