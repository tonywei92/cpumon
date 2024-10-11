package main

import (
	"flag"
	"fmt"
	"os/exec"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/shirou/gopsutil/v4/cpu"
)

func main() {
	var (
		flagInterval          = flag.Int("interval", 1, "Interval in seconds")
		flagSingle            = flag.Bool("single", true, "Single CPU")
		flagNotifyCPUMoreThan = flag.Int("notify-on-cpu-more-than", 0, "Notify on change, if notify-with is not set it will beep by default")
		flagNoBeep            = flag.Bool("no-beep", false, "Disable beep")
		flagNotifyWith        = flag.String("notify-with", "", "Run command on notification")
		flagVersion           = flag.Bool("version", false, "Show version")
	)

	flag.Parse()

	if *flagVersion {
		fmt.Println("cpumon v1.0.0 © 2024 Tony Soekirman")
		return
	}

	pw := progress.NewWriter()
	pw.Style().Visibility.Time = false
	pw.Style().Visibility.Value = false
	pw.Style().Colors.Percent = text.Colors{text.FgGreen}
	pw.Style().Colors.Tracker = text.Colors{text.FgBlue}
	pw.Style().Colors.Message = text.Colors{text.FgCyan}
	go pw.Render()

	if *flagNotifyCPUMoreThan > 100 {
		fmt.Println("CPU usage threshold cannot be more than 100%")
		return
	}

	if *flagNotifyWith != "" && *flagNotifyCPUMoreThan == 0 {
		fmt.Println("Notify command cannot be used without a CPU Notify threshold")
		return
	}

	fmt.Printf("Monitoring CPU usage every %d seconds. Press Ctrl+C to quit.\n", *flagInterval)
	warningMessage := ""
	commandRunMessage := ""

	if *flagSingle {

		tracker := progress.Tracker{
			Message: "CPU Usage",
			Total:   100,
			Units:   progress.UnitsDefault,
		}

		pw.AppendTracker(&tracker)

		for {
			percentagesSingle, err := cpu.Percent(0, false)
			warningMessage, commandRunMessage = notify(pw, warningMessage, commandRunMessage, true, percentagesSingle[0], *flagNotifyCPUMoreThan, !*flagNoBeep, *flagNotifyWith)
			if err != nil {
				fmt.Println("Error fetching CPU usage", err)
				return
			}

			// Update the progress bar
			tracker.SetValue(int64(percentagesSingle[0]))
			time.Sleep(time.Duration(*flagInterval) * time.Second)
		}

	} else {
		cpuCount, err := cpu.Counts(false)
		if err != nil {
			fmt.Println("Error fetching CPU count", err)
			return
		}

		trackers := make([]progress.Tracker, cpuCount)
		for i := 0; i < cpuCount; i++ {
			trackers[i] = progress.Tracker{
				Message: fmt.Sprintf("CPU #%d Usage", i),
				Total:   100,
				Units:   progress.UnitsDefault,
			}
			pw.AppendTracker(&trackers[i])
		}

		for {
			percentagesSingle, err := cpu.Percent(0, false)
			if err != nil {
				fmt.Println("Error fetching CPU usage", err)
				return
			}

			warningMessage, commandRunMessage = notify(pw, warningMessage, commandRunMessage, false, percentagesSingle[0], *flagNotifyCPUMoreThan, !*flagNoBeep, *flagNotifyWith)

			percentages, err := cpu.Percent(0, true)
			if err != nil {
				fmt.Println("Error fetching CPU usage", err)
				return
			}

			for i, percentage := range percentages {
				trackers[i].SetValue(int64(percentage))
			}

			time.Sleep(time.Duration(*flagInterval) * time.Second)
		}
	}

}

func notify(pw progress.Writer, lastWarningMessage string, lastCommandRunMessage string, skipSingleCpuUsageOutput bool, cpuUsage float64, cpuThreshold int, beep bool, command string) (string, string) {
	warningMessage := ""
	commandRunMessage := ""
	if cpuThreshold > 0 && (cpuUsage >= float64(cpuThreshold)) {
		if beep {
			beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		}
		if command != "" {
			cmd := exec.Command("sh", "-c", command)
			err := cmd.Run()
			if err != nil {
				commandRunMessage = fmt.Sprintf("Error running notify command: %s", err.Error())
			} else {
				commandRunMessage = fmt.Sprintf("Command ran successfully: %s", command)
			}
		}

		warningMessage = fmt.Sprintf("\nWARNING: CPU usage is more than %d%% on %s", cpuThreshold, time.Now().Format("2006-01-02 15:04:05"))
	}

	if warningMessage == "" {
		warningMessage = lastWarningMessage
	}

	if commandRunMessage == "" {
		commandRunMessage = lastCommandRunMessage
	}

	cpuUsageMessage := fmt.Sprintf("\nTotal CPU usage: %.2f%%", cpuUsage)
	if skipSingleCpuUsageOutput {
		pw.SetPinnedMessages(warningMessage, commandRunMessage)
	} else {
		pw.SetPinnedMessages(cpuUsageMessage, warningMessage, commandRunMessage)
	}

	return warningMessage, commandRunMessage
}
