# Remote Computer Monitor

A comprehensive monorepo solution for remotely monitoring and controlling computers through a mobile application. This project combines mobile frontend, desktop GUI, and system monitoring agents to provide seamless remote computer management.

## 📋 Project Overview

Remote Computer Monitor is a full-stack application that enables users to:
- Monitor computer system metrics remotely (CPU, memory, network, etc.)
- Control and manage remote computers from a mobile device
- Install and manage monitoring agents on target computers
- Provide a user-friendly interface for system administration

## 🏗️ Project Structure

This is a monorepo managed with Turbo and npm workspaces:

```
Remote-Computer-Monitor/
├── apps/
│   ├── mobile/                 # React Native/Expo mobile app
│   ├── agent/                  # Go-based system monitoring agent
│   ├── GUI-installer/          # Go GUI installer for agent setup
│   └── agent-installer/        # Agent installer module
├── packages/
│   └── shared/                 # Shared TypeScript utilities and types
└── agent/                      # C++ native agent code
```

### Key Components

#### 1. **Mobile App** (`apps/mobile`)
- **Stack:** React Native (Expo), TypeScript
- **Features:**
  - Cross-platform support (iOS, Android, Web)
  - Real-time system monitoring dashboard
  - Remote computer control interface
  - Navigation and UI components using React Navigation
  - Camera and haptic feedback support
- **Dependencies:** Expo framework, React Router, Native modules

#### 2. **System Monitoring Agent** (`apps/agent`)
- **Stack:** Go 1.26.1
- **Features:**
  - Collects system metrics (CPU, memory, disk, network)
  - Runs as a system service
  - Communicates with mobile app
  - Cross-platform support (Windows, Linux, macOS)
- **Key Dependencies:** `gopsutil` (system metrics), `service` (Windows service management)

#### 3. **GUI Installer** (`apps/GUI-installer`)
- **Stack:** Go 1.26.1 with Fyne GUI framework
- **Features:**
  - User-friendly graphical installer
  - QR code generation for easy agent pairing
  - File system monitoring during installation
  - Cross-platform GUI support
- **Key Dependencies:** Fyne, fsnotify, QR code library

#### 4. **Shared Package** (`packages/shared`)
- **Stack:** TypeScript
- **Purpose:** Common types, utilities, and interfaces shared across the monorepo

## 🛠️ Tech Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Mobile | React Native / Expo | 54.0.33 |
| Frontend Frameworks | React Navigation, Expo Router | Latest |
| Agent | Go | 1.26.1 |
| GUI Installer | Fyne | 2.7.3 |
| Language Composition | TypeScript (61.6%), Go (33.6%), JavaScript (3.3%), PowerShell (1.5%) | - |
| Build Tool | Turbo | 2.8.16+ |

## 📦 Installation & Setup

### Prerequisites
- Node.js 18+ (for mobile and monorepo tools)
- Go 1.26.1+ (for agent and installer)
- CMake 3.10+ (for native agent build)
- npm or yarn

### Initial Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/HarshGaur387R/Remote-Computer-Monitor.git
   cd Remote-Computer-Monitor
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

### Building Components

#### Mobile App
```bash
# Start development server
npm run mobile

# Build for Android
npm run android

# Build for iOS
npm run ios

# Build web version
npm run web

# Production build
npm run build:mobile
```

#### System Agent
```bash
# Configure CMake
npm run cmake:agent

# Build (debug)
npm run debug:agent

# Build (release)
npm run build:agent

# Configure and build debug
npm run cmake:debug:agent
```

## 🚀 Usage

### Mobile Application
1. Install the mobile app on your device
2. Follow the setup wizard to connect to remote computers
3. Use the dashboard to monitor system metrics in real-time
4. Control remote systems through the mobile interface

### Installing the Agent
1. Download and run the GUI installer on the target computer
2. Follow the on-screen instructions
3. Use the QR code to pair with your mobile app
4. The agent will start monitoring the system

### Command Reference
- `npm run mobile` - Start mobile development server
- `npm run build:agent` - Build release version of monitoring agent
- `npm run cmake:agent` - Configure native agent build
- `npm run test` - Run tests (currently placeholder)

## 🏛️ Architecture

### Mobile-to-Agent Communication
- Mobile app communicates with agents running on remote computers
- Real-time metric updates via efficient data serialization
- Secure pairing mechanism using QR codes

### System Monitoring
- Lightweight Go agent runs as a system service
- Gathers CPU, memory, disk, and network metrics
- Cross-platform compatibility (Windows, macOS, Linux)

### Native Components
- Optional C++ native agent for enhanced performance
- CMake-based build system
- Integrates with Go agent wrapper

## 📄 License

ISC License

## 🤝 Contributing

Contributions are welcome! Please feel free to:
1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Open a pull request

## 📞 Support

For issues and questions, please use the GitHub Issues section of this repository.

---

**Last Updated:** June 2026  
**Repository:** [Remote Computer Monitor](https://github.com/HarshGaur387R/Remote-Computer-Monitor)
