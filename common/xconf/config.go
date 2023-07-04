package xconf

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
)

type LogConfig struct {
	// Mode represents the logging mode, default is `console`.
	// console: log to console.
	// file: log to file.
	// volume: used in k8s, prepend the hostname to the log file name.
	Mode string `json:",default=console,options=[console,file,volume]"`
	// Encoding represents the encoding type, default is `json`.
	// json: json encoding.
	// plain: plain text encoding, typically used in development.
	Encoding string `json:",default=json,options=[json,plain]"`
	// TimeFormat represents the time format, default is `2006-01-02T15:04:05.000Z07:00`.
	TimeFormat string `json:",optional"`
	// Path represents the log file path, default is `logs`.
	Path string `json:",default=logs"`
	// Level represents the log level, default is `info`.
	Level string `json:",default=info,options=[debug,info,error,severe]"`
	// MaxContentLength represents the max content bytes, default is no limit.
	MaxContentLength uint32 `json:",optional"`
	// Compress represents whether to compress the log file, default is `false`.
	Compress bool `json:",optional"`
	// Stat represents whether to log statistics, default is `true`.
	Stat bool `json:",default=true"`
	// KeepDays represents how many days the log files will be kept. Default to keep all files.
	// Only take effect when Mode is `file` or `volume`, both work when Rotation is `daily` or `size`.
	KeepDays int `json:",optional"`
	// StackCooldownMillis represents the cooldown time for stack logging, default is 100ms.
	StackCooldownMillis int `json:",default=100"`
	// MaxBackups represents how many backup log files will be kept. 0 means all files will be kept forever.
	// Only take effect when RotationRuleType is `size`.
	// Even thougth `MaxBackups` sets 0, log files will still be removed
	// if the `KeepDays` limitation is reached.
	MaxBackups int `json:",default=0"`
	// MaxSize represents how much space the writing log file takes up. 0 means no limit. The unit is `MB`.
	// Only take effect when RotationRuleType is `size`
	MaxSize int `json:",default=0"`
	// Rotation represents the type of log rotation rule. Default is `daily`.
	// daily: daily rotation.
	// size: size limited rotation.
	Rotation string `json:",default=daily,options=[daily,size]"`
}

type EtcdConf struct {
	Hosts              []string
	KeyPrefix          string `json:",optional"`
	ID                 int64  `json:",optional"`
	User               string `json:",optional"`
	Pass               string `json:",optional"`
	CertFile           string `json:",optional"`
	CertKeyFile        string `json:",optional=CertFile"`
	CACertFile         string `json:",optional=CertFile"`
	InsecureSkipVerify bool   `json:",optional"`
}

type DiscoveryConfig struct {
	// Etcd represents the etcd configurations. If not set, then use k8s.
	Etcd EtcdConf `json:",optional"`
	// K8sNamespace represents the k8s namespace, If not set, then use localhost.
	// If you want to use k8s, You have to create service for all rpc services in k8s.
	// The service name should be `${serviceName}-svc`, and the port name suggested to be `rpc`.
	// example: user-rpc-svc, tokengating-rpc-svc
	// And k8s deployment should set serviceAccount, which must have find-endpoint permission.
	K8sNamespace string `json:",optional"`
}

type Config struct {
	Gateway   GatewayConfig
	App       AppConfig
	Log       LogConfig
	Discovery DiscoveryConfig
	Mode      string `json:",default=pro,options=dev|test|rt|pre|pro"`
}

func (c Config) LogConf(name string) logx.LogConf {
	return logx.LogConf{
		ServiceName:         name,
		Mode:                c.Log.Mode,
		Encoding:            c.Log.Encoding,
		TimeFormat:          c.Log.TimeFormat,
		Path:                c.Log.Path,
		Level:               c.Log.Level,
		MaxContentLength:    c.Log.MaxContentLength,
		Compress:            c.Log.Compress,
		Stat:                c.Log.Stat,
		KeepDays:            c.Log.KeepDays,
		StackCooldownMillis: c.Log.StackCooldownMillis,
		MaxBackups:          c.Log.MaxBackups,
		MaxSize:             c.Log.MaxSize,
		Rotation:            c.Log.Rotation,
	}
}

func (c Config) Etcd(name string) discov.EtcdConf {
	if len(c.Discovery.Etcd.Hosts) == 0 {
		return discov.EtcdConf{}
	}
	return discov.EtcdConf{
		Hosts:              c.Discovery.Etcd.Hosts,
		Key:                c.Discovery.Etcd.KeyPrefix + name,
		ID:                 c.Discovery.Etcd.ID,
		User:               c.Discovery.Etcd.User,
		Pass:               c.Discovery.Etcd.Pass,
		CertFile:           c.Discovery.Etcd.CertFile,
		CertKeyFile:        c.Discovery.Etcd.CertKeyFile,
		CACertFile:         c.Discovery.Etcd.CACertFile,
		InsecureSkipVerify: c.Discovery.Etcd.InsecureSkipVerify,
	}
}
