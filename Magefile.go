//go:build mage

package main

import (
	"fmt"
	"os"

	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
)

// Default configures the default target.
var Default = build.BuildAll

// Cleans up local folder
func CleanLocal() error {
	fmt.Println("Cleans the local folder")
	err := os.RemoveAll("yandexcloud-logging-datasource/")
	if err != nil {
		fmt.Println(err)
	}
	return err
}
