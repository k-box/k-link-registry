# K-Link Registry

The K-Link Registry creates a application management solution for registrants on a K-Link Network.

Each registrant (e.g.: technical staff of an organization) can add applications, such as K-Boxes or websites with a K-Search client through the graphical user interface the K-Link Registry provides. The added application then becomes an integral part of the K-Link Network. For each application the registrant creates, an associated authentication token is generated for its authorization on the network.

These applications are then entitled to perform certain actions on the K-Link Network. Which ones entirely depends on the permission level assigned to them:

- `data-search` - allows an application to access the search endpoint.
- `data-edit` - allows an application to edit a document located in the K-Link's search index.
- `data-add` - allows an application to add a document to the K-Link's search index.
- `data-view` - allows an application to view a document located in the K-Link's search index.
- `data-remove-own` - allows an application to delete a document which was previously added by the same application.
- `data-remove-all` - allows an application to delete any document within the search-index (e.g.: a future manage-my-search-index-super-admin-app)

Applications within a K-Link network interact through the K-Search API. A K-Search engine requests the permission level of applications by using the API of the connected K-Link Registry. This registry responds if the application is A) part of the network and B) has a certain set of required permissions. Based on this information the K-Search engine will either reply or refuse the request of an application to the K-Search API.

## Requirements

The K-Link Registry requires the following components to be installed in order to operate:

- `Docker 1.8+`

If configuring, then these additional components are needed as well:

- `composer`

Due to the nature of Docker, at least 2GB of RAM is required & an Intel Pentium 4 equivalent or newer/faster is required to operate this software. A Debian-based operating system is recommended & at least 1.5GB of free hard disk space is required for installation. Other environments have not been tested extensively but should be adequate if resources are doubled & the above requirements are met.

## Configuration

The K-Link Registry can be run using defaults out-of-the-box without configuration, however, it is strongly recommended to revise these defaults for optimal deployment & compatibility. Configuration parameters can be provided as environment variables directly through the `Dockerfile` or through the shell environment directly. These settings can then be deployed into an `.env` file directly by the Docker process. An `.env.dist` file is provided with defaults as an example. The parameters are explained with their significance below.

### DATABASE_HOST

* Default: `127.0.0.1`
* Type: IP address (ASCII string)
* This variable designates the database server IP address to be used with the K-Link Registry. If another server is specified, the local one bundled with the registry will be ignored & disabled.

### DATABASE_PORT

* Default: `3306`
* Type: integer
* This variable designates the database server port where the registry may communicate.

### DATABASE_NAME

* Default: `kregistry`
* Type: string (ASCII only)
* This variable designates the database to be used on the database server for storage of K-Link Registry data. Upon initialization, it will be created if the database server is pointing to the default instance.

### DATABASE_USER

* Default: `kregistry`
* Type: string
* This variable designates the database user to be used on the database server for access to K-Link Registry data. Upon initialization, it will be created if the database server is pointing to the default instance.

### DATABASE_ROOT_PASSWORD

* Default: `kregistry`
* Type: string
* This variable designates the database root password to be used on the database server for access to K-Link Registry data. Upon initialization, it will be created if the database server is pointing to the default instance.

### DATABASE_PASSWORD

* Default: `kregistry`
* Type: string
* This variable designates the database user password to be used on the database server for access to K-Link Registry data. Upon initialization, it will be created if the database server is pointing to the default instance.

### KREGISTRY_BASE_URL_PATH

* Default: `/registry/`
* Type: path (string)
* K-Suite software is designed to operate behind a proxy to provide a simple yet powerful user experience & interface. Thus, it is possible to operate the K-Link Registry behind a virtual sub folder or path of choice relative to the site root. This variable can be set to `/` otherwise to disable this feature completely or to configure a sub domain. Do not omit the trailing slash! No proxy configuration will be able to handle reversing & forwarding properly without it anyway.

### KREGISTRY_BASE_PATH

* Default: `/`
* Type: path (string)
* Due to the configurable nature of the platform, it is possible to specify a relative path to the project root to serve files. This is particularly useful if multiple web application installations are merged together into a single folder.

### KREGISTRY_BASE_PROTOCOL

* Default: `https`
* Type: protcol (ASCII string)
* Mainly `http` or `https` should be used here. Others will likely be problematic as requests are served through `http` & will need to take that into account if tunnelling.

### KREGISTRY_DOMAIN

