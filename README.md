# Encoder Backend

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

Support for a sqlite3 database may be added in the future should demand warrant it.

#### Configuration
All configuration is handled via. the `.env` file, copy the `.env.example` file and fill it out accordingly