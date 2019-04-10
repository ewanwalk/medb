# Encoder Backend

### Dependencies

#### HandBrake
Must be installed or provided to the environment, the application will look for the
binary `HandBrakeCLI` as indicated by the default install:
```bash
sudo apt install handbrake-cli
```

#### Database

The following database servers are supported:

- MySQL 5.7+
- MariaDB

Support for a sqlite3 database may be added in the future should demand warrant it.