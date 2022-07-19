package service

// variables refering to the redis exporter port
const (
	exporterPort                  = 9121
	sentinelExporterPort          = 9355
	exporterPortName              = "http-metrics"
	exporterContainerName         = "redis-exporter"
	sentinelExporterContainerName = "sentinel-exporter"
	exporterDefaultRequestCPU     = "25m"
	exporterDefaultLimitCPU       = "50m"
	exporterDefaultRequestMemory  = "50Mi"
	exporterDefaultLimitMemory    = "100Mi"
)

const (
	baseName = "rf"
	//sentinelName           = "s"
	sentinelName           = "sentinel"
	sentinelRoleName       = "sentinel"
	sentinelConfigFileName = "sentinel.conf"
	redisConfigFileName    = "redis.conf"
	//redisName              = "r"
	redisName         = "redis"
	redisMasterName   = "redis-master"
	redisExporterName = "redis-exporter"
	//redisShutdownName   = "r-s"
	redisShutdownName = "redis-shutdown"
	//redisReadinessName  = "r-readiness"
	redisReadinessName  = "redis-readiness"
	redisRoleName       = "redis"
	appLabel            = "redis-failover"
	hostnameTopologyKey = "kubernetes.io/hostname"
)

const (
	redisRoleLabelKey    = "redisfailovers-role"
	redisRoleLabelMaster = "master"
	redisRoleLabelSlave  = "slave"
)
