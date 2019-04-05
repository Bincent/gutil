package gutil

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"time"
)

type Consul struct {
	Client 	*api.Client
}

type ConsulDiscover struct {
	ServiceName                    string
	Balance                        string
}

type ConsulRegister struct {
	ServiceID					   string
	ServiceName                    string
	ServiceTag                     []string
	Port                           int
	MonitorAddress				   string
	Timeout                        time.Duration
	Interval                       time.Duration
	DeregisterCriticalServiceAfter time.Duration
}

// 注册服务
func (this *Consul) Register(register *ConsulRegister) error {
	ip := LocalIP()

	if len(register.ServiceName) > 0 {
		return errors.New("must need Service Name")
	}

	//if (register.Port == nil) {
	//	return errors.New("must need Service Port")
	//}

	if len(register.ServiceID) > 0 {
		register.ServiceID = register.ServiceName + "-" + ip
	}

	if len(register.MonitorAddress) > 0 {
		register.MonitorAddress = "tcp://" + ip
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

	service := &api.AgentServiceRegistration{
		ID:      register.ServiceID,
		Name:    register.ServiceName,
		Port:    register.Port,
		Address: ip,
		Tags:    register.ServiceTag,
		Check: &api.AgentServiceCheck{
			TCP:     register.MonitorAddress,
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
func (this *Consul) Discover(discover *ConsulDiscover) (agentService []*api.AgentService, err error) {
	var _agentService []*api.AgentService

	services, _, err := this.Client.Catalog().Services(&api.QueryOptions{})
	for name := range services {
		servicesData, _, err := this.Client.Health().Service(name, "", true,
			&api.QueryOptions{})
		if err != nil { return _agentService, err }

		for _, entry := range servicesData {
			if discover.ServiceName != entry.Service.Service {
				continue
			}

			for _, health := range entry.Checks {
				if health.ServiceName != discover.ServiceName {
					continue
				}

				_agentService = append(_agentService, entry.Service)
			}
		}
	}

	return _agentService, nil
}

func (this *Consul) SetClient(client *api.Client) {
	this.Client = client
}