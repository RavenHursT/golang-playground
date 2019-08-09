package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type ParsedJSONMap map[string]interface{} 

type RequestLogItem struct {
	Method string `json:"method"`
	Headers http.Header `json:"headers"`
	Body ParsedJSONMap `json:"body"`
}

type responseLogItem struct {
	Status  int         `json:"status" xml:"status" form:"status" query:"status"`
	Headers http.Header `json:"headers" xml:"headers" form:"headers" query:"headers"`
	Body ParsedJSONMap `json:"body"`
}

func RequestResponseLogger(context echo.Context, reqBody []byte, respBody []byte) {
	var jsonParseError error
	var reqBodyMap map[string]interface{}
	jsonParseError = json.Unmarshal([]byte(reqBody), &reqBodyMap)
	if jsonParseError != nil {
		reqBodyMap = make(map[string]interface{})
	}
	req := context.Request()
	reqLogItem := RequestLogItem{req.Method, req.Header, reqBodyMap}
	reqLogItemJSON, _ := json.Marshal(reqLogItem)
	fmt.Printf("%s\n", reqLogItemJSON)

	var respBodyMap map[string]interface{}
	jsonParseError = json.Unmarshal([]byte(respBody), &respBodyMap)
	if jsonParseError != nil {
		respBodyMap = make(map[string]interface{})
	}
	resp := context.Response()
	respLogItem := responseLogItem{resp.Status, resp.Header(), respBodyMap}
	respLogItemJSON, _ := json.Marshal(respLogItem)
	fmt.Printf("%s\n", respLogItemJSON)
}
