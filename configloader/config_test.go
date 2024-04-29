package configloader

import (
	"gitlab.com/grpasr/common/tests"
	"testing"
)

// TODO
// test case all GOENV

func Test_create_config_with_LoaderConfig_instance(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	c, errConf := NewConfig("invalidEnv", "TLoaderConfig")

	err := c.LDRLoadConfigs("configTests", "yaml", "./tests/loaderConfig/")

	tests.MaybeFail("test_create_config_with_LoaderConfig_instance", err,
		tests.Expect(errConf, nil),
		tests.Expect(c.LDRGetGOENV(), "localhost"),
		tests.Expect(len(c.LDRGetGrpcTypes()), 2),
		tests.Expect(len(c.LDRGetConfigsFiles()), 1),
		tests.Expect(c.LDRGetConfigsFiles()[0], "configs"),
		tests.Expect(c.CFGgoenv(), "localhost"),
	)
}

func Test_create_config_with_invalid_instance(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	c, err := NewConfig("localhost", "TInvalidConfig")

	tests.MaybeFail("test_create_config_with_invalid_instance",
		tests.Expect(c == nil, true),
		tests.Expect(err.Error(), "unknown ConfigType"))
}
