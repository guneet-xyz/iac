package service

import (
	"iac/utils/docker/compose"
	"log/slog"
)

type ServiceInfoContainerRunningInfo struct {
	ContainerId   string
	ContainerName string
}

type ServiceInfoContainerConfigInfo struct {
	ContainerName string
	Name          string
}

type ServiceInfoContainer struct {
	ConfigInfoFound  bool
	ConfigInfo       ServiceInfoContainerConfigInfo
	RunningInfoFound bool
	RunningInfo      ServiceInfoContainerRunningInfo
}

type ServiceInfo struct {
	Name       string
	Containers []ServiceInfoContainer
}

func GetServiceInfo(svcName string) (ServiceInfo, error) {
	slog.Debug("Getting service info", "service", svcName)
	configPath, err := GetComposePathFromServiceName(svcName)
	if err != nil {
		return ServiceInfo{}, err
	}

	info := ServiceInfo{
		Name:       svcName,
		Containers: []ServiceInfoContainer{},
	}

	config, err := compose.ReadConfig(configPath)

	for containerSvcName, containerSvc := range config.Services {
		container := ServiceInfoContainer{
			ConfigInfoFound: true,
			ConfigInfo: ServiceInfoContainerConfigInfo{
				ContainerName: containerSvc.ContainerName,
				Name:          containerSvcName,
			},
			RunningInfoFound: false,
			RunningInfo:      ServiceInfoContainerRunningInfo{},
		}
		info.Containers = append(info.Containers, container)
	}

	stats, err := compose.GetStats(configPath)
	if err != nil {
		return ServiceInfo{}, err
	}

	for _, statsContainer := range stats.Containers {
		found := false
		for i, infoContainer := range info.Containers {
			if statsContainer.Name == infoContainer.ConfigInfo.ContainerName {
				info.Containers[i].RunningInfoFound = true
				info.Containers[i].RunningInfo = ServiceInfoContainerRunningInfo{
					ContainerId:   statsContainer.Id,
					ContainerName: statsContainer.Name,
				}
				found = true
				break
			}
		}
		if !found {
			container := ServiceInfoContainer{
				ConfigInfoFound:  false,
				ConfigInfo:       ServiceInfoContainerConfigInfo{},
				RunningInfoFound: true,
				RunningInfo: ServiceInfoContainerRunningInfo{
					ContainerId:   statsContainer.Id,
					ContainerName: statsContainer.Name,
				},
			}
			info.Containers = append(info.Containers, container)
		}
	}

	return info, nil
}
