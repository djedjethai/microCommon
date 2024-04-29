package observability

import (
	"gitlab.com/grpasr/common/observability/logging"
	"gitlab.com/grpasr/common/observability/metrics"
	"gitlab.com/grpasr/common/observability/tracing"
)

var (
	Tracing *tracing.TracingFacade
	Metrics *metrics.MetricsFacade
	Logging *logging.LoggingFacade
)

func SetObservabilityFacade(svcName ...string) {
	Tracing = tracing.NewTracingfacade()
	Metrics = metrics.NewMetricsFacade()
	Logging = logging.NewLoggingFacade(svcName...)
}

// // GetTls returns a configuration that enables the use of mutual TLS.
// func GetTls() (*tls.Config, error) {
// 	clientAuth, err := tls.LoadX509KeyPair("./confs/client.crt", "./confs/client.key")
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	caCert, err := os.ReadFile("./confs/rootCA.crt")
// 	if err != nil {
// 		return nil, err
// 	}
// 	caCertPool := x509.NewCertPool()
// 	caCertPool.AppendCertsFromPEM(caCert)
//
// 	c := &tls.Config{
// 		RootCAs:      caCertPool,
// 		Certificates: []tls.Certificate{clientAuth},
// 	}
//
// 	return c, nil
// }
