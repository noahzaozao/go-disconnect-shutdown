package main

import (
    "log"
    "time"
    "os"
    "os/exec"
    "fmt"
    "github.com/prometheus-community/pro-bing"
)

func pingCheck(target string) bool {
    pinger, err := probing.NewPinger("192.168.1.1")
    if err != nil {
        panic(err)
        return false
    }
    pinger.Count = 1
    err = pinger.Run() // Blocks until finished.
    if err != nil {
        panic(err)
        return false
    }
    _ = pinger.Statistics() // get send/receive/duplicate/rtt stats
    // fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
    // fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
	// 	stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
	// fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
	// 	stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
    return true
}

func main() {
    logFile, err := os.OpenFile("shutdown.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal("Error opening log file:", err)
    }
    defer log.Printf("%s Closing log file...", currentTime.Format("2006-01-02 15:04:05"))
    defer logFile.Close()

    log.SetOutput(logFile)

    startTime := time.Now()
    for {
        currentTime := time.Now()
        if pingCheck("192.168.1.1") {
            fmt.Printf("%s %s\n", currentTime.Format("2006-01-02 15:04:05"), "Ping successful, continuing operation...")
            log.Printf("%s %s\n", currentTime.Format("2006-01-02 15:04:05"), "Ping successful, continuing operation...")
            time.Sleep(10 * time.Second) // Check every 10 seconds
        } else {
            if time.Since(startTime) > time.Minute {
                fmt.Printf("Ping failed for more than 1 minute, shutting down...")
                log.Println("Ping failed for more than 1 minute, shutting down...")
                err := exec.Command("shutdown", "/s", "/t", "60").Run() // Shutdown after 60 seconds
                if err != nil {
                    fmt.Printf("%s Error shutting down: %s", currentTime.Format("2006-01-02 15:04:05"), err)
                    log.Printf("%s Error shutting down: %s", currentTime.Format("2006-01-02 15:04:05"), err)
                }
                break
            } else {
                fmt.Printf("%s Ping failed, waiting for 1 minute...", currentTime.Format("2006-01-02 15:04:05"), )
                log.Println("%s Ping failed, waiting for 1 minute...", currentTime.Format("2006-01-02 15:04:05"), )
                time.Sleep(time.Minute) // Wait for 1 minute
            }
        }
    }
}
