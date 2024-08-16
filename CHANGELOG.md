# Changelog

## [Unreleased]

### Added

### Fixed

### Changed

## [v1.0.0] - 2024-08-10

### Added

- **API Versioning**: Introduced versioning for the API under the `/v1` path prefix.
  - **CPU Endpoint**: Added the `/v1/cpu` endpoint to retrieve CPU load averages (1, 5, 15 minutes).
  - **RAM Endpoint**: Added the `/v1/ram` endpoint to retrieve RAM information including total, available, free, and used memory.
  - **Storage Endpoint**: Added the `/v1/storage` endpoint to retrieve storage information for devices and partitions, including mount points, filesystem types, and space usage.
  - **Network Endpoint**: Added the `/v1/network` endpoint to retrieve network interface statistics including Rx/Tx packets, bytes, errors, and drops, as well as link bitrate.
- **Service Layer**: Implemented service layer with CPU, RAM, Storage, and Network services to encapsulate business logic.
- **Repository Layer**: Implemented repository adapters for CPU, RAM, Storage, and Network to read system information from files and commands.
- **Handler Layer**: Implemented HTTP handlers for CPU, RAM, Storage, and Network endpoints using `gorilla/mux` for routing.
- **Project Structure**: Organized the project structure following hexagonal architecture principles with separate directories for `core`, `adapters`, and `cmd`.
