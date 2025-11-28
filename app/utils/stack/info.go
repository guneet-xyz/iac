package stack

import (
	"iac/utils/docker"
	"iac/utils/docker/compose"
	"log/slog"
	"sync"
	"time"
)

type StackInfoContainerRunningInfo struct {
	ContainerId   string
	ContainerName string
}

type StackInfoContainerConfigInfo struct {
	ContainerName string
	Name          string
}

type StackInfoContainerInspectInfo struct {
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

type StackInfoContainer struct {
	ConfigInfoFound   bool
	ConfigInfo        StackInfoContainerConfigInfo
	RunningInfoStatus FetchingStatus
	RunningInfo       StackInfoContainerRunningInfo
	InspectInfoStatus FetchingStatus
	InspectInfo       StackInfoContainerInspectInfo
}

type StackInfo struct {
	Name       string
	Containers []StackInfoContainer
}

type InspectResultChannelStruct struct {
	ContainerIndex int
	InspectResult  docker.InspectResult
}

func SendStackInfoToChannel(stackName string, stackInfoChannel chan StackInfo) error {
	slog.Debug("Getting stack info", "stackName", stackName)
	configPath, err := GetComposePathFromStackName(stackName)
	if err != nil {
		return err
	}

	info := StackInfo{
		Name:       stackName,
		Containers: []StackInfoContainer{},
	}

	config, err := compose.ReadConfig(configPath)
	if err != nil {
		return err
	}

	for containerSvcName, containerSvc := range config.Services {
		container := StackInfoContainer{
			ConfigInfoFound: true,
			ConfigInfo: StackInfoContainerConfigInfo{
				ContainerName: containerSvc.ContainerName,
				Name:          containerSvcName,
			},
			RunningInfoStatus: FetchingStatusFetching,
			RunningInfo:       StackInfoContainerRunningInfo{},
			InspectInfoStatus: FetchingStatusFetching,
			InspectInfo:       StackInfoContainerInspectInfo{},
		}
		info.Containers = append(info.Containers, container)
	}

	stackInfoChannel <- info

	stats, err := compose.GetStats(configPath)
	if err != nil {
		return err
	}

	for _, statsContainer := range stats.Containers {
		found := false
		for i, infoContainer := range info.Containers {
			if statsContainer.Name == infoContainer.ConfigInfo.ContainerName {
				info.Containers[i].RunningInfoStatus = FetchingStatusFound
				info.Containers[i].RunningInfo = StackInfoContainerRunningInfo{
					ContainerId:   statsContainer.Id,
					ContainerName: statsContainer.Name,
				}
				found = true
				break
			}
		}
		if !found {
			container := StackInfoContainer{
				ConfigInfoFound:   false,
				ConfigInfo:        StackInfoContainerConfigInfo{},
				RunningInfoStatus: FetchingStatusNotFound,
				RunningInfo: StackInfoContainerRunningInfo{
					ContainerId:   statsContainer.Id,
					ContainerName: statsContainer.Name,
				},
				InspectInfoStatus: FetchingStatusFetching,
				InspectInfo:       StackInfoContainerInspectInfo{},
			}
			info.Containers = append(info.Containers, container)
		}
	}
	stackInfoChannel <- info

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
		inspectInfo := StackInfoContainerInspectInfo{
			RestartCount: payload.InspectResult.RestartCount,
			Running:      payload.InspectResult.State.Running,
			StartedAt:    payload.InspectResult.State.StartedAt,
		}
		info.Containers[payload.ContainerIndex].InspectInfoStatus = FetchingStatusFound
		info.Containers[payload.ContainerIndex].InspectInfo = inspectInfo
		stackInfoChannel <- info
	}

	return nil
}
