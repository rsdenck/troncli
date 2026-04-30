package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Output struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Total   int         `json:"total,omitempty"`
	Items   interface{} `json:"items,omitempty"`
}

func NewSuccess(data interface{}) *Output {
	return &Output{
		Status: "success",
		Data:   data,
	}
}

func NewError(err string, code string) *Output {
	return &Output{
		Status: "error",
		Error:  err,
		Code:   code,
	}
}

func NewList(items interface{}, total int) *Output {
	return &Output{
		Status: "success",
		Total:  total,
		Items:  items,
	}
}

func (o *Output) WithMessage(msg string) *Output {
	o.Message = msg
	return o
}

func (o *Output) Print(jsonOutput bool) {
	if jsonOutput {
		data, _ := json.MarshalIndent(o, "", "  ")
		fmt.Println(string(data))
		return
	}
	
	if o.Status == "success" {
		if o.Message != "" {
			fmt.Println(o.Message)
		}
		if o.Items != nil {
			switch items := o.Items.(type) {
			case []map[string]interface{}:
				for _, item := range items {
					for k, v := range item {
						fmt.Printf("%s: %v\n", k, v)
					}
					fmt.Println("---")
				}
			case []string:
				for _, item := range items {
					fmt.Println(item)
				}
			default:
				if o.Data != nil {
					data, _ := json.MarshalIndent(o.Data, "", "  ")
					fmt.Println(string(data))
				}
			}
		} else if o.Data != nil {
			data, _ := json.MarshalIndent(o.Data, "", "  ")
			fmt.Println(string(data))
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", o.Error)
	}
}