* Default: none
* Type: domain (string)
* The domain name to used to serve client requests. Do note that the default value will not work in a public set up. Though this value does not need to be reciprocated by the proxy or environment, it does need to resolve to the correct IP address in order to work & validate as a trusted request. Email verification links will fail otherwise.

### KREGISTRY_ADMIN_USERNAME

* Default: `admin@example.com`
* Type: string
* This variable designates the global administrator username used to access the K-Link Registry system. It cannot be changed after installation.

### KREGISTRY_ADMIN_PASSWORD

* Default: none
* Type quoted & escaped string
* This variable is used to override & permanently assign the global administrator password which is used to administer the K-Link Registry system. The password must be hashed using the `bcrypt` algorithm & will default to `admin` if none is specified. Note that `$` symbols must be escaped as `\$` as per bash syntax & the resulting hash should be quoted using double quotes to ensure escaped characters are recognized.

### KSEARCH_CORE_IP_ADDRESS

* Default: `127.0.0.1,172.16.0.0/12,192.168.0.0/16,10.0.0.0/8,100.64.0.0/10,192.0.2.0/24,198.51.100.0/24,203.0.113.0/24,198.18.0.0/15`
* Type: comma separated CIDR values (string)
* This variable overrides default trusted networks that the registry server should trust. In particular, it controls which intermediate requests to the RPC-based application permission service are trusted & approved. This setting is critical for Docker-based set ups with multiple images that must communicate with one another. The default setting includes all private networks which can prove to be a security risk if the Docker network is not trustable, not in complete control of the network configuration or hosts a bridge that does not use restricted NAT.

### MAILER_TRANSPORT

* Default: `smtp`
* Type: protocol (ASCII string)
* This variable designates the protocol to be used by the mailer transport sub-system to send messages. This setting can be overridden in the K-Link Registry settings dashboard post-install.

### MAILER_HOST

* Default: `127.0.0.1`
* Type: IP address (ASCII string)
* This variable designates the mail server host to be used to send messages. This setting can be overridden in the K-Link Registry settings dashboard post-install.

### MAILER_USER

* Default: unset
* Type: string
* This variable designates the use of a username to access the mail server system. If `null`, no username will be used. This setting can be overridden in the K-Link Registry settings dashboard post-install.

### MAILER_PASSWORD

* Default: unset
* Type: string
* This variable designates the use of a password to access the mail server system. If `null`, no password will be used. This setting can be overridden in the K-Link Registry settings dashboard post-install.

### MAILER_PORT

* Default: `587`
* Type: integer
* This variable designates the port number with which to communicate with the mail server. This setting can be overridden in the K-Link Registry settings dashboard post-install.

### MAILER_SENDER_NAME

* Default: `"K-Registry Mailer Daemon"`
* Type: quoted & escaped string
* This variable designates the sender name used to send emails from the mailer sub system. This setting can be overridden in the K-Link Registry settings dashboard post-install.

### MAILER_SENDER_ADDRESS

* Default: `admin@example.com`
* Type: string
* This variable designates the sender from address used to send emails from the mailer sub system. This setting can be overridden in the K-Link Registry settings dashboard post-install.

### TOKEN_EXPIRATION_SECONDS

* Default: `1800`
* Type: integer
* Determines the interval of time in seconds a token issued is valid from the moment it is generated.

### APP_DEV

* Default: `prod`
* Type: enum (one of [`test`, `prod`, `dev`])
* This variable determines the regime under which the K-Link Registry web application will be released.

### APP_DEBUG

* Default: `0`
* Type: integer (anything non-zero will resolve to enabling debugging mode)
* This variable indicates whether trusted or sensitive information will be reported by the application to the end user. En bref, it enables debugging for developers or logging errors to retrieve report information on issues encountered during use.

### APP_SECRET

* Default: unset
* Type: string
* This variable is used to identify the application with a unique secret hash that can be used for private distribution network identification & use. Currently, it is not used as it is not published outside of the public license system.

## Installation

The K-Link Registry is deployed using a docker image that can be downloaded from the [K-Link](https://docker.klink.asia/main/klinkdocker_kregistry) repository using the following Docker terminal command:

`docker pull docker.klink.asia/main/klinkdocker_kregistry`

## Use

The K-Link Registry can be started once installed using the following terminal Docker command:

`docker run -p 80:80 docker.klink.asia/main/klinkdocker_kregistry`

## License

This project is licensed under the AGPL v3 license, see [LICENSE.txt](./LICENSE.txt).
