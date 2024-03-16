package contants

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"
)

type fieldKey string
type FieldMap map[fieldKey]string

func (f FieldMap) resolve(key fieldKey) string {
	if k, ok := f[key]; ok {
		return k
	}

	return string(key)
}

// 自定义 formatter
type MyFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Force quoting of all values
	ForceQuote bool

	// DisableQuote disables quoting for all values.
	// DisableQuote will have a lower priority than ForceQuote.
	// If both of them are set to true, quote will be forced on all values.
	DisableQuote bool

	// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
	EnvironmentOverrideColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)

	// Disables the truncation of the level text to 4 characters.
	DisableLevelTruncation bool

	// PadLevelText Adds padding the level text so that all the levels output at the same length
	// PadLevelText is a superset of the DisableLevelTruncation option
	PadLevelText bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	// Whether the logger's out is to a terminal
	isTerminal bool

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &TextFormatter{
	//     FieldMap: FieldMap{
	//         FieldKeyTime:  "@timestamp",
	//         FieldKeyLevel: "@level",
	//         FieldKeyMsg:   "@message"}}
	FieldMap FieldMap

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	terminalInitOnce sync.Once

	// The max length of the level text, generated dynamically on init
	levelTextMaxLength int
}

// Format implement the Formatter interface
func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Caller.File = filepath.Base(entry.Caller.File)
	entry.Caller.Function = filepath.Base(entry.Caller.Function)

	data := make(logrus.Fields)
	for k, v := range entry.Data {
		data[k] = v
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	var funcVal, fileVal string
	fixedKeys := make([]string, 0, 4+len(data))

	if !f.DisableTimestamp {
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(logrus.FieldKeyTime))
	}
	fixedKeys = append(fixedKeys, f.FieldMap.resolve(logrus.FieldKeyLevel))
	if entry.Message != "" {
		fixedKeys = append(fixedKeys, f.FieldMap.resolve(logrus.FieldKeyMsg))
	}
	//if entry.err != "" {
	//	fixedKeys = append(fixedKeys, f.FieldMap.resolve(logrus.FieldKeyLogrusError))
	//}

	if entry.HasCaller() {
		funcVal = entry.Caller.Function
		fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)

		if funcVal != "" {
			fixedKeys = append(fixedKeys, f.FieldMap.resolve(logrus.FieldKeyFunc))
		}
		if fileVal != "" {
			fixedKeys = append(fixedKeys, f.FieldMap.resolve(logrus.FieldKeyFile))
		}
	}

	if !f.DisableSorting {
		if f.SortingFunc == nil {
			sort.Strings(keys)
			fixedKeys = append(fixedKeys, keys...)
		}
	} else {
		fixedKeys = append(fixedKeys, keys...)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}
	for _, key := range fixedKeys {
		var value interface{}
		switch {
		case key == f.FieldMap.resolve(logrus.FieldKeyTime):
			value = entry.Time.Format(timestampFormat)
		case key == f.FieldMap.resolve(logrus.FieldKeyLevel):
			value = entry.Level.String()
		case key == f.FieldMap.resolve(logrus.FieldKeyMsg):
			value = entry.Message
		//case key == f.FieldMap.resolve(logrus.FieldKeyLogrusError):
		//	value = entry.err
		case key == f.FieldMap.resolve(logrus.FieldKeyFunc) && entry.HasCaller():
			value = funcVal
		case key == f.FieldMap.resolve(logrus.FieldKeyFile) && entry.HasCaller():
			value = fileVal
		default:
			value = data[key]
		}
		f.appendKeyValue(b, key, value)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}
func (f *MyFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}
func (f *MyFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}
func (f *MyFormatter) needsQuoting(text string) bool {
	if f.ForceQuote {
		return true
	}
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	if f.DisableQuote {
		return false
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}
