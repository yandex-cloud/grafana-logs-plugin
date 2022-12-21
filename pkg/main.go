package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"

	"yc-logging/pkg/plugin"
)

func main() {
	if err := datasource.Manage("yandexcloud-logging-datasource", plugin.NewLoggingDatasource, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error("yandexcloud-logging-datasource error: %s", err.Error())
		os.Exit(1)
	}
}
