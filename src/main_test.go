package main

import (
	"net"
	"os/exec"
	"testing"
	"time"
)

// runMainServer starts the main.go application as a subprocess using "go run".
// It applies the environment variable settings already set in t.Setenv.
// The function returns the started process and an error if any.
func runMainServer(t *testing.T) *exec.Cmd {
	cmd := exec.Command("go", "run", "main.go")
	// Forward stdout/stderr for debugging if needed.
	// Note that these could be directed to /dev/null if not desired.
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start main.go: %v", err)
	}
	// Give the server time to start.
	time.Sleep(2 * time.Second)
	return cmd
}

func assertPortListening(t *testing.T, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Expected server to be listening on %s: %v", addr, err)
	}
	conn.Close()
}

func TestMainDefaultPort(t *testing.T) {
	// Unset PORT so that default "8080" is used.
	t.Setenv("PORT", "")
	cmd := runMainServer(t)
	// Ensure the process is killed when the test ends.
	defer cmd.Process.Kill()

	assertPortListening(t, "localhost:8080")
}

func TestMainCustomPort(t *testing.T) {
	const customPort = "9090"
	t.Setenv("PORT", customPort)
	cmd := runMainServer(t)
	defer cmd.Process.Kill()

	assertPortListening(t, "localhost:"+customPort)
}
