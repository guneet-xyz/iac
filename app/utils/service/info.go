package service

import (
	"iac/utils/docker"
	"iac/utils/docker/compose"
	"log/slog"
	"sync"
	"time"
)

type ServiceInfoContainerRunningInfo struct {
	ContainerId   string
	ContainerName string
}

type ServiceInfoContainerConfigInfo struct {
	ContainerName string
	Name          string
}

type ServiceInfoContainerInspectInfo struct {
	RestartCount int
	Running      bool
	StartedAt    time.Time
}

type FetchingStatus string

const (
	FetchingStatusFetching FetchingStatus = "fetching"
	FetchingStatusFound    FetchingStatus = "fetched"
	FetchingStatusNotFound FetchingStatus = "not_found"
)

type ServiceInfoContainer struct {
	ConfigInfoFound   bool
	ConfigInfo        ServiceInfoContainerConfigInfo
	RunningInfoStatus FetchingStatus
	RunningInfo       ServiceInfoContainerRunningInfo
	InspectInfoStatus FetchingStatus
	InspectInfo       ServiceInfoContainerInspectInfo
}

type ServiceInfo struct {
	Name       string
	Containers []ServiceInfoContainer
}

type InspectResultChannelStruct struct {
	ContainerIndex int
	InspectResult  docker.InspectResult
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
			RunningInfoStatus: FetchingStatusFetching,
			RunningInfo:       ServiceInfoContainerRunningInfo{},
			InspectInfoStatus: FetchingStatusFetching,
			InspectInfo:       ServiceInfoContainerInspectInfo{},
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
				info.Containers[i].RunningInfoStatus = FetchingStatusFound
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
				RunningInfoStatus: FetchingStatusNotFound,
				RunningInfo: ServiceInfoContainerRunningInfo{
					ContainerId:   statsContainer.Id,
					ContainerName: statsContainer.Name,
				},
				InspectInfoStatus: FetchingStatusFetching,
				InspectInfo:       ServiceInfoContainerInspectInfo{},
			}
			info.Containers = append(info.Containers, container)
		}
	}
	svcInfoChannel <- info

	var wg sync.WaitGroup
	inspectResultChannel := make(chan InspectResultChannelStruct)

	for i, infoContainer := range info.Containers {
		if infoContainer.RunningInfoStatus == FetchingStatusFound {
			wg.Go(func() {
				inspectResults, err := docker.Inspect(infoContainer.RunningInfo.ContainerId)
				if err != nil {
					slog.Error("Failed to inspect container", "container", infoContainer.RunningInfo.ContainerId, "error", err)
					return
				}
				if len(inspectResults) > 0 {
					payload := InspectResultChannelStruct{
						ContainerIndex: i,
						InspectResult:  inspectResults[0],
					}
					inspectResultChannel <- payload
				}
			})
		} else {
			info.Containers[i].InspectInfoStatus = FetchingStatusNotFound
		}
	}

	go func() {
		wg.Wait()
		close(inspectResultChannel)
	}()

	for payload := range inspectResultChannel {
		inspectInfo := ServiceInfoContainerInspectInfo{
			RestartCount: payload.InspectResult.RestartCount,
			Running:      payload.InspectResult.State.Running,
			StartedAt:    payload.InspectResult.State.StartedAt,
		}
		info.Containers[payload.ContainerIndex].InspectInfoStatus = FetchingStatusFound
		info.Containers[payload.ContainerIndex].InspectInfo = inspectInfo
		svcInfoChannel <- info
	}

	return nil
}
