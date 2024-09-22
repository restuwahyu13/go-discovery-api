package helpers

import (
	"net/http"

	"github.com/goccy/go-json"
)

type (
	Error struct {
		Name    string `json:"name,omitempty"`
		Message string `json:"message"`
		Code    int    `json:"code,omitempty"`
		Stack   any    `json:"stack,omitempty"`
	}

	Request[T any] struct {
		Req   T
		Body  T
		Param T
		Query T
	}

	Response struct {
		StatCode   float64        `json:"stat_code,omitempty"`
		StatMsg    any            `json:"stat_msg,omitempty"`
		ErrCode    any            `json:"err_code,omitempty"`
		ErrMsg     any            `json:"err_msg,omitempty"`
		Data       any            `json:"data,omitempty"`
		Errors     *[]Error       `json:"errors,omitempty"`
		Pagination map[string]any `json:"pagination,omitempty"`
	}
)

func ApiResponse(rw http.ResponseWriter, options *Response) {
	rw.Header().Set("Content-Type", "application/json")

	var (
		parser     IParser        = NewParser()
		res        Response       = Response{StatCode: http.StatusInternalServerError}
		config     map[string]any = nil
		errCode    string         = "GENERAL_ERROR"
		errMessage string         = "API is busy please try again later!"
	)

	optionsByte, err := parser.Marshal(options)
	if err != nil {
		res.ErrCode = errCode
		res.ErrCode = errMessage
	}

	if err := parser.Unmarshal(optionsByte, &config); err != nil {
		res.ErrCode = errCode
		res.ErrCode = errMessage
	}

	if config["stat_code"] == nil && config["stat_msg"] == nil && config["err_msg"] != nil {
		res.ErrCode = errCode
		res.ErrMsg = config["err_msg"]
	}

	if statCode := config["stat_code"]; statCode != nil {
		res.StatCode = statCode.(float64)
	}

	for key, value := range config {
		switch key {

		case "stat_msg":
			if v, ok := value.(string); ok {
				res.StatMsg = &v
			}

		case "err_code":
			if v, ok := value.(string); ok {
				res.ErrCode = &v
			}

		case "err_msg":
			res.ErrMsg = value

		case "data":
			res.Data = value

		case "pagination":
			if v, ok := value.(map[string]any); ok {
				res.Pagination = v
			}
		}
	}

	if options.StatCode >= 400 && options.StatCode <= 500 {
		rw.WriteHeader(int(options.StatCode))
		json.NewEncoder(rw).Encode(res)
	} else {
		json.NewEncoder(rw).Encode(res)
	}
}
