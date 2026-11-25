# Studylog

A lightweight command-line application for tracking study sessions and aggregating study time across university modules.

Studylog is written in Go and uses the standard library exclusively. It is designed to provide a simple, terminal-based alternative to manual time tracking while keeping recorded data accessible through standard file formats.

## Project Status

Studylog is currently in early development. The core functionality is implemented, but further work on reliability, testing and integration is planned before the first stable release.

## Features

- Start and stop study sessions.
- Display the currently active session and elapsed time.
- Attach optional notes to completed sessions.
- Aggregate daily and weekly study time by module.
- Store completed sessions in CSV format.
- Persist active sessions between program executions.
- Configure the data directory through an environment variable.
- Migrate older CSV schemas automatically.

## Installation

### Requirements

- Go 1.26.10 or a compatible newer version
- Git

### Build from source

```sh
git clone https://github.com/milisaur/studylog.git
cd studylog
go build .
```

### Install

```sh
go install .
```

Ensure that the Go binary directory is included in your `PATH`.

## Usage

Start a study session:

```sh
studylog start 61111
```

Display the active session:

```sh
studylog status
```

Complete a session:

```sh
studylog stop
```

Complete a session with a note:

```sh
studylog stop "Linear independence exercises"
```

Display daily statistics:

```sh
studylog today
```

Display weekly statistics:

```sh
studylog week
```

Weeks are defined as Monday through Sunday.

## Configuration

By default, Studylog stores its data in:

```text
~/.local/share/studylog/
```

A custom storage directory can be specified using the `STUDYLOG_DATA_DIR` environment variable:

```sh
export STUDYLOG_DATA_DIR="$HOME/studylog-data"
```

## Data Storage

Studylog uses two files:

| File | Purpose |
|---|---|
| `active.json` | Currently active study session |
| `sessions.csv` | Completed study sessions |

Completed sessions are stored using the following CSV schema:

```csv
module,start_time,end_time,duration_seconds,note
```

Timestamps are stored in RFC 3339 format with fractional-second precision. Durations are stored as whole seconds.

The CSV format allows recorded data to be processed independently of Studylog, including through spreadsheet applications and external scripts.

## Project Structure

```text
studylog/
├── main.go
├── config.go
├── session.go
├── storage.go
├── statistics.go
├── go.mod
├── README.md
└── LICENSE
```

- `main.go` — application entry point and command routing.
- `config.go` — storage configuration.
- `session.go` — session data structures and lifecycle management.
- `storage.go` — JSON and CSV persistence.
- `statistics.go` — filtering and aggregation of completed sessions.

## Known Limitations

- Only one active session is supported per data directory.
- Concurrent command execution is not fully protected.
- Interrupted writes or failed session cleanup may result in incomplete or duplicate records.
- Sessions spanning midnight are attributed entirely to their starting date.
- Statistics exclude currently active sessions.
- Calendar boundaries are determined using the system's local time zone.
- Editing, deleting and pausing sessions are not currently supported.
- Error handling and exit codes are not yet consistent across all commands.
- Automated test coverage has not yet been established.
- Multi-device synchronization is not currently supported or guaranteed to be safe.
- The application is primarily developed for Linux. Other platforms have not been systematically tested.

## Roadmap

Planned improvements include:

- Consistent CLI error handling and exit codes.
- Improved data integrity and session recovery.
- Protection against duplicate records.
- Automated unit and integration tests.
- XDG Base Directory support.
- Nix packaging.
- Spreadsheet integration.
- Safe multi-device workflows.

The scope and implementation of future features may change as the project develops.

## Dependencies

Studylog currently uses only the Go standard library.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
