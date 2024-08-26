# pi-monitor-api

<p align="center">
  <img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/alvmarrod/pi-monitor-api">
  <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/alvmarrod/pi-monitor-api?filename=src%2Fgo.mod">
  <img alt="GitHub License" src="https://img.shields.io/github/license/alvmarrod/pi-monitor-api">
</p>

## Overview

**pi-monitor-api** is a lightweight monitoring API designed to retrieve system information such as CPU load averages, RAM usage, storage details, and network statistics from a Raspberry Pi or similar linux based devices.

This program may be useful to you if you require the minimum functionality, avoiding any unnecesary features that may consume the scarce resources of your edge device.

The API is structured using hexagonal architecture principles, ensuring a clean separation of concerns and ease of testing and maintenance.

## Resource consumption

The following metrics have been measured using the containezed version and `docker stats`, not profiling tools over the raw binary.

| **Version** | **CPU** | **RAM** |
|:---:|:---:|:---:|
| v1.0.0 | 10m | ~10MB |

## Features

- **API Versioning**: All endpoints are grouped under a versioned path (`/v1`).
  - **CPU Monitoring**: Retrieve CPU load averages for the past 1, 5, and 15 minutes.
  - **RAM Monitoring**: Get detailed information about RAM usage, including total, available, free, and used memory.
  - **Storage Monitoring**: Access information on devices and partitions, including mount points, filesystem types, and storage utilization.
  - **Network Monitoring**: Fetch network interface statistics, including Rx/Tx packets, bytes, errors, drops, and link bitrate.

## Project Structure

The project is organized according to hexagonal architecture principles, with clear separation between core business logic, adapters, and infrastructure.

```bash
├── cmd
│   └── main.go                     # Application entry point
├── internal
│   ├── adapters                    # Implement the concrete versions of the ports for each domain
│   │   ├── handler                 #   HTTP handlers for endpoints: map them to service methods
│   │   └── repository              #   Repositories for accessing system information: interact with databases, files, or other storage systems to provide data
│   └── core
│       ├── domain                  # Domain models for CPU, RAM, Storage, and Network
│       ├── ports                   # Interfaces for interacting with repositories and services
│       └── services                # Business logic implementations for each domain
```

## Installation

### Prerequisites

- **Go 1.18+**: Ensure you have Go installed on your machine.
- **Git**: To clone the repository.

### Steps

1. Clone the repository and go into the src folder:

   ```bash
   git clone https://github.com/alvmarrod/pi-monitor-api.git
   cd pi-monitor-api
   ```

Now, short version:

```bash
make build
make run
```

Long version:

2. Install dependencies:

   ```bash
   cd src
   go mod tidy
   ```

3. Build the application:

   ```bash
   version=`cat ./../version.txt`
   go build -o pi_monitor_api_$version cmd/main.go
   ```

4. Run the application:

   ```bash
   ./pi-monitor-api
   ```

## API Endpoints

### CPU

- **GET `/v1/cpu`**
  - Retrieves CPU load averages for the last 1, 5, and 15 minutes.

### RAM

- **GET `/v1/ram`**
  - Returns total, available, free, and used memory.

### Storage

- **GET `/v1/storage`**
  - Provides information about storage devices and partitions, including mount points and usage.

### Network

- **GET `/v1/network`**
  - Fetches network interface statistics, including Rx/Tx packets, bytes, errors, drops, and link bitrate.

## Usage

You can call the API using tools like `curl` or Postman:

```bash
curl http://localhost:8080/v1/cpu
```

## Testing

1. Make sure you are at the root of the application: `./src`
2. Add any missing dependencies: `go mod tidy`
3. Start the service: `go run cmd/main.go`
4. Test with cURL or your favourite tool

## FAQ

- Q: does the service works out of the box as docker container?
- A: Although the service is deployable as a docker container, beware of unexpected results due to overlay network and filesystem for the container.

## Contributing

Contributions are welcome! Please submit issues or pull requests with your changes. Make sure to follow the existing code style and add tests where applicable.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
