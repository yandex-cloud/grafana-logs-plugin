# Yandex Cloud Logging plugin for Grafana

Yandex Cloud Logging plugin for Grafana is an extension for Grafana that allows you to add [Cloud Logging](https://cloud.yandex.com/services/logging) as a data source.

With the plugin, services and applications can read log group records using the [filter expression language](https://cloud.yandex.com/docs/logging/concepts/filter).

## Requirements

1. [Create a service account](https://cloud.yandex.com/docs/iam/operations/sa/create#create-sa) and assign it the `logging.reader` role for your folder.
1. [Create an authorized key](https://cloud.yandex.com/docs/iam/operations/authorized-key/create) for the service account to authenticate with the Cloud Logging API.
1. [Create a log group](https://cloud.yandex.com/docs/logging/operations/create-group) and [add](https://cloud.yandex.com/docs/logging/operations/write-logs) records to it.

## Support

To ask questions about the plugin or report plugin issues, if any, open an issue in the [yandex-cloud/grafana-logs-plugin](https://github.com/yandex-cloud/grafana-logs-plugin) repository.

## Contribution

We appreciate the community's contribution to the plugin development. Check this [guide](https://github.com/yandex-cloud/grafana-logs-plugin/blob/master/CONTRIBUTING.md) before making a pull request (PR).

## License

[Apache 2.0 license](https://github.com/yandex-cloud/grafana-logs-plugin/blob/master/CONTRIBUTING.md)
