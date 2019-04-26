package gutil

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"time"
)

type Consul struct {
	Config	*api.Config
	Client 	*api.Client
}

type ServiceInfo struct {
	Host      string
	Port      int
}

type RegisterInfo struct {
	ServiceInfo					   ServiceInfo
	ServiceName                    string
	Timeout                        time.Duration
	Interval                       time.Duration
	DeregisterCriticalServiceAfter time.Duration
}

// 连接Consul
func (this *Consul) Connect()  {
	if (this.Client == nil) {
		config := api.DefaultConfig()
		config.Scheme = this.Config.Scheme
		config.Address = this.Config.Address
		client, err := api.NewClient(config)
		if (err != nil) {
			fmt.Println("consul client error : ", err.Error())
			panic(err)
		}

		this.Client = client
	}
}

// 注册服务
func (this *Consul) Register(register *RegisterInfo) error {
	if len(register.ServiceInfo.Host) > 0 {
		register.ServiceInfo.Host = LocalIP()
	}

	if len(register.ServiceName) > 0 {
		return errors.New("must need Service Name")
	}

	if (register.Timeout.String() == "") {
		register.Timeout = 3
	}

	if (register.Interval.String() == "") {
		register.Interval = 5
	}

	if (register.DeregisterCriticalServiceAfter.String() == "") {
		register.DeregisterCriticalServiceAfter = 30
	}

	fmt.Println(this.Client)

	service := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", register.ServiceName, register.ServiceInfo.Host, register.ServiceInfo.Port),
		Name:    register.ServiceName,
		Port:    register.ServiceInfo.Port,
		Address: register.ServiceInfo.Host,
		Tags:    []string{register.ServiceName},
		Check: &api.AgentServiceCheck{
			TCP:     fmt.Sprintf("tcp://%s-%d", register.ServiceInfo.Host, register.ServiceInfo.Port),
			Interval: string(register.Interval * time.Second),
			Timeout:  string(register.Timeout * time.Second),
			DeregisterCriticalServiceAfter: string(register.DeregisterCriticalServiceAfter * time.Second),
		},
	}

	if err := this.Client.Agent().ServiceRegister(service); err != nil {
		return err
	}

	return nil
}

// 服务发现
func (this *Consul) Discover(service_name string) (agentService []*api.AgentService, err error) {
	var AgentService []*api.AgentService

	services, _, err := this.Client.Catalog().Services(&api.QueryOptions{})
	if err != nil {
		return AgentService, err
	}

	for name := range services {
		servicesData, _, err := this.Client.Health().Service(name, "", true,
			&api.QueryOptions{})
		if err != nil {
			return AgentService, err
		}

		for _, entry := range servicesData {
			if service_name != entry.Service.Service {
				continue
			}

			for _, health := range entry.Checks {
				if service_name != health.ServiceName {
					continue
				}
				AgentService = append(AgentService, entry.Service)
			}
		}
	}

	return AgentService, nil
}