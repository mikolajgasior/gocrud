package logger

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"reflect"
	"time"
)

const (
	logAttrError   = "err"
	logAttrModule  = "mod"
	logAttrService = "svc"
	logAttrHandler = "hand"
	logAttrPath    = "path"
)

func SetLogger(envValue string) {
	logLevel := logLevelFromEnv(envValue)
	slog.Info("log level from env", "level", logLevel)

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))
	slog.SetDefault(logger)
}

func AttrError(err error) slog.Attr {
	return slog.String(logAttrError, err.Error())
}

func AttrModule(obj interface{}) slog.Attr {
	name := reflect.TypeOf(obj).Elem().PkgPath() + "." + reflect.TypeOf(obj).Elem().Name()

	return slog.String(logAttrModule, name)
}

func AttrService(obj interface{}, action string) slog.Attr {
	name := reflect.TypeOf(obj).Elem().PkgPath() + "." + reflect.TypeOf(obj).Elem().Name() + "." + action

	return slog.String(logAttrService, name)
}

func AttrHandler(obj interface{}) slog.Attr {
	name := reflect.TypeOf(obj).Elem().PkgPath() + "." + reflect.TypeOf(obj).Elem().Name()

	return slog.String(logAttrHandler, name)
}

func AttrPath(path string) slog.Attr {
	return slog.String(logAttrPath, path)
}

func logLevelFromEnv(s string) slog.Level {
	switch s {
	case "err":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	case "debug":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

func generateRandomAlphanumeric(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func LogError(msg string, attrs ...interface{}) string {
	code := fmt.Sprintf("ERR-%s-%d", generateRandomAlphanumeric(8), time.Now().UnixNano())
	attrs = append(attrs, slog.String("err_code", code))
	slog.Error(msg, attrs...)
	return code
}
