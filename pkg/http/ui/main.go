package ui

import (
	"embed"
	"fmt"
	"net/http"
	"reflect"
	"regexp"

	"codeberg.org/mikolajgasior/gocrud/pkg/http/cors"
	svccrud "codeberg.org/mikolajgasior/gocrud/pkg/service"
)

const (
	CodeServiceError = "SERVICE_ERROR"
)

const (
	pathHome  = "/"
	pathLogin = "/login/"
)

const (
	pathPartList   = "list"
	pathPartCreate = "create"
	pathPartUpdate = "update"
	pathPartDelete = "delete"
)

const (
	htmlInputIdPrefix   = "field_input_"
	htmlInputNamePrefix = "field_"
)

const (
	filterValPrefix = "filter_val_"
	filterOpPrefix  = "filter_op_"
)

type ContextValue string

const (
	CtxValueUserID = "LoggedUserID"
)

var (
	pathCreateRegexp      = regexp.MustCompile("^" + pathPartCreate + "[/]{0,1}$")
	pathFormCreateRegexp  = regexp.MustCompile("^form/" + pathPartCreate + "[/]{0,1}$")
	pathReadRegexp        = regexp.MustCompile("^read/[0-9]+[/]{0,1}$")
	pathUpdateRegexp      = regexp.MustCompile("^" + pathPartUpdate + "/[0-9]+[/]{0,1}$")
	pathFormUpdateRegexp  = regexp.MustCompile("^form/" + pathPartUpdate + "/[0-9]+[/]{0,1}$")
	pathFetchDeleteRegexp = regexp.MustCompile("^fetch/" + pathPartDelete + "/[0-9]+[/]{0,1}$")
	pathListRegexp        = regexp.MustCompile("^" + pathPartList + "[/]{0,1}$")
)

var (
	filterValRegexp = regexp.MustCompile("^filter_val_[a-zA-Z0-9_]+$")
	filterOpRegexp  = regexp.MustCompile("^filter_op_[a-zA-Z0-9_]+$")
)

//go:embed html
var embedHTML embed.FS

type Handler struct {
	embedHTML  embed.FS
	cors       *cors.CORS
	paths      map[string]func() interface{}
	svc        *svccrud.CRUD
	layout     Layout
	pathPrefix string
}

type HandlerInput struct {
	EmbedHTML  *embed.FS
	Svc        *svccrud.CRUD
	CORS       *cors.CORS
	Paths      map[string]func() interface{}
	Layout     Layout
	PathPrefix string
}

func New(input HandlerInput) *Handler {
	handler := &Handler{
		cors:       input.CORS,
		paths:      input.Paths,
		svc:        input.Svc,
		layout:     input.Layout,
		pathPrefix: input.PathPrefix,
	}

	if input.EmbedHTML != nil {
		handler.embedHTML = *input.EmbedHTML
	} else {
		handler.embedHTML = embedHTML
	}

	return handler
}

func redirect(w http.ResponseWriter, urlPath string) {
	w.Header().Set("Location", urlPath)
	w.WriteHeader(http.StatusSeeOther)
}

func requestUser(r *http.Request) (string, string) {
	ctxID := r.Context().Value(ContextValue(CtxValueUserID))

	userID := "0"
	if ctxID != nil {
		userID = fmt.Sprintf("%d", ctxID)
	}
	userName := ""

	return userID, userName
}

func structName(obj interface{}) string {
	objValue := reflect.ValueOf(obj)
	objIndirectValue := reflect.Indirect(objValue)
	objType := objIndirectValue.Type()

	if objType.String() == "reflect.Value" {
		objType = reflect.ValueOf(obj.(reflect.Value).Interface()).Type().Elem().Elem()
	}

	objTypeName := objType.Name()

	return objTypeName
}

func structFieldNames(obj interface{}) []string {
	objValue := reflect.ValueOf(obj)
	objIndirectValue := reflect.Indirect(objValue)

	if objIndirectValue.Kind() != reflect.Struct {
		return []string{}
	}

	objType := objIndirectValue.Type()

	if objType.String() == "reflect.Value" {
		objType = reflect.ValueOf(obj.(reflect.Value).Interface()).Type().Elem().Elem()
	}

	fieldNames := make([]string, 0, objType.NumField())
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldNames = append(fieldNames, field.Name)
	}

	return fieldNames
}

func structToMap(s interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		out[t.Field(i).Name] = v.Field(i).Interface()
	}
	return out
}
