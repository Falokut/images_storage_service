# Images storage service #

[![Go Report Card](https://goreportcard.com/badge/github.com/Falokut/images_storage_service)](https://goreportcard.com/report/github.com/Falokut/images_storage_service)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/Falokut/images_storage_service)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Falokut/images_storage_service)
[![Go](https://github.com/Falokut/images_storage_service/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/Falokut/images_storage_service/actions/workflows/go.yml) ![](https://changkun.de/urlstat?mode=github&repo=Falokut/images_storage_service)
[![Go Coverage](https://github.com/Falokut/images_storage_service/wiki/coverage.svg)](https://raw.githack.com/wiki/Falokut/images_storage_service/coverage.html)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)


# Content
+ [Images storage service](#images-storage-service)
+ [Docs](#swagger-docs)
+ [Params info](#configuration-params-info)
    + [Jaeger config](#jaeger-config)
    + [Prometheus config](#prometheus-config)
    + [Minio config](#minio-config)
+ [Author](#author)
+ [License](#license)

# Configuration params info
if supported values is empty, then any type values are supported

| yml name | yml section | env name | param type| description | supported values |
|-|-|-|-|-|-|
| log_level   |      | LOG_LEVEL  |   string   |      logging level        | panic, fatal, error, warning, warn, info, debug, trace|
| healthcheck_port   |      | HEALTHCHECK_PORT  |   string   |     port for healthcheck| any valid port that is not occupied by other services. The string should not contain delimiters, only the port number|
| storage_mode   |      | STORAGE_MODE  |   string   |service storage mode| MINIO or LOCAL, LOCAL is default, case insensitive|
| base_local_storage_path   |      | BASE_LOCAL_STORAGE_PATH  |   string   |path of images storage(relative or absolute path)||
| max_image_size   |      | MAX_IMAGE_SIZE  |   int   |max image size in bytes| only positive values|
| host   |  listen    | HOST  |   string   |  ip address or host to listen   |  |
| port   |  listen    | PORT  |   string   |  port to listen   | The string should not contain delimiters, only the port number|
| max_request_size   |  listen    | MAX_REQUEST_SIZE  |   int32   |  max request size in mb, by default 4 mb  |only > 0|
| max_response_size   |  listen    | MAX_RESPONSE_SIZE  |   int32   |  max response size in mb, by default 4 mb   |only > 0|
| server_mode   |  listen    | SERVER_MODE  |   string   | Server listen mode, Rest API, gRPC or both | GRPC, REST, BOTH|
|service_name|  prometheus    | PROMETHEUS_SERVICE_NAME | string |  service name, thats will show in prometheus  ||
|server_config|  prometheus    |   | nested yml configuration  [metrics server config](#prometheus-config) | |
|jaeger|||nested yml configuration  [jaeger config](#jaeger-config)|configuration for jaeger connection ||
|minio|||nested yml configuration  [minio config](#minio-config)|configuration for minio connection ||

## Jaeger config

|yml name| env name|param type| description | supported values |
|-|-|-|-|-|
|address|JAEGER_ADDRESS|string|ip address(or host) with port of jaeger service| all valid addresses formatted like host:port or ip-address:port |
|service_name|JAEGER_SERVICE_NAME|string|service name, thats will show in jaeger in traces||
|log_spans|JAEGER_LOG_SPANS|bool|whether to enable log scans in jaeger for this service or not||

## Prometheus config
|yml name| env name|param type| description | supported values |
|-|-|-|-|-|
|host|METRIC_HOST|string|ip address or host to listen for prometheus service||
|port|METRIC_PORT|string|port to listen for  of prometheus service| any valid port that is not occupied by other services. The string should not contain delimiters, only the port number|



# Minio config
|yml name| env name|param type| description | supported values |
|-|-|-|-|-|
|endpoint|MINIO_ENDPOINT|string|ip address or host to connect to the minio||
|access_key_id|MINIO_ACCESS_KEY_ID|string|secret id(like login) to access to minio||
|secret_access_key|MINIO_SECRET_ACCESS_KEY|string|secret(like password) to access to minio||
|secure|MINIO_SECURE|bool|enable ssh or not||

# Docs
[Swagger docs](swagger/docs/images_storage_service_v1.swagger.json)
 
 # Author

- [@Falokut](https://github.com/Falokut) - Primary author of the project

# License

This project is licensed under the terms of the [MIT License](https://opensource.org/licenses/MIT).

---