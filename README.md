<h1 align="center">
  <br>
  <a href="https://github.com/ad-on-is/resticity"><img src="./images/resticity-logo.svg" alt="Resticity" width="100"></a>
  <br>
  Resticity
  <br>
</h1>

<h4 align="center">A beautiful cross-platform UI for <a href="https://restic.readthedocs.io/en/stable/" target="_blank">restic</a>, built with <a href="https://wails.io" target="_blank">Wails</a>.</h4>

![screenshot](./images/resticity_screenshot.png)

## Key Features

- Easy to use
- Cross platform
  - Linux
  - Windows
  - MacOS
- Docker image to run on self-hosted servers
- Scheduled backups
- Supports local and remote repositories
  - Local folder or mounted network drive
  - AWS/Backblaze
  - Google
- System tray support

## How To Use

```bash
# Run in GUI mode
$ resticity

# Run in background mode (useful for autostart)
$ resticity --no-gui

# Run in server-only mode (this is the default mode for Docker images)
$ resticity --server

# Run with custom configuration path
$ resticity --config /path/to/config.json

# Run with Docker
# Assign the volumes/paths that you want resticity to grant access to
$ docker run -d --name resticity -p 11278:11278 -v /path/to/config.json:/config.json -v /mnt:/mnt -v /home:/home ghcr.io/ad-on-is/resticity/resticity
```

## Configuration

Resticity looks for a configuration file in the following order:

1. Custom file location with the `--config path/to/config.json` flag
2. `RESTICITY_SETTINGS_FILE` environment variable
3. `$XDG_CONFIG_HOME/resticity/config.json`

## Installation

## Build yourself

```bash
# Clone this repo
$ git clone https://github.com/ad-on-is/resticity

# cd into resticity
$ cd resticity

# Build using wails
$ wails build

```

## License

MIT

---

> [adisdurakovic.com](https://adisdurakovic.com) &nbsp;&middot;&nbsp;
> GitHub [@ad-on-is](https://github.com/ad-on-is)
