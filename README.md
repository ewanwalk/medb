# Media Encoding Dashboard
The goal of this application is to uniformly encode media libraries to be served through whatever applicable means (Plex / Other). Libraries are independantly configurable
which allows you to set quality standards however you would like!

### Dependencies

#### HandBrake
Must be installed or provided to the environment, the application will look for the
binary `HandBrakeCLI` as indicated by the default install:
```bash
sudo apt install handbrake-cli
```

Handbrake must have the `scan` option (latest versions)

#### Database

The following database servers are supported:

- MySQL 5.7+
- MariaDB

#### Configuration
All configuration is handled via. the `.env` file, copy the `.env.example` file and fill it out accordingly