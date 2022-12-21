package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/logging/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"
)

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var keyFile string

var rootCmd = func() cobra.Command {
	cmd := cobra.Command{
		Use:  "logcli",
		Args: cobra.NoArgs,
	}
	cmd.PersistentFlags().StringVarP(&keyFile, "key-file", "k", "./key", "path to file with key")
	cmd.AddCommand(&tokenCmd)
	return cmd
}()
var tokenCmd = func() cobra.Command {
	cmd := cobra.Command{
		Use:  "groups",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			l, err := listGroups(cmd.Context())
			if err != nil {
				return fmt.Errorf("list groups: %w", err)
			}
			for _, g := range l {
				cmd.Println(g.Name, g.Id, g.Status.String())
			}
			return nil
		},
	}
	return cmd
}()

func listGroups(ctx context.Context) ([]*logging.LogGroup, error) {
	key, err := getKey()
	if err != nil {
		return nil, err
	}
	creds, err := ycsdk.ServiceAccountKey(key)
	if err != nil {
		return nil, fmt.Errorf("account credentials: %w", err)
	}
	fmt.Println("got creds", creds)
	sdk, err := ycsdk.Build(ctx, ycsdk.Config{
		Credentials:        creds,
		DialContextTimeout: time.Second,
		Endpoint:           "api.cloud-preprod.yandex.net:443",
	})
	if err != nil {
		return nil, fmt.Errorf("build sdk: %w", err)
	}

	var groups []*logging.LogGroup
	pt := ""
	for {
		resp, err := sdk.Logging().LogGroup().List(ctx, &logging.ListLogGroupsRequest{
			FolderId:  "aoeet5aqphu1um3hfjsj",
			PageToken: pt,
		})
		if err != nil {
			return nil, fmt.Errorf("create token: %w", err)
		}
		groups = append(groups, resp.GetGroups()...)
		if pt = resp.GetNextPageToken(); pt == "" {
			break
		}
	}

	return groups, nil
}

func getKey() (key *iamkey.Key, err error) {
	f, err := os.Open(keyFile)
	if err != nil {
		return nil, fmt.Errorf("key file: %w", err)
	}
	dec := json.NewDecoder(f)
	if err := dec.Decode(&key); err != nil {
		return nil, fmt.Errorf("key file decode: %w", err)
	}

	return key, nil
}
