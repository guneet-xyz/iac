package services

import (
	"fmt"
	"iac/utils/ansi"
	"iac/utils/errors"
	"iac/utils/exitcodes"
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
	RunE:  RunE,
}

func init() {
	Cmd.Flags().StringVar(&name, "name", "", "Specify a service to get info about")
}

func RunE(cmd *cobra.Command, args []string) error {
	services, err := service.GetServiceNames()
	if err != nil {
		return err
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
			return nil
		}
	} else {
		filteredServices = services
		if len(filteredServices) == 0 {
			slog.Info("No services found")
			os.Exit(exitcodes.NoServicesFound)
			return nil
		}
	}

	svcInfos := []service.ServiceInfo{}
	svcInfoChannel := make(chan service.ServiceInfo)

	svcInfoWg := sync.WaitGroup{}

	for _, svcName := range filteredServices {
		svcInfoWg.Go(func() { service.SendServiceInfoToChannel(svcName, svcInfoChannel) })
	}

	svcInfoWgFinished := false

	go func() {
		svcInfoWg.Wait()
		svcInfoWgFinished = true
		close(svcInfoChannel)
	}()

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	wg.Go(func() {
		for svcInfo := range svcInfoChannel {
			mutex.Lock()
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
			PrintServiceInfo(svcInfos)
			mutex.Unlock()
		}
	})

	wg.Go(func() {
		for {
			mutex.Lock()
			if svcInfoWgFinished {
				return
			}
			spinnerIteration++
			PrintServiceInfo(svcInfos)
			mutex.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	})

	wg.Wait()

	return nil
}

var linesPreviouslyPrinted int
var spinnerIteration int

var spinnerCharset = []string{
	"⠋",
	"⠙",
	"⠹",
	"⠸",
	"⠼",
	"⠴",
	"⠦",
	"⠧",
	"⠇",
	"⠏",
}

func spinnerChar(iteration int) string {
	return color.BlueString(spinnerCharset[iteration%len(spinnerCharset)])
}

var symbolX = color.RedString("✗")
var symbolQuestion = color.YellowString("?")
var symbolCheck = color.GreenString("✓")

func PrintServiceInfo(svcInfos []service.ServiceInfo) {
	lines := []string{}
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
				symbol = spinnerChar(spinnerIteration)
			case service.FetchingStatusNotFound:
				symbol = symbolX
			case service.FetchingStatusFound:
				switch container.ConfigInfoFound {
				case false:
					symbol = symbolQuestion
				case true:
					symbol = symbolCheck
				}
			}

			var extraInfo string
			switch container.InspectInfoStatus {
			case service.FetchingStatusFetching:
				extraInfo = spinnerChar(spinnerIteration)
			case service.FetchingStatusNotFound:
				symbol = symbolX
				extraInfo = ""
			case service.FetchingStatusFound:
				if !container.InspectInfo.Running {
					symbol = symbolX
					break
				}
				restarts := container.InspectInfo.RestartCount
				runningSince := time.Since(container.InspectInfo.StartedAt)
				var runningSinceString string
				if runningSince > 24*time.Hour {
					runningSinceString = fmt.Sprintf("%d days, %d hours", int(runningSince.Hours())/24, int(runningSince.Hours())%24)
				} else if runningSince > time.Hour {
					runningSinceString = fmt.Sprintf("%d hours, %d minutes", int(runningSince.Hours()), int(runningSince.Minutes())%60)
				} else if runningSince > time.Minute {
					runningSinceString = fmt.Sprintf("%d minutes, %d seconds", int(runningSince.Minutes()), int(runningSince.Seconds())%60)
				} else {
					runningSinceString = fmt.Sprintf("%d seconds", int(runningSince.Seconds()))
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

	for i := 0; i < linesPreviouslyPrinted; i++ {
		fmt.Print(ansi.LineUp + ansi.LineClear)
	}

	linesPreviouslyPrinted = len(lines)
	spinnerIteration++

	for _, line := range lines {
		fmt.Println(line)
	}
}
