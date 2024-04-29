package configloader

import (
	"gitlab.com/grpasr/common/tests"
	"testing"
)

func Test_load_config(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	lc := NewLoaderConfig("localhost")

	err := lc.LDRLoadConfigs("configTests", "yaml", "./tests/loaderConfig/")

	tests.MaybeFail("test_load_config", err,
		tests.Expect(lc.LDRGetGOENV(), "localhost"),
		tests.Expect(lc.GetServiceEndpointURL(), "127.0.0.1"),
		tests.Expect(lc.GetServiceEndpointPort(), "4000"),
		tests.Expect(len(lc.LDRGetGrpcTypes()), 2),
		tests.Expect(len(lc.LDRGetConfigsFiles()), 1),
		tests.Expect(lc.LDRGetConfigsFiles()[0], "configs"),
		tests.Expect(lc.GetServiceEndpointFormatedURL(), "127.0.0.1:4000"),
		tests.Expect(lc.LcfgGetVersion(), "1"), // add
		tests.Expect(lc.LcfgGetStoragePathConfigs(), "../broker/configs"),
		tests.Expect(lc.LcfgGetDownloadPathConfigs(), "/configs"),
		tests.Expect(lc.LcfgGetStoragePathGrpc(), "../broker/api"),
		tests.Expect(lc.LcfgGetDownloadPathGrpc(), "/grpc"),
		tests.Expect(lc.LcfgGetDelayBetweenReqRetry(), int8(5)),
		tests.Expect(lc.LcfgGetReqRetry(), int8(3)),
	)
}

func Test_load_config_production(t *testing.T) {
	tests.MaybeFail = tests.InitFailFunc(t)

	lc := NewLoaderConfig("production")

	err := lc.LDRLoadConfigs("configTests", "yaml", "./tests/loaderConfig/")

	tests.MaybeFail("test_load_config_production", err,
		tests.Expect(lc.LDRGetGOENV(), "production"),
		tests.Expect(lc.GetServiceEndpointURL(), "production-url"),
		tests.Expect(lc.GetServiceEndpointPort(), "4000"),
		tests.Expect(lc.GetServiceEndpointFormatedURL(), "production-url:4000"))

}

// TODO ........
// func Test_load_config_invalid_varenv(t *testing.T) {
// 	tests.MaybeFail = tests.InitFailFunc(t)
//
// 	lc := NewLoaderConfig("development")
//
// 	err := lc.LDRLoadConfigs("configTests", "yaml", "./tests/loaderConfig/")
//
// 	tests.MaybeFail("test_load_config", err,
// 		tests.Expect(lc.LDRGetGOENV(), "development"),
// 		tests.Expect(lc.GetServiceEndpointURL(), "127.0.0.1"),
// 		tests.Expect(lc.GetServiceEndpointPort(), "4000"),
// 		tests.Expect(len(lc.LDRGetGrpcTypes()), 2),
// 		tests.Expect(len(lc.LDRGetConfigsFiles()), 1),
// 		tests.Expect(lc.LDRGetConfigsFiles()[0], "configs"),
// 		tests.Expect(lc.GetServiceEndpointFormatedURL(), "127.0.0.1:4000"))
// }
//
// func Test_load_config_production_invalid_endpointURL(t *testing.T) {
// 	tests.MaybeFail = tests.InitFailFunc(t)
//
// 	lc := NewLoaderConfig("ggg")
//
// 	err := lc.LDRLoadConfigs("configTests", "yaml", "./tests/loaderConfig/")
//
// 	tests.MaybeFail("Test_load_config_production_invalid_endpointURL",
// 		tests.Expect(err.Error(), "registryURL is invalid"))
// }
//
// func Test_load_config_production_invalid_endpointPort(t *testing.T) {
// 	tests.MaybeFail = tests.InitFailFunc(t)
//
// 	lc := NewLoaderConfig("ggg")
//
// 	err := lc.LDRLoadConfigs("configTests", "yaml", "./tests/loaderConfig/")
//
// 	tests.MaybeFail("Test_load_config_production_invalid_endpointPort",
// 		tests.Expect(err.Error(), "registryPort is invalid"))
//
// }
