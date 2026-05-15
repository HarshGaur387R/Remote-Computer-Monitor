package constants

import (
	"os"
	"path/filepath"
)

const APP_NAME = "RCMA-Installer"
const AGENT_NAME = "rcma"
const AGENT_BINARY_NAME = "rcma.exe"
const AGENT_DISPLAY_NAME = "RCM AGENT"

const URL = "https://github.com/HarshGaur387R/Remote-Computer-Monitor/releases/download/v1.0.3-beta/v1.0.3_rcma_testing.exe"

var RCMA_BINARY_DIR_PATH = filepath.Join(os.Getenv("LOCALAPPDATA"), "RCMA")
var RCMA_BINARY_PATH = filepath.Join(RCMA_BINARY_DIR_PATH, AGENT_BINARY_NAME)

const CONFIG_DIR = `C:\ProgramData\RCMA`
const CONFIG_FILE_PATH = CONFIG_DIR + `\config.json`
const LOG_FILE_PATH = CONFIG_DIR + `\agent_logs.ndjson`