package gutil

import (
	"fmt"
	"github.com/bincent/gutil"
	"github.com/hashicorp/consul/api"
	"time"
)

// ConsulRegister consul service register
type ConsulRegister struct {
	Address                        string
	Service                        string
	Tag                            []string
	Port                           int
	CheckUri					   string
	Timeout                        time.Duration
	Interval                       time.Duration
	DeregisterCriticalServiceAfter time.Duration
}

// Register register service
func (this *ConsulRegister) Register() error {
	config := api.DefaultConfig()
	config.Address = this.Address
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	agent := client.Agent()

	IP := gutil.LocalIP()

	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", this.Service, IP, this.Port), // 服务节点的名称
		Name:    fmt.Sprintf("grpc.health.v1.%v", this.Service),    // 服务名称
		Tags:    this.Tag,                                          // tag，可以为空
		Port:    this.Port,                                         // 服务端口
		Address: IP,                                             // 服务 IP
		Check: &api.AgentServiceCheck{ // 健康检查
			HTTP:                           this.CheckUri,								  // 健康检测地址
			Timeout:                        this.Timeout.String(),						  // 超时
			Interval: 						this.Interval.String(),                       // 健康检查间隔
			DeregisterCriticalServiceAfter: this.DeregisterCriticalServiceAfter.String(), // 注销时间，相当于过期时间
		},
	}

	if err := agent.ServiceRegister(reg); err != nil {
		return err
	}

	return nil
}