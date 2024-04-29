package mongo

import (
	"context"
	"fmt"
	e "gitlab.com/grpasr/common/errors/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	connectionTimeoutDefault = 10
	requestTimeoutDefault    = 30
	totalConnectionDuration  = 300 // 10 * 30
)

type IConfigStrategy interface {
	getDSN() string
	getUsername() string
	getPassword() string
}

type StoreConfig struct {
	goEnv             string
	databaseName      string // each service have its own database
	connectionTimeout int
	requestTimeout    int
	isReplicaSet      bool
	configStrategy    IConfigStrategy
}

func (sc *StoreConfig) GetDatabaseName() string {
	return sc.databaseName
}

func (sc *StoreConfig) GetConnectionTimeout() int {
	return sc.connectionTimeout
}

func (sc *StoreConfig) GetRequestTimeout() int {
	return sc.requestTimeout
}

type NonReplicaSetConfig struct {
	url      string
	username string
	password string
}

func NewNonReplicaSetConfig(url, username, password string) NonReplicaSetConfig {
	return NonReplicaSetConfig{url, username, password}
}

type ReplicaSetConfig struct {
	url            string
	replicaSetName string
}

func NewReplicaSetConfig(url, replicaSetName string) ReplicaSetConfig {
	return ReplicaSetConfig{url, replicaSetName}
}

type NonReplicaSetStrategy struct {
	nonReplicaSetConfig NonReplicaSetConfig
}

func (nrs *NonReplicaSetStrategy) getDSN() string {
	return nrs.nonReplicaSetConfig.url
}

func (nrs *NonReplicaSetStrategy) getUsername() string {
	return nrs.nonReplicaSetConfig.username
}
func (nrs *NonReplicaSetStrategy) getPassword() string {
	return nrs.nonReplicaSetConfig.password
}

type ReplicaSetStrategy struct {
	replicaSetConfig ReplicaSetConfig
}

func (rs *ReplicaSetStrategy) getDSN() string {
	return rs.replicaSetConfig.url
}

func (nrs *ReplicaSetStrategy) getUsername() string { return "" }
func (nrs *ReplicaSetStrategy) getPassword() string { return "" }

func NewStoreConfig(goEnv, databaseName string, connectionTimeout, requestTimeout int, isReplicaSet bool, replicaSetCfg ReplicaSetConfig, nonReplicaSetCfg NonReplicaSetConfig) *StoreConfig {
	sc := &StoreConfig{
		goEnv:             goEnv,
		databaseName:      databaseName,
		isReplicaSet:      isReplicaSet,
		connectionTimeout: connectionTimeoutDefault,
		requestTimeout:    requestTimeoutDefault,
	}

	if connectionTimeout > 0 {
		sc.connectionTimeout = connectionTimeout
	}

	if requestTimeout > 0 {
		sc.requestTimeout = requestTimeout
	}

	if isReplicaSet {
		sc.configStrategy = &ReplicaSetStrategy{replicaSetConfig: replicaSetCfg}
	} else {
		sc.configStrategy = &NonReplicaSetStrategy{nonReplicaSetConfig: nonReplicaSetCfg}
	}
	return sc
}

// connection staff
func ClientConnect(client **mongo.Client, options *options.ClientOptions, storeCfg *StoreConfig) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(storeCfg.connectionTimeout)*time.Second)
	defer cancel()

	*client, err = mongo.Connect(ctx, options)
	return
}

func ClientPing(client **mongo.Client, options *options.ClientOptions, storeCfg *StoreConfig) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(storeCfg.requestTimeout)*time.Second)
	defer cancel()

	err = (*client).Ping(ctx, nil)
	if err != nil {
		return
	}
	return
}

type Runner func(**mongo.Client, *options.ClientOptions, *StoreConfig) error

func Retry(run Runner, retry, d int8, ctx context.Context) Runner {
	delay := time.Duration(int64(d))
	baseDelay := delay

	return func(client **mongo.Client, options *options.ClientOptions, storeCfg *StoreConfig) error {
		for r := int8(0); ; r++ {
			err := run(client, options, storeCfg)
			if err == nil || r > retry {
				return err
			}

			delay = time.Duration(baseDelay*time.Duration(r+1)) * time.Second
			fmt.Printf("Attempt %d failed; retrying in %v", r+1, delay)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func ClientProvider(storeCfg *StoreConfig, retry, d int8) (*mongo.Client, e.IError) {

	dsn := storeCfg.configStrategy.getDSN()

	clientOptions := options.Client().ApplyURI(dsn)

	if !storeCfg.isReplicaSet {
		clientOptions.SetAuth(options.Credential{
			Username: storeCfg.configStrategy.getUsername(),
			Password: storeCfg.configStrategy.getPassword(),
		})
	}

	var client *mongo.Client

	ctx, cancel := context.WithTimeout(
		context.Background(), totalConnectionDuration*time.Second)
	defer cancel()

	retryConnect := Retry(ClientConnect, retry, d, ctx)
	err := retryConnect(&client, clientOptions, storeCfg)
	if err != nil {
		// obs.Logging.NewLogHandler(obs.Logging.LLHError()).
		// 	Err(err).
		// 	Msg("error creating mongoDB client")
		return nil, e.NewCustomHTTPStatus(e.StatusBadRequest, "", err.Error())
	}

	retryPing := Retry(ClientPing, retry, d, ctx)
	err = retryPing(&client, nil, storeCfg)
	if err != nil {
		// obs.Logging.NewLogHandler(obs.Logging.LLHError()).
		// 	Err(err).
		// 	Msg("error ping mongoDB client")
		return nil, e.NewCustomHTTPStatus(e.StatusBadRequest, "", err.Error())
	}

	return client, nil
}
