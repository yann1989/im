package logger

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestName(t *testing.T) {
	InitLogger()
	logrus.WithFields(logrus.Fields{
		"url":    "request.Request.RequestURI",
		"method": "request.Request.Method",
		"ip":     "request.Request.RemoteAddr",
	}).Error("ceshi")
}
