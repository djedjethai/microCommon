package configs

// import (
// 	"errors"
// )
//
// type TracingConfigs struct {
// 	jaegerEndpoint string
// 	service        string
// 	id             int32
// 	environment    string
// 	tracingLibrary string
// 	samplingRatio  float32
// 	level          string // verbose, debug, prod
// }
//
// func NewTracingConfigsDefault(svc, env, level string, splRatio float32, id int32) (*TracingConfigs, error) {
// 	tc := &TracingConfigs{
// 		jaegerEndpoint: "http://127.0.0.1:14268/api/traces",
// 		service:        svc,
// 		id:             id,
// 		environment:    env,
// 		tracingLibrary: "go.opentelemetry.io/otel/trace",
// 	}
//
// 	if level != "verbose" && level != "debug" && level != "prod" {
// 		return nil, errors.New("Jaeger, invalid level")
// 	} else {
// 		tc.level = level
// 	}
//
// 	if splRatio > 1 || splRatio < 0 {
// 		return nil, errors.New("Jaeger, invalid ration")
// 	} else {
// 		tc.samplingRatio = splRatio
// 	}
//
// 	return tc, nil
//
// }
//
// func NewTracingConfigs(svc, env, level, endpoint, tracingLibrary string, splRatio float32, id int32) (*TracingConfigs, error) {
// 	trcCfig, err := NewTracingConfigsDefault(svc, env, level, splRatio, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(endpoint) > 0 {
// 		trcCfig.jaegerEndpoint = endpoint
// 	}
//
// 	if len(tracingLibrary) > 0 {
// 		trcCfig.tracingLibrary = tracingLibrary
// 	}
//
// 	return trcCfig, nil
// }
//
// func (tc *TracingConfigs) GetService() string {
// 	return tc.service
// }
//
// func (tc *TracingConfigs) GetID() int32 {
// 	return tc.id
// }
//
// func (tc *TracingConfigs) GetLevel() string {
// 	return tc.level
// }
//
// func (tc *TracingConfigs) GetJaegerEndpoint() string {
// 	return tc.jaegerEndpoint
// }
//
// func (tc *TracingConfigs) GetEnvironment() string {
// 	return tc.environment
// }
