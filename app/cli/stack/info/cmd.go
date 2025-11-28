package info

import (
	"fmt"
	"iac/utils/errors"
	"iac/utils/exitcodes"
	"iac/utils/out"
	"iac/utils/out/symbols"
	"iac/utils/stack"
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
	Use:   "info",
	Short: "info about the stack",
	RunE:  func(_ *cobra.Command, _ []string) error { return ShowInfoAboutStack(name) },
}

func init() {
	Cmd.Flags().StringVar(&name, "name", "", "Specify a stack to get info about")
}

func ShowInfoAboutStack(name string) error {
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	stackInfoWgFinished := false
	stackInfos := []stack.StackInfo{}
	stackInfoChannel := make(chan stack.StackInfo)

	wg.Go(func() {

		stacks, err := stack.GetStackNames()
		if err != nil {
			return
		}

		var filteredStacks []string

		if name != "" {
			for _, stack := range stacks {
				if stack == name {
					filteredStacks = append(filteredStacks, stack)
					break
				}
			}
			if len(filteredStacks) == 0 {
				slog.Info("Stack not found", "name", name)
				os.Exit(exitcodes.StackNotFound)
				return
			}
		} else {
			filteredStacks = stacks
			if len(filteredStacks) == 0 {
				slog.Info("No stacks found")
				os.Exit(exitcodes.NoStacksFound)
				return
			}
		}

		stackInfoWg := sync.WaitGroup{}

		for _, stackName := range filteredStacks {
			stackInfoWg.Go(func() { stack.SendStackInfoToChannel(stackName, stackInfoChannel) })
		}

		stackInfoWg.Wait()
		stackInfoWgFinished = true
		close(stackInfoChannel)

	})

	wg.Go(func() {
		for stackInfo := range stackInfoChannel {
			found := false
			for i, existingStackInfo := range stackInfos {
				if existingStackInfo.Name == stackInfo.Name {
					stackInfos[i] = stackInfo
					found = true
					break
				}
			}
			if !found {
				stackInfos = append(stackInfos, stackInfo)
			}
			mutex.Lock()
			printStackInfo(stackInfos)
			mutex.Unlock()
		}
	})

	wg.Go(func() {
		for {
			if stackInfoWgFinished {
				return
			}
			mutex.Lock()
			printStackInfo(stackInfos)
			mutex.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	})

	wg.Wait()

	return nil
}

func printStackInfo(stackInfos []stack.StackInfo) {
	lines := []string{}
	if len(stackInfos) == 0 {
		lines = append(lines, out.SpinnerChar())
	}

	for _, info := range stackInfos {
		lines = append(lines, info.Name)
		for _, container := range info.Containers {
			var name string
			if container.ConfigInfoFound {
				name = container.ConfigInfo.Name
			} else if container.RunningInfoStatus == stack.FetchingStatusFound {
				name = container.RunningInfo.ContainerName
			} else {
				err := errors.New("This should never happen. Container config info not found, but running info also not found")
				panic(err)
			}

			var symbol string
			switch container.RunningInfoStatus {
			case stack.FetchingStatusFetching:
				symbol = out.SpinnerChar()
			case stack.FetchingStatusNotFound:
				symbol = symbols.X
			case stack.FetchingStatusFound:
				switch container.ConfigInfoFound {
				case false:
					symbol = symbols.Question
				case true:
					symbol = symbols.Check
				}
			}

			var extraInfo string
			switch container.InspectInfoStatus {
			case stack.FetchingStatusFetching:
				extraInfo = out.SpinnerChar()
			case stack.FetchingStatusNotFound:
				symbol = symbols.X
				extraInfo = ""
			case stack.FetchingStatusFound:
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
