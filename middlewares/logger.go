package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"prs/models"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func LogrusLogger() echo.MiddlewareFunc {
	/* ... logger 초기화 */
	logger := logrus.New()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logEntry := logrus.NewEntry(logger)
			// var logResponse config.Logger
			data := make(map[string]interface{})

			// var httpBody *http.body

			// request_id를 가져와 logEntry에 셋팅
			id := c.Request().Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = c.Response().Header().Get(echo.HeaderXRequestID)
			}

			var getBodyData []string
			values, _ := c.FormParams()
			for k, v := range values {
				value := fmt.Sprintf("%s: %s", k, strings.Join(v, "&"))
				getBodyData = append(getBodyData, value)
			}

			form, err := c.MultipartForm()
			if err != nil {
				if err.Error() != "request Content-Type isn't multipart/form-data" {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
			}
			if form != nil {
				files := form.File["photo"]
				for idx, file := range files {
					value := fmt.Sprintf("photo%03d: %s", idx, file.Filename)
					getBodyData = append(getBodyData, value)
				}
			}

			// logrus에 저장
			data["request_id"] = id
			data["body"] = strings.Join(getBodyData, ", ")
			data["connect_ip"] = c.RealIP()
			data["request_url"] = c.Request().URL.RequestURI()
			data["user_agent"] = c.Request().UserAgent()

			logEntry = logEntry.WithFields(data)
			// logEntry를 Context에 저장
			req := c.Request()
			c.SetRequest(req.WithContext(
				context.WithValue(
					req.Context(),
					"LOG",
					logEntry,
				),
			))

			return next(c)
		}
	}
}

func CreateLogger(db *gorm.DB, logger *logrus.Entry, status int, err error) {

	logNew := logrus.New()

	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.WarnLevel)

	logConf := models.Logger{
		Body:       fmt.Sprintf("%s", logger.Data["body"]),
		ConnectIp:  fmt.Sprintf("%s", logger.Data["connect_ip"]),
		RequestId:  fmt.Sprintf("%s", logger.Data["request_id"]),
		RequestUrl: fmt.Sprintf("%s", logger.Data["request_url"]),
		Status:     status,
		Backoff:    time.Second.Milliseconds(),
		UserAgent:  fmt.Sprintf("%s", logger.Data["user_agent"]),
		CreatedAt:  time.Now(),
	}

	log := logger.WithFields(logrus.Fields{
		"backoff":     logConf.Backoff,
		"body":        logConf.Body,
		"created":     logConf.CreatedAt,
		"IP":          logConf.ConnectIp,
		"request-id":  logConf.RequestId,
		"request-url": logConf.RequestUrl,
		"status":      logConf.Status,
		"user-agent":  logConf.UserAgent,
	})

	if err != nil {
		log.Error(err.Error())
		logConf.Message = err.Error()
	} else {
		log.Info("")
	}

	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.
	logNew.Out = os.Stdout

	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile(fmt.Sprintf("log/log_%s.log", time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logNew.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// db.Create(&logConf)
}
