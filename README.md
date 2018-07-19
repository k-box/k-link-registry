# K-Link-Registry

The K-Link-Registry is a Web Application that allows its Users, the
"Registrants", to modify Applications served by K-Link.

## Installation
Precompiled binaries as well as docker images are provided for easy
installation.

## Migrating from K-Link-registry 1
K-Link registry builds ontop of the API schema for the last iteration of
the registry. It can be used with a preexisting Database, but this database
needs to be migrated first. Refer to the documentation of the `migrate`
subcommand.

## Usage
K-Link registry exposes two APIs: Version 1.0 is currently used as an
endpoint for application validation, Version 2.0 is used by the frontend
to get access to the managed resources, but can be used by other apps in
the same way, once they are authenticated.

## Configuration
The software can be configured via a config file, environment variables or
flags. To learn about the various options, run the binary using the `help`
command, i.e. `klinkregistry help`.

#### Base config
Used by every command

| flag      | ENV                   | description                                                     |
|-----------|-----------------------|-----------------------------------------------------------------|
| config    | -                     | config file to use                                              |
| assets    | `REGISTRY_ASSETS_DIR` | Assets dir (default empty, embedded assets will be used)        |
| db-host   | `REGISTRY_DB_HOST`    | Database host (default: "database")                             |
| db-port   | `REGISTRY_DB_PORT`    | Database Port (default: "3306")                                 |
| db-user   | `REGISTRY_DB_USER`    | Database User (default: "kregistry")                            |
| db-pass   | `REGISTRY_DB_PASS`    | Database Password (default: "kregistry")                        |
| db-name   | `REGISTRY_DB_NAME`    | Database Name (default: "kregistry")                            |
| smtp-host | `REGISTRY_SMTP_HOST`  | Mail Host (default: empty, logger will be used to output mails) |
| smtp-port | `REGISTRY_SMTP_PORT`  | Outgoing mail Port (default: 25)                                |
| smtp-user | `REGISTRY_SMTP_USER`  | Mail user (default: kregistry)                                  |
| smtp-pass | `REGISTRY_SMTP_HOST`  | Mail Password (default: registry)                               |
| smtp-from | `REGISTRY_SMTP_HOST`  | From Address (default: registry@example.com)                    |

### `server` config
Used by the `server` subcommand

| flag               | ENV                           | description                                                  |
|--------------------|-------------------------------|--------------------------------------------------------------|
| http               | `REGISTRY_HTTP_LISTEN`        | Address for the HTTP server to listen on (default: ":80")    |
| http-read-timeout  | `REGISTRY_HTTP_READ_TIMEOUT`  | Timeout duration for HTTP read (default: "10s")              |
| http-write-timeout | `REGISTRY_HTTP_WRITE_TIMEOUT` | Timeout duration for HTTP write (default: "10s")             |
| http-max-header    | `REGISTRY_HTTP_MAX_HEADER`    | Maximal HTTP Header size, in bytes. (default: 1MB)           |
| domain             | `REGISTRY_HTTP_DOMAIN`        | Domain used for generation of links (default: "example.com") |
| base-path          | `REGISTRY_HTTP_BASE_PATH`     | Base path the application is served on (default: "/")        |
| http-secret        | `REGISTRY_HTTP_SECRET`        | Secret string for session generation (default: generated)    |
| admin-username     | `REGISTRY_ADMIN_USERNAME`     | Username (email) for admin account                           |
| admin-password     | `REGISTRY_ADMIN_PASSWORD`     | Password for admin account                                   |

###  `migrate` config
This command uses the base configuration

Supported arguments:
* `up` migrates to the latest database revision, this is probably what you have in mind
* `down` cleans the database
* `1` migrates to the specific revision number, '1' in this case

## Known Bugs
Users cannot change their name once registred, except if they are owner