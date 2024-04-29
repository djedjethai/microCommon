package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/grpasr/common/tests"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	outC     chan string
	r        *os.File
	w        *os.File
	buf      bytes.Buffer
	out      string
	res      []string
	logEntry LogEntry
)

type LogEntry struct {
	Level      string    `json:"level"`
	Service    string    `json:"service"`
	KeyStr     string    `json:"keyStr"`
	KeyInt     int       `json:"keyInt"`
	KeyFloat64 float64   `json:"keyFloat64"`
	Err        string    `json:"error"`
	Time       time.Time `json:"time"`
	Caller     string    `json:"caller"`
	Message    string    `json:"message"`
}

func Test_production_logs_from_facade(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)
	ret := teardown()

	logEnv = production

	lf := NewLoggingFacade("serviceTest")

	lf.NewLogHandler(debugLevel).Msg("DebugLevel")
	fmt.Fprint(os.Stdout, "::")
	lf.NewLogHandler(infoLevel).Msg("InfoLevel")
	fmt.Fprint(os.Stdout, "::")
	lf.NewLogHandler(warnLevel).Msg("WarnLevel")
	fmt.Fprint(os.Stdout, "::")
	lf.NewLogHandler(errorLevel).Msg("ErrorLevel")
	// NewLogHandler(FatalLevel).Msg("FatalLevel") // stop the prog exec
	// fmt.Fprint(os.Stdout, "::")
	// NewLogHandler(PanicLevel).Msg("PanicLevel")

	ret()
	res = outputSplitter()

	tests.MaybeFail("create_new_codeError_basic", tests.Expect(res[0], ""))
	tests.MaybeFail("create_new_codeError_basic", tests.Expect(res[1], ""))
	tests.MaybeFail("create_new_codeError_basic", tests.Expect(res[2], ""))
	err := json.Unmarshal([]byte(res[3]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "error"),
		tests.Expect(logEntry.Service, "serviceTest"),
		tests.Expect(logEntry.Message, "ErrorLevel"),
	)

	res = []string{}
	buf.Reset()
}

func Test_development_logs_from_facade(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)
	ret := teardown()

	lf := NewLoggingFacade("serviceTest")
	lf.SetLoggingEnvToDevelopment()

	lf.NewLogHandler(debugLevel).Msg("DebugLevel")
	fmt.Fprint(os.Stdout, "::")
	lf.NewLogHandler(infoLevel).Msg("InfoLevel")
	fmt.Fprint(os.Stdout, "::")
	lf.NewLogHandler(warnLevel).Msg("WarnLevel")
	fmt.Fprint(os.Stdout, "::")
	lf.NewLogHandler(errorLevel).Msg("ErrorLevel")
	fmt.Fprint(os.Stdout, "::")
	// NewLogHandler(FatalLevel).Msg("FatalLevel") // stop the prog exec
	// fmt.Fprint(os.Stdout, "::")
	// NewLogHandler(PanicLevel).Msg("PanicLevel")

	ret()
	res = outputSplitter()

	err := json.Unmarshal([]byte(res[0]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "debug"),
		tests.Expect(logEntry.Message, "DebugLevel"),
	)
	err = json.Unmarshal([]byte(res[1]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "info"),
		tests.Expect(logEntry.Message, "InfoLevel"),
	)
	err = json.Unmarshal([]byte(res[2]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "warn"),
		tests.Expect(logEntry.Message, "WarnLevel"),
	)
	err = json.Unmarshal([]byte(res[3]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "error"),
		tests.Expect(logEntry.Message, "ErrorLevel"),
	)

	res = []string{}
	buf.Reset()

}

func Test_Msgf_log_from_facade(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)
	ret := teardown()

	lf := NewLoggingFacade("serviceTest")

	lf.NewLogHandler(errorLevel).Float64("KeyFloat64", 55.55).Msgf("response", 55)
	fmt.Fprint(os.Stdout, "::")
	lf.NewLogHandler(errorLevel).Float64("keyFloat64", 55.55).Msgf("success", "successMsg")

	ret()
	res = outputSplitter()

	err := json.Unmarshal([]byte(res[0]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "error"),
		tests.Expect(logEntry.KeyFloat64, 55.55),
		tests.Expect(logEntry.Message, "response: 55"),
	)
	err = json.Unmarshal([]byte(res[1]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "error"),
		tests.Expect(logEntry.KeyFloat64, 55.55),
		tests.Expect(logEntry.Message, "success: successMsg"),
	)

	res = []string{}
	buf.Reset()
}

func Test_Setting_env_log_from_facade(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	lf := NewLoggingFacade("serviceTest")

	lf.SetLoggingEnvToDevelopment()
	tests.MaybeFail("test_env_to_development", nil,
		tests.Expect(lf.GetLoggingEnv(), "development"),
	)

	lf.SetLoggingEnvToProduction()
	tests.MaybeFail("test_env_to_production", nil,
		tests.Expect(lf.GetLoggingEnv(), "production"),
	)

	lf.SetLoggingEnvToDevelopment()
	tests.MaybeFail("test_env_to_development", nil,
		tests.Expect(lf.GetLoggingEnv(), "development"),
	)
}

func teardown() func() {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ = os.Pipe()
	os.Stdout = w

	outC = make(chan string, 1) // Adding a buffer size of 1

	return func() {
		w.Close()
		os.Stdout = old  // restoring the real stdout
		io.Copy(&buf, r) // Read from the pipe
		outC <- buf.String()
	}
}

func outputSplitter() []string {
	select {
	case out = <-outC:
		return strings.Split(out, "::")
	default:
		return []string{}
	}
}
