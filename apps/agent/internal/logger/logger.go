package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/sys/windows/svc/eventlog"
)

const sourceName = "RCMA"

// LogEntry is the structure written as a single JSON line to agent_logs.ndjson.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// Logger writes progress entries to a file and errors to the Windows Event Log.
type Logger struct {
	mu       sync.Mutex
	file     *os.File
	eventLog *eventlog.Log
}

// New opens (or creates) the log file at logDir/agent_logs.ndjson and registers
// a Windows Event Log source. Call Close when the service stops.
func New(logDir string) (*Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}

	logPath := filepath.Join(logDir, "agent_logs.ndjson")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open agent_logs: %w", err)
	}

	_ = eventlog.InstallAsEventCreate(sourceName, eventlog.Error|eventlog.Warning|eventlog.Info)

	elog, err := eventlog.Open(sourceName)
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("open event log: %w", err)
	}

	return &Logger{file: f, eventLog: elog}, nil
}

func (l *Logger) write(level, msg string) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
	}
	line, _ := json.Marshal(entry)

	l.mu.Lock()
	defer l.mu.Unlock()
	l.file.Write(append(line, '\n')) // one JSON object per line (NDJSON)
}

func (l *Logger) Info(msg string) { l.write("INFO", msg) }

func (l *Logger) Infof(format string, args ...any) { l.Info(fmt.Sprintf(format, args...)) }

func (l *Logger) Error(msg string) {
	l.write("ERROR", msg)
	if l.eventLog != nil {
		_ = l.eventLog.Error(1, msg)
	}
}

func (l *Logger) Errorf(format string, args ...any) { l.Error(fmt.Sprintf(format, args...)) }

func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		l.file.Close()
	}
	if l.eventLog != nil {
		l.eventLog.Close()
	}
}