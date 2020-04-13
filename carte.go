package carte

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Jsonable interface {
	Json() (name string, value string)
}

// Log writes the writer specified in the severity struct
func Log(sev *severity, message string, details ...interface{}) {
	_, _ = log(sev, message, details...)
}

// Rlog is the same as log, but returns the values from the call to Write()
func Rlog(sev *severity, message string, details ...interface{}) (int, error) {
	return log(sev, message, details...)
}

// Made the log one line to allow for a quick synchronous write to a custom writer with mutex locking
func log(sev *severity, message string, details ...interface{}) (int, error) {
	// Get date
	date := getDate()

	// Get func name
	funcName := getFuncName()

	// Rough estimate of all the required wrappers to a log, NOT the info
	// TODO: Cant efficiently calculate the size of the final log, find another way to save allocations
	// Reduces the number of allocations
	// Base of 42 + (if len dtls > 0) -> 9 + len dtls * 6
	baseLogLen := 42
	if len(details) > 0 {
		baseLogLen += 9
		baseLogLen += len(details) * 6
	}
	jsonLog := make([]byte, 0, baseLogLen)

	// DATE
	jsonLog = append(jsonLog, `{"Time":"`...)
	jsonLog = append(jsonLog, date...)

	// FUNC
	if funcName != nil {
		jsonLog = append(jsonLog, `","Func":"`...)
		jsonLog = append(jsonLog, funcName...)
	}

	// TYPE
	jsonLog = append(jsonLog, `","Severity":"`...)
	jsonLog = append(jsonLog, sev.name...)

	// MESSAGE
	jsonLog = append(jsonLog, `","Message":"`...)
	jsonLog = append(jsonLog, message...)
	jsonLog = append(jsonLog, '"')

	// DETAILS
	if len(details) > 0 {
		jsonLog = append(jsonLog, `,"Details":{`...)
		for i, dtl := range details {
			if dtl == nil {
				continue
			}

			// Conversions get more expensive further down the type switch
			// Jsonable is the fastest and clearest in expressing information
			var name, value string
			switch v := dtl.(type) {
			case Jsonable:
				name, value = v.Json()
			case fmt.Stringer:
				name = reflect.TypeOf(dtl).String()
				value = v.String()
			case error:
				name = reflect.TypeOf(dtl).String()
				value = strings.ReplaceAll(v.Error(), "\"", "'")
			default:
				name = reflect.TypeOf(dtl).String()
				if val, err := json.Marshal(v); err != nil {
					value = strings.ReplaceAll(err.Error(), "\"", "'")
				} else {
					value = strings.ReplaceAll(string(val), "\"", "\\\"")
				}
			}

			jsonLog = append(jsonLog, '"')
			jsonLog = append(jsonLog, name...)
			jsonLog = append(jsonLog, `":"`...)
			jsonLog = append(jsonLog, value...)
			jsonLog = append(jsonLog, '"')

			// If this is not the last detail append a comma
			if i != len(details)-1 {
				jsonLog = append(jsonLog, ","...)
			}
		}
		jsonLog = append(jsonLog, '}')
	}
	jsonLog = append(jsonLog, "}\n"...)

	if hook := sev.getHook(); hook != nil {
		go hook(jsonLog)
	}

	// TODO: add ability to write to StdOut as well if the writer is not
	// Probably in settings

	return sev.GetWriter().Write(jsonLog)
}
