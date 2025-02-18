package redisfailover

import (
	"context"
	"github.com/spotahome/kooper/v2/controller/leaderelection"
	"k8s.io/client-go/kubernetes"
	"time"

	"github.com/spotahome/kooper/v2/controller"
	kooperlog "github.com/spotahome/kooper/v2/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/spotahome/redis-operator/log"
	"github.com/spotahome/redis-operator/metrics"
	rfservice "github.com/spotahome/redis-operator/operator/redisfailover/service"
	"github.com/spotahome/redis-operator/service/k8s"
	"github.com/spotahome/redis-operator/service/redis"
)

const (
	//resync       = 30 * time.Second
	resync       = 10 * time.Second
	operatorName = "redis-operator"

	leaderElectionKey = "leader-election-redis-operator"
)

// New will create an operator that is responsible of managing all the required stuff
// to create redis failovers.
func New(cfg Config, k8sService k8s.Services, redisClient redis.Client, kooperMetricsRecorder metrics.Recorder,
	logger log.Logger, kubecli kubernetes.Interface) (controller.Controller, error) {
	// Create internal services.
	rfService := rfservice.NewRedisFailoverKubeClient(k8sService, logger)
	rfChecker := rfservice.NewRedisFailoverChecker(k8sService, redisClient, logger)
	rfHealer := rfservice.NewRedisFailoverHealer(k8sService, redisClient, logger)

	// Create the handlers.
	rfHandler := NewRedisFailoverHandler(cfg, rfService, rfChecker, rfHealer, k8sService, kooperMetricsRecorder, logger)
	rfRetriever := NewRedisFailoverRetriever(k8sService)

	kooperlogger := kooperlogger{Logger: logger.WithField("operator", "redisfailover")}

	controllerCfg := controller.Config{
		Handler:         rfHandler,
		Retriever:       rfRetriever,
		MetricsRecorder: kooperMetricsRecorder,
		Logger:          kooperlogger,
		Name:            "redisfailover",
		ResyncInterval:  resync,
	}

	if cfg.EnableLeaderElection {
		// Leader election service.
		lesvc, err := leaderelection.NewDefault(leaderElectionKey, cfg.OperatorNameSpace, kubecli, kooperlogger)
		if err != nil {
			return nil, err
		}
		controllerCfg.LeaderElector = lesvc
	}

	// Create our controller.
	return controller.New(&controllerCfg)
}

func NewRedisFailoverRetriever(cli k8s.Services) controller.Retriever {
	return controller.MustRetrieverFromListerWatcher(&cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return cli.ListRedisFailovers(context.Background(), "", options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return cli.WatchRedisFailovers(context.Background(), "", options)
		},
	})
}

type kooperlogger struct {
	log.Logger
}

func (k kooperlogger) WithKV(kv kooperlog.KV) kooperlog.Logger {
	return kooperlogger{Logger: k.Logger.WithFields(kv)}
}
