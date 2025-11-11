package services

import (
	"fmt"
	"iac/utils/errors"
	"iac/utils/exitcodes"
	"iac/utils/out"
	"iac/utils/out/symbols"
	"iac/utils/service"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	name string
)

var Cmd = &cobra.Command{
	Use:   "services",
	Short: "info about services",
	RunE:  func(_ *cobra.Command, _ []string) error { return ShowInfoAboutService(name) },
}

func init() {
	Cmd.Flags().StringVar(&name, "name", "", "Specify a service to get info about")
}

func ShowInfoAboutService(name string) error {
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	svcInfoWgFinished := false
	svcInfos := []service.ServiceInfo{}
	svcInfoChannel := make(chan service.ServiceInfo)

	wg.Go(func() {

		services, err := service.GetServiceNames()
		if err != nil {
			return
		}

		var filteredServices []string

		if name != "" {
			for _, svc := range services {
				if svc == name {
					filteredServices = append(filteredServices, svc)
					break
				}
			}
			if len(filteredServices) == 0 {
				slog.Info("Service not found", "name", name)
				os.Exit(exitcodes.ServiceNotFound)
				return
			}
		} else {
			filteredServices = services
			if len(filteredServices) == 0 {
				slog.Info("No services found")
				os.Exit(exitcodes.NoServicesFound)
				return
			}
		}

		svcInfoWg := sync.WaitGroup{}

		for _, svcName := range filteredServices {
			svcInfoWg.Go(func() { service.SendServiceInfoToChannel(svcName, svcInfoChannel) })
		}

		svcInfoWg.Wait()
		svcInfoWgFinished = true
		close(svcInfoChannel)

	})

	wg.Go(func() {
		for svcInfo := range svcInfoChannel {
			found := false
			for i, existingSvcInfo := range svcInfos {
				if existingSvcInfo.Name == svcInfo.Name {
					svcInfos[i] = svcInfo
					found = true
					break
				}
			}
			if !found {
				svcInfos = append(svcInfos, svcInfo)
			}
			mutex.Lock()
			printServiceInfo(svcInfos)
			mutex.Unlock()
		}
	})

	wg.Go(func() {
		for {
			if svcInfoWgFinished {
				return
			}
			mutex.Lock()
			printServiceInfo(svcInfos)
			mutex.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	})

	wg.Wait()

	return nil
}

func printServiceInfo(svcInfos []service.ServiceInfo) {
	lines := []string{}
	if len(svcInfos) == 0 {
		lines = append(lines, out.SpinnerChar())
	}

	for _, info := range svcInfos {
		lines = append(lines, fmt.Sprintf("%s", info.Name))
		for _, container := range info.Containers {
			var name string
			if container.ConfigInfoFound {
				name = container.ConfigInfo.Name
			} else if container.RunningInfoStatus == service.FetchingStatusFound {
				name = container.RunningInfo.ContainerName
			} else {
				err := errors.New("This should never happen. Container config info not found, but running info also not found")
				panic(err)
			}

			var symbol string
			switch container.RunningInfoStatus {
			case service.FetchingStatusFetching:
				symbol = out.SpinnerChar()
			case service.FetchingStatusNotFound:
				symbol = symbols.X
			case service.FetchingStatusFound:
				switch container.ConfigInfoFound {
				case false:
					symbol = symbols.Question
				case true:
					symbol = symbols.Check
				}
			}

			var extraInfo string
			switch container.InspectInfoStatus {
			case service.FetchingStatusFetching:
				extraInfo = out.SpinnerChar()
			case service.FetchingStatusNotFound:
				symbol = symbols.X
				extraInfo = ""
			case service.FetchingStatusFound:
				if !container.InspectInfo.Running {
					symbol = symbols.X
					break
				}
				restarts := container.InspectInfo.RestartCount
				runningSince := time.Since(container.InspectInfo.StartedAt)
				var runningSinceString string
				if runningSince > 24*7*time.Hour {
					runningSinceString = fmt.Sprintf("%d days", int(runningSince.Hours())/24)
				} else if runningSince > time.Hour {
					if runningSince.Hours() > 1 {
						runningSinceString = fmt.Sprintf("%d hours", int(runningSince.Hours()))
					} else {
						// can replace this with static string but I'm curious whether this can ever be negative
						runningSinceString = fmt.Sprintf("%d hour", int(runningSince.Hours()))
					}
				} else if runningSince > time.Minute {
					if runningSince.Minutes() > 1 {
						runningSinceString = fmt.Sprintf("%d minutes", int(runningSince.Minutes()))
					} else {
						runningSinceString = fmt.Sprintf("%d minute", int(runningSince.Minutes()))
					}
				} else {
					if runningSince.Seconds() > 1 {
						runningSinceString = fmt.Sprintf("%d seconds", int(runningSince.Seconds()))
					} else {
						runningSinceString = fmt.Sprintf("%d second", int(runningSince.Seconds()))
					}
				}
				extraInfo += fmt.Sprintf("up %s", runningSinceString)
				if restarts > 0 {
					extraInfo += fmt.Sprintf(", %d restarts", restarts)
				}
				extraInfo = color.HiBlackString(extraInfo)
			}

			if extraInfo == "" {
				lines = append(lines, fmt.Sprintf("  [%s] %s", symbol, name))
			} else {
				extraInfo = color.HiBlackString("(") + extraInfo + color.HiBlackString(")")
				lines = append(lines, fmt.Sprintf("  [%s] %s %s", symbol, name, extraInfo))
			}
		}
	}

	out.RepaintLines(lines)
}
