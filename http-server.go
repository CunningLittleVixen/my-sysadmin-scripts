package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

var (
	scriptMu  sync.Mutex
	scriptCmd *exec.Cmd
	scriptPid int
)

const (
	scriptPath = "/usr/local/bin/script.sh"
	logFile    = "/var/log/script/monitor.log"
)

func startHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	scriptMu.Lock()
	defer scriptMu.Unlock()

	if scriptCmd != nil && scriptCmd.Process != nil {
		_, err := fmt.Fprintf(w, "Script already running (PID %d)\n", scriptPid)
		if err != nil {
			log.Println(err)
		}
		return
	}

	cmd := exec.Command(scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start script: %v", err), http.StatusInternalServerError)
		return
	}

	scriptCmd = cmd
	scriptPid = cmd.Process.Pid
	log.Printf("Script started with PID %d", scriptPid)

	go func() {
		err := cmd.Wait()
		scriptMu.Lock()
		if scriptCmd == cmd {
			scriptCmd = nil
			scriptPid = 0
		}
		scriptMu.Unlock()
		log.Printf("Script exited: %v", err)
	}()

	_, err := fmt.Fprintf(w, "Script started (PID %d)\n", scriptPid)
	if err != nil {
		log.Println(err)
	}
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	scriptMu.Lock()
	defer scriptMu.Unlock()

	if scriptCmd == nil || scriptCmd.Process == nil {
		_, err := fmt.Fprintf(w, "Script not running\n")
		if err != nil {
			log.Println(err)
		}
		return
	}

	pid := scriptPid
	if err := scriptCmd.Process.Signal(syscall.SIGTERM); err != nil {
		http.Error(w, fmt.Sprintf("Failed to stop script: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("Sent SIGTERM to PID %d", pid)

	// состояние сбросит горутина после Wait(); SIGKILL — если не успел за 5с
	go func(cmd *exec.Cmd) {
		time.Sleep(5 * time.Second)
		scriptMu.Lock()
		defer scriptMu.Unlock()
		if scriptCmd == cmd {
			log.Printf("Script did not exit, sending SIGKILL to PID %d", scriptPid)
			_ = cmd.Process.Kill()
		}
	}(scriptCmd)

	_, err := fmt.Fprintf(w, "Stop signal sent (PID %d)\n", pid)
	if err != nil {
		log.Println(err)
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, err := os.Open(logFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to open log: %v", err), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close log file: %v", err)
		}
	}()

	const tailLines = 19

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > tailLines {
			lines = lines[1:]
		}
	}
	if err := scanner.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to read log: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	for _, line := range lines {
		_, err = fmt.Fprintln(w, line)
		if err != nil {
			log.Println(err)
		}
	}
}

func okHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/stop", stopHandler)
	http.HandleFunc("/log", logHandler)
	http.HandleFunc("/", okHandler)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
