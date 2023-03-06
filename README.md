# Yandex Cloud Logging Data Source

## Introduction

The Yandex Cloud Logging Grafana Plugin can be used to extend Grafana by adding
[Cloud Logging](https://cloud.yandex.com/en/services/logging) as a data source in Grafana.

The plugin allows you to  query entries stored in logging groups in cloud store by
services or by user applications using [filter expression language](https://cloud.yandex.com/en/docs/logging/concepts/filter).

## Prerequisites

You should have Yandex Cloud account and configured log group with some records see [quick start guide](https://cloud.yandex.com/en/docs/logging/quickstart)

When you have log group you should create [service account](https://cloud.yandex.com/en/docs/iam/concepts/users/service-accounts)
and generate [api key](https://cloud.yandex.com/en/docs/iam/concepts/authorization/api-key) for authorization in Logging API.

Service account should have following roles in folder with log group:
- `logging.reader` for getting log entries
- `logging.viewer` for log groups listing


## Help

Issues and questions about this plugin can be posted as an issue in this GitHub repository.

## Contributing

This project welcomes contributions from the community. Before submitting a pull
request, please review our contribution guide.

## License

Licensed under the Apache License, Version 2.0

See LICENSE
