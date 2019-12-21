package gutil

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"net"
	"strconv"
)

type Consul struct {
	Client 	*api.Client
}

type ServiceInfo struct {
	Host      string
	Port      int
}

type RegisterInfo struct {
	ServiceInfo					   ServiceInfo
	ServiceName                    string
	Timeout                        string
	Interval                       string
	Deregister					   string
}

// 连接Consul
func NewConsul(Scheme string, Address string) *Consul {
	consul := &Consul{}
	client, err := api.NewClient(&api.Config{
		Scheme:Scheme, Address:Address,
	})
	if (err != nil) {
		fmt.Println("consul client error : ", err.Error())
		panic(err)
	}
	consul.Client = client

	return consul
}

// 注册服务
func (this *Consul) Register(register *RegisterInfo) error {
	if len(register.ServiceName) == 0 {
		return errors.New("must need Service Name")
	}

	if len(register.ServiceInfo.Host) == 0 {
		register.ServiceInfo.Host = LocalIP()
	}

	if register.ServiceInfo.Port == 0 {
		return errors.New("must need Service Port")
	}

	if len(register.Timeout) == 0 {
		register.Timeout = "1s"
	}

	if len(register.Interval) == 0 {
		register.Interval = "10s"
	}

	if len(register.Deregister) == 0 {
		register.Deregister = "30s"
	}

	service := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", register.ServiceName, register.ServiceInfo.Host, register.ServiceInfo.Port),
		Name:    register.ServiceName,
		Address: register.ServiceInfo.Host,
		Port:    register.ServiceInfo.Port,
		Tags:    []string{register.ServiceName},
		Check: &api.AgentServiceCheck{
			TCP: net.JoinHostPort(register.ServiceInfo.Host, strconv.Itoa(register.ServiceInfo.Port)),
			Interval: register.Interval,
			Timeout:  register.Timeout,
			DeregisterCriticalServiceAfter: register.Deregister,
		},
	}

	if err := this.Client.Agent().ServiceRegister(service); err != nil {
		fmt.Println("register service for consul failted")
		return err
	}

	fmt.Println("register service for consul success")
	return nil
}

// 服务发现
func (this *Consul) Discover(service_name string) ([]*api.AgentService, error) {
	services, _, err := this.Client.Health().Service(service_name, "", true,
		&api.QueryOptions{})

	if err != nil { return nil, err }

	var agentService []*api.AgentService
	for _, entry := range services {
		agentService = append(agentService, entry.Service)
	}

	return agentService, nil
}