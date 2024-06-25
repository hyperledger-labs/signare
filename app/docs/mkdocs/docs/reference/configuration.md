# Configuration reference

This document describes the different signare command arguments and its static configuration file possible
attributes.

The target audience of this document is every user seeking precise configuration details.

## Static configuration

To configure the signare, you must use a YAML configuration file called `signare-config.yml`.
Here is a configuration example:

```yaml
logger:
  logLevel: 'debug'
database:
  postgresql:
    host: 'localhost'
    port: 5432
    scheme: 'postgres'
    database: 'db_signare'
    username: 'postgres'
    password: 'postgres'
    sslmode: 'disable'
metrics:
  prometheus:
    port: 9092
    path: /metrics
    maxRequestInFlight: 10
    timeoutInMillis: 30000
    namespace: 'signer'
hsmmodules:
  softhsm:
    lib: '/usr/local/lib/softhsm/libsofthsm2.so'
```

Let us dive into the different attributes:

| Name           | Format                                                  | Required | Description                       |
|----------------|---------------------------------------------------------|:--------:|-----------------------------------|
| **logger**     | [Logger configuration](#logger-configuration)           |    ✗     | Application logging configuration |
| **database**   | [Database configuration](#database-configuration)       |    ✔     | General database configuration    |
| **metrics**    | [Metrics configuration](#metrics-configuration)         |    ✗     | General metrics configuration     |
| **hsmmodules** | [HSM Modules configuration](#hsm-modules-configuration) |    ✔     | HSM Modules types configuration   |

### Logger configuration

| Name         | Type   | Required | Description                                            | Valid Values                 |
|--------------|--------|:--------:|--------------------------------------------------------|------------------------------|
| **logLevel** | string |    ✔     | Configuration of the  minimum level of logs to display | **INFO, WARN, ERROR, DEBUG** |

### Database configuration

| Name           | Type                                                    | Required | Description                                 |
|----------------|---------------------------------------------------------|:--------:|---------------------------------------------|
| **postgresql** | [PostgresSQL configuration](#postgressql-configuration) |    ✔     | Configuration of the postgres database type | 

!!! info
    The only supported databases are the ones that can be configured through this attribute.

#### PostgresSQL configuration

| Name          | Type                                                                   | Required | Description                                                        |
|---------------|------------------------------------------------------------------------|:--------:|--------------------------------------------------------------------|
| **host**      | string                                                                 |    ✔     | Host url of the database                                           | 
| **port**      | int                                                                    |    ✔     | Port number of the database                                        | 
| **scheme**    | string                                                                 |    ✔     | Scheme of the database system                                      | 
| **username**  | string                                                                 |    ✔     | Username to use for the database connection                        | 
| **password**  | string                                                                 |    ✔     | Password to use along with the username in the database connection | 
| **sslmode**   | string                                                                 |    ✔     | SSLMode to use in the database system                              | 
| **database**  | string                                                                 |    ✔     | Database name                                                      | 
| **sqlClient** | [PostgresSQL client configuration](#postgres-sql-client-configuration) |    ✗     | Configuration of the database client                               | 

#### Postgres SQL client configuration

| Name                      | Type | Required | Description                                         | Default Value (if any) |
|---------------------------|------|:--------:|-----------------------------------------------------|------------------------|
| **maxIdleConnections**    | int  |    ✗     | Max idle connections for the database/sql handle    | 2                      |
| **maxOpenConnections**    | int  |    ✗     | Max open connections for the database/sql handle    | 100                    |
| **maxConnectionLifetime** | int  |    ✗     | Max connection lifetime for the database/sql handle | 0                      |

### Metrics configuration

| Name           | Type                                                          | Required | Description                            |
|----------------|---------------------------------------------------------------|:--------:|----------------------------------------|
| **prometheus** | [Prometheus configuration](#prometheus-metrics-configuration) |    ✔     | Configuration of the Prometheus system | 

!!! info
    The only supported monitoring systems are the ones that can be configured through this attribute.

#### Prometheus metrics configuration

| Name                    | Type   | Required | Description                                          | Default Value (if any) |
|-------------------------|--------|:--------:|------------------------------------------------------|------------------------|
| **port**                | int    |    ✗     | Port number where Prometheus metrics will be exposed | 9780                   |
| **path**                | string |    ✗     | URL path where prometheus will listen                | /metrics               |
| **maxRequestsInFlight** | int    |    ✗     | Number of concurrent HTTP requests                   | 10                     |
| **timeoutInMillis**     | int    |    ✗     | Number of millis until timeout                       | 30000                  |
| **namespace**           | string |    ✗     | Namespace to prefix metric names                     | signer                 |

### HSM Modules configuration

signare requires configuration of at least one HSM type to function.

| Name        | Type                                            | Required | Description                            |
|-------------|-------------------------------------------------|:--------:|----------------------------------------|
| **softhsm** | [SoftHSM configuration](#softhsm-configuration) |    ✗     | Configuration of the Prometheus system | 

!!! info
    The only supported HSM systems are the ones that can be configured through this attribute.

#### SoftHSM Configuration

| Name        | Type   | Required | Description                              | Default Value (if any) |
|-------------|--------|:--------:|------------------------------------------|------------------------|
| **library** | string |    ✔     | Library path to the softHSM installation |                        |

## Command flags

When executing the signare binary, a multitude of flags are at your disposal in order to customize some of its
configuration.

Let us delve deeper into the specifics to further describe the flag options:

| Name                     | Type   | Required | Description                                                  | Default Value (if any) |
|--------------------------|--------|:--------:|--------------------------------------------------------------|------------------------|
| **signer-administrator** | string |    ✔     | Id of the signare's initial admin                      |                        |
| **config**               | string |    ✔     | Path to where the config yml file is stored                  |                        |
| **listen-address**       | string |    ✗     | Address where the signare will listen                  | 0.0.0.0                |
| **http-port**            | int    |    ✗     | Number of the port where REST API methods will be hosted     | 32325                  |
| **rpc-port**             | int    |    ✗     | Number of the port where JSON RPC API methods will be hosted | 4545                   |




