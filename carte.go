package carte

import (
	"strconv"
	"strings"
	"time"
)

// Made the log one line to allow for a quick synchronous write to a custom writer with mutex locking
func log(sev *severity, details ...Jsonable) (int, error) {
	w := sev.GetWriter()
	// Ignore nil writers
	if w == nil {
		return 0, nil
	}

	// Get func name
	funcName := getCaller()

	// Rough estimate of all the required wrappers to a log, NOT the info
	// TODO: Cant efficiently calculate the size of the final log, find another way to save allocations
	// Can most likely calculate it now, will need to reorder
	// Reduces the number of allocations, essentially just a guess since the latest patches
	baseLogLen := 42
	if len(details) > 0 {
		baseLogLen += 9
		baseLogLen += len(details) * 6
	}
	jsonLog := make([]byte, 0, baseLogLen)

	// DATE
	jsonLog = append(jsonLog, `{"TS":`...)
	jsonLog = append(jsonLog, strconv.FormatInt(time.Now().Unix(), 10)...)

	// FUNC
	if funcName != nil {
		jsonLog = append(jsonLog, `,"FN":"`...)
		jsonLog = append(jsonLog, funcName...)
		jsonLog = append(jsonLog, '"')
	}

	// TYPE
	jsonLog = append(jsonLog, `,"SV":"`...)
	jsonLog = append(jsonLog, sev.name...)
	jsonLog = append(jsonLog, '"')

	// DETAILS
	for i, dtl := range details {
		if dtl == nil {
			continue
		}

		name, val := dtl.Json()
		name = strings.ReplaceAll(name, "\"", "'")
		name = strings.ReplaceAll(name, "\n", "")
		val = strings.ReplaceAll(val, "\"", "'")
		val = strings.ReplaceAll(val, "\n", "")

		jsonLog = append(jsonLog, ",\""...)
		jsonLog = append(jsonLog, name...)
		jsonLog = append(jsonLog, `":"`...)
		jsonLog = append(jsonLog, val...)
		jsonLog = append(jsonLog, '"')

		// If this is not the last detail append a comma
		if i != len(details)-1 {
		}
	}

	// Newlines break json parsing
	jsonLog = []byte(strings.ReplaceAll(string(jsonLog), "\n", ""))
	jsonLog = append(jsonLog, "}\n"...)

	if hook := sev.getHook(); hook != nil {
		go hook(jsonLog)
	}

	// TODO: add ability to write to StdOut as well if the writer is not
	// Probably in settings

	n, err := w.Write(jsonLog)
	// Panic
	if sev == Panic {
		panic("carte.Panic called")
	}
	return n, err
}
