package configloader

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strconv"
	// "time"
)

const (
	localhost            string = "localhost"
	serviceUrl                  = "service_url"
	servicePort                 = "service_port"
	loaderconfigconf            = "loaderconfigconf"
	grpctypes                   = "grpctypes"
	configfiles                 = "configfiles"
	version                     = "version"
	storagePathConfig           = "storage_path_config"
	downloadPathConfig          = "download_path_config"
	storagePathGrpc             = "storage_path_grpc"
	downloadPathGrpc            = "download_path_grpc"
	delayBetweenReqRetry        = "delay_between_req_retry"
	reqRetry                    = "req_retry"
)

type loaderConfig struct {
	GOENV string
	*serviceEndpoint
	*loaderConfigConf
	grpcTypes   []string
	configFiles []string
}

func NewLoaderConfig(env string) *loaderConfig {
	return &loaderConfig{
		GOENV: env,
	}
}

func (lc *loaderConfig) LDRLoadConfigs(configName, configExtenssion, path string) error {

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	configFilePath := filepath.Join(currentDir, path)

	viper.SetConfigName(configName) // Config file name without extension
	viper.SetConfigType(configExtenssion)
	viper.AddConfigPath(configFilePath)

	// Read the configuration file
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	grpcTypes := []string{}
	configFiles := []string{}

	// set the registryConfigs
	se := NewServiceEndpoint()
	registryVars := viper.Sub(lc.LDRGetGOENV())
	if registryVars != nil {
		for key, value := range registryVars.AllSettings() {
			switch key {
			case serviceUrl:
				err := se.setServiceEndpointURL(value.(string))
				if err != nil {
					return err
				}
			case servicePort:
				err := se.setServiceEndpointPort(value.(string))
				if err != nil {
					return err
				}
			}
		}
	}
	lc.serviceEndpoint = se

	// set the loaderConfigConf
	lcc := NewLoaderConfigConf()
	lccDatas := viper.Sub(loaderconfigconf)
	if lccDatas != nil {
		for key, value := range lccDatas.AllSettings() {
			switch key {
			case version:
				lcc.setVersion(value.(string))
			case storagePathConfig:
				lcc.setStoragePathConfigs(value.(string))
			case downloadPathConfig:
				lcc.setDownloadPathConfigs(value.(string))
			case storagePathGrpc:
				lcc.setStoragePathGrpc(value.(string))
			case downloadPathGrpc:
				lcc.setDownloadPathGrpc(value.(string))
			case delayBetweenReqRetry:
				lcc.setDelayBetweenReqRetry(value.(string))
			case reqRetry:
				lcc.setReqRetry(value.(string))
			}

		}
	}
	lc.loaderConfigConf = lcc

	// set the dataConfigs
	grpcTypesVars := viper.Sub(grpctypes)
	if grpcTypesVars != nil {
		for _, value := range grpcTypesVars.AllSettings() {
			grpcTypes = append(grpcTypes, value.(string))
		}
	}
	lc.grpcTypes = grpcTypes

	configFilesVars := viper.Sub(configfiles)
	if configFilesVars != nil {
		for _, value := range configFilesVars.AllSettings() {
			configFiles = append(configFiles, value.(string))
		}
	}
	lc.configFiles = configFiles

	return nil
}

func (lc *loaderConfig) LDRGetGOENV() string {
	return lc.GOENV
}

func (lc *loaderConfig) LDRGetGrpcTypes() []string {
	return lc.grpcTypes
}

func (lc *loaderConfig) LDRGetConfigsFiles() []string {
	return lc.configFiles
}

// LoaderConfigs holds configurations to load datas
type loaderConfigConf struct {
	version              string
	storagePathConfigs   string
	downloadPathConfigs  string
	storagePathGrpc      string
	downloadPathGrpc     string
	delayBetweenReqRetry int8
	reqRetry             int8
}

func NewLoaderConfigConf() *loaderConfigConf {
	return &loaderConfigConf{}
}

func (l *loaderConfigConf) setVersion(v string) {
	if v != "" {
		l.version = v
	}
}
func (l *loaderConfigConf) LcfgGetVersion() string {
	return l.version
}

func (l *loaderConfigConf) setStoragePathConfigs(v string) {
	if v != "" {
		l.storagePathConfigs = v
	}
}
func (l *loaderConfigConf) LcfgGetStoragePathConfigs() string {
	return l.storagePathConfigs
}

func (l *loaderConfigConf) setDownloadPathConfigs(v string) {
	if v != "" {
		l.downloadPathConfigs = v
	}
}
func (l *loaderConfigConf) LcfgGetDownloadPathConfigs() string {
	return l.downloadPathConfigs
}

func (l *loaderConfigConf) setStoragePathGrpc(v string) {
	if v != "" {
		l.storagePathGrpc = v
	}
}
func (l *loaderConfigConf) LcfgGetStoragePathGrpc() string {
	return l.storagePathGrpc
}

func (l *loaderConfigConf) setDownloadPathGrpc(v string) {
	if v != "" {
		l.downloadPathGrpc = v
	}
}
func (l *loaderConfigConf) LcfgGetDownloadPathGrpc() string {
	return l.downloadPathGrpc
}

func (l *loaderConfigConf) setDelayBetweenReqRetry(v string) {
	if v != "" {
		num, err := strconv.ParseInt(v, 10, 8)
		if err == nil {
			l.delayBetweenReqRetry = int8(num)
		}
	}
}

func (l *loaderConfigConf) LcfgGetDelayBetweenReqRetry() int8 {
	return l.delayBetweenReqRetry
}

func (l *loaderConfigConf) setReqRetry(v string) {
	if v != "" {
		num, err := strconv.ParseInt(v, 10, 8)
		if err == nil {
			l.reqRetry = int8(num)
		}
	}
}
func (l *loaderConfigConf) LcfgGetReqRetry() int8 {
	return l.reqRetry
}

// ServiceEndpoint holds the serviceEndpoint informations
type serviceEndpoint struct {
	name string
	url  string
	port string
}

func NewServiceEndpoint() *serviceEndpoint {
	return &serviceEndpoint{}
}

func (se *serviceEndpoint) setServiceEndpointURL(url string) error {
	if len(url) > 0 {
		se.url = url
		return nil
	}
	return fmt.Errorf("registryURL is invalid")
}

func (se *serviceEndpoint) GetServiceEndpointURL() string {
	return se.url
}

func (se *serviceEndpoint) setServiceEndpointPort(port string) error {
	if len(port) > 0 {
		se.port = port
		return nil
	}
	return fmt.Errorf("registryPort is invalid")
}

func (se *serviceEndpoint) GetServiceEndpointPort() string {
	return se.port
}

func (se *serviceEndpoint) GetServiceEndpointFormatedURL() string {
	return fmt.Sprintf("%s:%s", se.url, se.port)
}
