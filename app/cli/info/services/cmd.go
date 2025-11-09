package services

import (
	"fmt"
	"iac/utils/ansi"
	"iac/utils/exitcodes"
	"iac/utils/service"
	"log/slog"
	"os"
	"sync"
	"time"

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
			_spinnerIteration++
			PrintServiceInfo(svcInfos)
			mutex.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	})

	wg.Wait()

	return nil
}

var _linesPreviouslyPrinted int
var _spinnerIteration int

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
	return spinnerCharset[iteration%len(spinnerCharset)]
}

func PrintServiceInfo(svcInfos []service.ServiceInfo) {
	lines := []string{}
	for _, info := range svcInfos {
		lines = append(lines, fmt.Sprintf("+ %s", info.Name))
		for _, container := range info.Containers {
			switch container.RunningInfoStatus {
			case service.RunningInfoStatusFetching:
				lines = append(lines, fmt.Sprintf(" [%s] %s", spinnerChar(_spinnerIteration), container.ConfigInfo.Name))
			case service.RunningInfoStatusNotFound:
				lines = append(lines, fmt.Sprintf(" [✗] %s", container.ConfigInfo.Name))
			case service.RunningInfoStatusFound:
				lines = append(lines, fmt.Sprintf(" [✓] %s", container.ConfigInfo.Name))
			}
		}
	}

	for i := 0; i < _linesPreviouslyPrinted; i++ {
		fmt.Print(ansi.LineUp + ansi.LineClear)
	}

	_linesPreviouslyPrinted = len(lines)
	_spinnerIteration++

	for _, line := range lines {
		fmt.Println(line)
	}
}
