package logging

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/grpasr/common/tests"
	"os"
	"testing"
)

// NOTE this MUST be the first test
func Test_service_name_registration_and_default_variables(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	_ = NewLoggingFacade()
	tests.MaybeFail("service_name_should_be_undefined", tests.Expect(serviceName, "serviceName undefined"))

	lf := NewLoggingFacade("service")
	tests.MaybeFail("service_name_should_be_service", tests.Expect(serviceName, "service"))

	lf.SetSvcName("updated")
	tests.MaybeFail("service_name_should_be_updated", tests.Expect(serviceName, "updated"))

	// test the default logEnv
	tests.MaybeFail("service_name_should_be_updated", tests.Expect(logEnv, production))
}

func Test_production_logs(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)
	ret := teardown()

	_ = NewLoggingFacade("serviceTest")

	newLogHandler(debugLevel).Msg("DebugLevel")
	fmt.Fprint(os.Stdout, "::")
	newLogHandler(infoLevel).Msg("InfoLevel")
	fmt.Fprint(os.Stdout, "::")
	newLogHandler(warnLevel).Msg("WarnLevel")
	fmt.Fprint(os.Stdout, "::")
	newLogHandler(errorLevel).Msg("ErrorLevel")
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

func Test_development_logs(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)
	ret := teardown()

	lf := NewLoggingFacade("serviceTest")
	lf.SetLoggingEnvToDevelopment()

	newLogHandler(debugLevel).Msg("DebugLevel")
	fmt.Fprint(os.Stdout, "::")
	newLogHandler(infoLevel).Msg("InfoLevel")
	fmt.Fprint(os.Stdout, "::")
	newLogHandler(warnLevel).Msg("WarnLevel")
	fmt.Fprint(os.Stdout, "::")
	newLogHandler(errorLevel).Msg("ErrorLevel")
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

func Test_Msgf_log(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)
	ret := teardown()

	_ = NewLoggingFacade("serviceTest")

	newLogHandler(errorLevel).Float64("KeyFloat64", 55.55).Msgf("response", 55)
	fmt.Fprint(os.Stdout, "::")
	newLogHandler(errorLevel).Float64("keyFloat64", 55.55).Msgf("success", "successMsg")

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

func Test_Float64_log(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)
	ret := teardown()

	_ = NewLoggingFacade("serviceTest")

	newLogHandler(errorLevel).Float64("KeyFloat64", 55.55).Send()
	fmt.Fprint(os.Stdout, "::")
	newLogHandler(errorLevel).Float64("keyFloat64", 55.55).Msg("ErrorLevel")

	ret()
	res = outputSplitter()

	err := json.Unmarshal([]byte(res[0]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "error"),
		tests.Expect(logEntry.KeyFloat64, 55.55),
	)
	err = json.Unmarshal([]byte(res[1]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "error"),
		tests.Expect(logEntry.KeyFloat64, 55.55),
		tests.Expect(logEntry.Message, "ErrorLevel"),
	)

	res = []string{}
	buf.Reset()
}

func Test_Err_log(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)
	ret := teardown()

	_ = NewLoggingFacade("serviceTest")

	errT := errors.New("error test")

	newLogHandler(errorLevel).Err(errT).Send()
	fmt.Fprint(os.Stdout, "::")
	newLogHandler(errorLevel).Err(errT).Msg("ErrorLevel")

	ret()
	res = outputSplitter()

	err := json.Unmarshal([]byte(res[0]), &logEntry)

	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "error"),
		tests.Expect(logEntry.Err, "error test"),
	)
	err = json.Unmarshal([]byte(res[1]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "error"),
		tests.Expect(logEntry.Err, "error test"),
		tests.Expect(logEntry.Message, "ErrorLevel"),
	)

	res = []string{}
	buf.Reset()
}

func Test_Int_log(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)
	ret := teardown()

	_ = NewLoggingFacade("serviceTest")

	newLogHandler(errorLevel).Int("KeyInt", 555).Send()
	fmt.Fprint(os.Stdout, "::")
	newLogHandler(errorLevel).Int("keyInt", 555).Msg("ErrorLevel")

	ret()
	res = outputSplitter()

	err := json.Unmarshal([]byte(res[0]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "error"),
		tests.Expect(logEntry.KeyInt, 555),
	)
	err = json.Unmarshal([]byte(res[1]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "error"),
		tests.Expect(logEntry.KeyInt, 555),
		tests.Expect(logEntry.Message, "ErrorLevel"),
	)

	res = []string{}
	buf.Reset()
}

func Test_Str_log(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)
	ret := teardown()

	lf := NewLoggingFacade("serviceTest")
	lf.SetLoggingEnvToDevelopment()

	newLogHandler(infoLevel).Str("keyStr", "valStr").Send()
	fmt.Fprint(os.Stdout, "::")
	newLogHandler(errorLevel).Str("keyStr", "valStr").Msg("ErrorLevel")

	ret()
	res = outputSplitter()

	err := json.Unmarshal([]byte(res[0]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "info"),
		tests.Expect(logEntry.KeyStr, "valStr"),
	)
	err = json.Unmarshal([]byte(res[1]), &logEntry)
	tests.MaybeFail("create_new_codeError_basic", err,
		tests.Expect(logEntry.Level, "error"),
		tests.Expect(logEntry.KeyStr, "valStr"),
		tests.Expect(logEntry.Message, "ErrorLevel"),
	)

	res = []string{}
	buf.Reset()
}
