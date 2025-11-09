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

type RunningInfoStatus string

const (
	RunningInfoStatusFetching RunningInfoStatus = "fetching"
	RunningInfoStatusFound    RunningInfoStatus = "running"
	RunningInfoStatusNotFound RunningInfoStatus = "not_found"
)

type ServiceInfoContainer struct {
	ConfigInfoFound   bool
	ConfigInfo        ServiceInfoContainerConfigInfo
	RunningInfoStatus RunningInfoStatus
	RunningInfo       ServiceInfoContainerRunningInfo
}

type ServiceInfo struct {
	Name       string
	Containers []ServiceInfoContainer
}

func SendServiceInfoToChannel(svcName string, svcInfoChannel chan ServiceInfo) error {
	slog.Debug("Getting service info", "service", svcName)
	configPath, err := GetComposePathFromServiceName(svcName)
	if err != nil {
		return err
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
			RunningInfoStatus: RunningInfoStatusFetching,
			RunningInfo:       ServiceInfoContainerRunningInfo{},
		}
		info.Containers = append(info.Containers, container)
	}

	svcInfoChannel <- info

	stats, err := compose.GetStats(configPath)
	if err != nil {
		return err
	}

	for _, statsContainer := range stats.Containers {
		found := false
		for i, infoContainer := range info.Containers {
			if statsContainer.Name == infoContainer.ConfigInfo.ContainerName {
				info.Containers[i].RunningInfoStatus = RunningInfoStatusFound
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
				ConfigInfoFound:   false,
				ConfigInfo:        ServiceInfoContainerConfigInfo{},
				RunningInfoStatus: RunningInfoStatusNotFound,
				RunningInfo: ServiceInfoContainerRunningInfo{
					ContainerId:   statsContainer.Id,
					ContainerName: statsContainer.Name,
				},
			}
			info.Containers = append(info.Containers, container)
		}
	}

	svcInfoChannel <- info
	return nil
}
