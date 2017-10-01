# Media Encoding Database
medb is a tool made to automatically encode video libraries. This document will serve as a basic outline on the requirements for its use.

## Prerequisites

### MySQL / MariaDB

- download and install either MySQL or MariaDB (preferred)
  - https://downloads.mariadb.org/ (mariadb)
  - https://dev.mysql.com/downloads/mysql/ (mysql)
- create a database and name it whatever you like (this will be used in your configuration)

### SQLite (coming soon)

- provide a valid path to where you wish to store the database on disk

### Encoder

- download and install either ffmpeg or handbrakecli (preferred)
  - https://handbrake.fr/downloads2.php (handbrake)
  - https://www.ffmpeg.org/download.html (ffmpeg)
- make sure either encoder is accessible via. command line

handbrake
```
handbrakecli --version
```

ffmpeg
```
ffmpeg -version
```

## Configuration
All configuration is handled via. the **config.toml** file which should sit at the same level as the executable itself.  The format is of course **toml** (documentation of which may be [found here](https://github.com/toml-lang/toml))


### **log level**
option: `log_level`

example: `log_level = "info"`

determines the level of verbosity when logging, in order of verbosity (most to least)
- debug
- info (default)
- error

### **log path**
option: `log_path`

example: `log_path = "./medb.log"`

if provided, will log output to the location provided as well as to stdout

### **paths**
option: `[[path]]`

example: 
```
[[path]]
enabled = true
directory = "/path/to/library"
priority = 5
quality = 21
watch_interval = 500
minimum_size = "125MB"
```
Paths determine what your libraries are to the encoder. It is recommended that you provide one `[[path]]` per library of media to better assign profiles / priority to your media.

### **path sub-settings**
option: `enabled`

example: `enabled = true`

determines if the parent path is enabled or not (useful for disabling encoding on the provided library temporarily)

---
option: `directory`

example: `directory = "./media/shows"`

determines which directory the path will be linked to, this may be a relative (`./`) or absolute (`/`) path

---
option: `priority`

example: `priority = 10`

determines the priority of the path, this can be any number greater than zero, use this to weight which libraries should be encoded first

---
option: `quality`

example: `quality = 21`

determines the quality of the encode, use this to tune space to quality savings, this can range from 15-30 (lower being higher quality)

---
option: `watch_interval`

example: `watch_interval = 1000`

determines the time between event polling (looking for changes in files - add, remove, etc) this is a value in milliseconds.

---
option: `minimum_size`

example: `minimum_size = "125MB"`

determines the minimum size needed to qualify for encoding, if you adjust this it will be retroactivly applied to all existing media (e.g. lowering the limit will encode all those which previously did not meet the requirement)

### **media**
option: `[media]`

example:
```
[media]
extensions = [
  "mkv",
  "mp4",
  "avi",
  "flv",
  "m4v",
  "mov",
  "wmv",
  "webm"
]
```

determines the extensions we wish to keep track of, this is useful for ensuring you exclude all files of a certain type (e.g. don't keep track of images, etc)


### **encoder**
option: `[encoder]`

example:

```
[encoder]
override = "/bin/handbrakecli"
override_type = "handbrake"
cores = 2.0
concurrency = 1
staging = ""
```

overall settings for the encoder itself, used to tune resource allocation / usage

### **encoder sub-settings**
option: `override`

example: `override = "/bin/handbrakecli"`

used when you have a non-standard encoder setup and the encoder may not be reachable globally.

---
option: `override_type`

example: `override_type = "handbrake"`

used only when an `override` is provided. used to determine what kind of encoder you're providing (out of the currently supported two types)
- handbrake
- ffmpeg

---
option: `cores`

example: `cores = 0.33`

determines the number of cores to use, in the above example we are asking for 0.33 (or 33%) of the cores available. This must be provided as a float and will swap between whole (any number greater than or equal to 1) and percentage (any number less than 1)
- 0.33 = 33% of cores
- 3 = 3 cores

---
option: `concurrency`

example: `concurrency = 2`

determines the max number of concurrent encoding jobs that may run at a given point in time. This should be used sparingly on a system with few cores. This setting is representative of the following: `cores * concurrency = total cores`
- (cores = 2 * concurrency = 2) = 4 total cores

---
option: `staging`

example: `staging = "./staging"`

determines the path in which we want to stage our encoding jobs, it is important to set this to a path that has sufficient space to hold the largest of your files. Accepts both a relative and absolute path.
- ./staging (default)

### **database**
option: `[database]`

example:
```
[database]
type = "mysql"
path = ""
username = "user"
password = "pass"
hostname = "127.0.0.1"
name = "dbname"
```

determines the database configuration to be used for this setup. We use a database to keep track of your files and which of those have been encoded. We chose this method over attempting to alter a file by injecting data into it due to keeping thing sanitary and expandable.

### **database sub-settings**
option: `type`

example: `type = "mysql"`

determines which kind of database to use (of those supported). What you provide for this option determines what you will need to provide for the other options.
- mysql (default)
- sqlite (coming soon)

---
requirement: `type = "sqlite"`

option: `path`

example: `path = "./media.db"`

the location on disk to store the sqlite database

---
requirement: `type = "mysql"`

option: `username`

example: `username = "myuser"`

determines the username to use when connecting to the SQL server

---
requirement: `type = "mysql"`

option: `password`

example: `password = "mypass"`

determines the password to use when connecting to the SQL server

---
requirement: `type = "mysql"`

option: `hostname`

example: `hostname = "127.0.0.1"`

determines the host in which we attempt to form the SQL server connection

---
requirement: `type = "mysql"`

option: `name`

example: `name = "medb"`

determines the database to use with the SQL server


### **http**
option: `[http]`

example:
```
[http]
port = "4000"
```

determines the web server settings it should be noted that you maybe simply naviagte to `localhost:4000` (based on the above example, to view the dashboard.

### **http sub-settings**
option: `port`

example: `port = "4000"`

determines the port in which we wish to bind to for the purposes of viewing the web dashboard, you may bind this to port `80` taken it is not already in use - in which case you may navigate to use `localhost` to view the dashboard.

---
option: `username`

example: `username = "myuser"`

requirement: `password`

determines the username used for authentication on the dashboard

---
option: `password`

example: `password = "mypass"`

requirement: `username`

determines the password used for authentication on the dashboard

---
option: `ip_whitelist`

example: 
```
ip_whitelist = [
  "123.1.2.3",
  "123.1.2.4
]
```

restricts access to the dashboard via. an IP whitelist - does not work properly behind a reverse proxy (e.g. nginx) this may be changed in the future.

### A note on the source code
Currently I have not open sourced this tool, mostly due to things not being where I really want them to be at present. Once I get things cleaned up and moved into a better format I may add all the source!
