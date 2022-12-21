package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var (
	_ backend.QueryDataHandler      = (*loggingDatasource)(nil)
	_ backend.CallResourceHandler   = (*loggingDatasource)(nil)
	_ backend.CheckHealthHandler    = (*loggingDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*loggingDatasource)(nil)
)

type loggingDatasource struct {
	logger log.Logger

	sdk      *sdk
	folderID string
}

func NewLoggingDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var pubConfig loggingConfig
	if err := json.Unmarshal(settings.JSONData, &pubConfig); err != nil {
		return nil, fmt.Errorf("unmarshal plugin config: %w", err)
	}
	sdk, err := newSDK(pubConfig.APIEndpoint, settings.DecryptedSecureJSONData[apiKeyJsonInSettings])
	if err != nil {
		return nil, fmt.Errorf("yc sdk: %w", err)
	}

	return &loggingDatasource{
		logger:   log.DefaultLogger,
		sdk:      sdk,
		folderID: pubConfig.FolderID,
	}, nil
}

func (o *loggingDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	resp := backend.NewQueryDataResponse()

	for _, query := range req.Queries {
		var lr loggingRequest
		if err := json.Unmarshal(query.JSON, &lr); err != nil {
			o.logger.Error("unmarshal query issue", "error", err.Error(), "ref_id", query.RefID, "query", query.JSON)
			return &backend.QueryDataResponse{}, fmt.Errorf("can not unmarshal query: %w", err)
		}
		if lr.GroupID == "" {
			continue
		}

		values := logEntriesValues{extractFields: lr.AddPayloadFields}
		respD := resp.Responses[query.RefID]

		entries, err := o.sdk.readEntries(ctx, lr, query.TimeRange.From, query.TimeRange.To)
		if err != nil {
			o.logger.Error("fetch logs error", "error", err.Error(), "ref_id", query.RefID)
			respD.Error = err
			resp.Responses[query.RefID] = respD
			continue
		}

		o.logger.Info("got entries", "count", len(entries), "ref_id", query.RefID)

		for _, ent := range entries {
			values.append(
				ent.Timestamp.AsTime(),
				levelFromSDK(ent.Level),
				ent.Uid,
				ent.StreamName,
				ent.Resource.Type,
				ent.Resource.Id,
				ent.Message,
				ent.JsonPayload.AsMap(),
			)
		}

		frame := data.NewFrame(query.RefID,
			data.NewField("timestamp", nil, values.timestamps),
			data.NewField("content", data.Labels{"group": lr.GroupID}, values.contents),
			data.NewField("level", nil, values.levels),
			data.NewField("id", nil, values.ids),
			data.NewField("stream", nil, values.streams),
			data.NewField("resource.type", nil, values.resourceTypes),
			data.NewField("resource.id", nil, values.resourceIDs),
			data.NewField("message", nil, values.messages),
			data.NewField("json_payload", nil, values.payloads),
		)
		frame.SetMeta(&data.FrameMeta{
			PreferredVisualization: data.VisTypeLogs,
		})

		respD.Frames = append(respD.Frames, frame)
		resp.Responses[query.RefID] = respD
	}
	return resp, nil
}

func (o *loggingDatasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	switch req.Path {
	case "suggestQuery":
		return o.suggestQuery(ctx, req, sender)
	default:
		o.logger.Warn("unknown resource call", "path", req.Path, "body", string(req.Body))
		return o.sendResourceJSON(sender, http.StatusNotFound, resourceError{Error: fmt.Sprintf("unknown path %q", req.Path)})
	}

}

func (o *loggingDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if err := o.sdk.CheckEndpointConnection(ctx, "log-reading"); err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "OK",
	}, nil
}

func (o *loggingDatasource) Dispose() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := o.sdk.Shutdown(ctx); err != nil {
		o.logger.Error("plugin dispose error", "error", err.Error())
	}

}

func (o *loggingDatasource) suggestQuery(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	var sreq suggestQueryRequest
	if err := json.Unmarshal(req.Body, &sreq); err != nil {
		o.logger.Error("invalid suggest query response", "error", err.Error(), "body", string(req.Body))
		return o.sendResourceJSON(sender, http.StatusBadRequest, resourceError{Error: err.Error()})
	}

	resp := suggestQueryResponse{
		Groups:        []string{},
		ResourceTypes: []string{},
		ResourceIDs:   []string{},
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if o.folderID == "" {
			return nil
		}
		groups, err := o.sdk.listGroups(ctx, o.folderID)
		if err != nil {
			return err
		}
		resp.Groups = make([]string, 0, len(groups))
		for _, g := range groups {
			resp.Groups = append(resp.Groups, g.Id)
		}
		return nil
	})

	eg.Go(func() error {
		if sreq.GroupID == "" {
			return nil
		}
		resources, err := o.sdk.listResources(ctx, sreq.GroupID)
		if err != nil {
			return err
		}

		resp.ResourceTypes = make([]string, 0, len(resources))
		for _, r := range resources {
			resp.ResourceTypes = append(resp.ResourceTypes, r.Type)
			if r.Type == sreq.ResourceType {
				resp.ResourceIDs = r.Ids
			}
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		o.logger.Error("suggestions failed", "error", err.Error())
		return err
	}

	return o.sendResourceJSON(sender, http.StatusOK, resp)
}

func (o *loggingDatasource) sendResourceJSON(sender backend.CallResourceResponseSender, status int, body any) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("resource json marshal: %w", err)
	}
	return sender.Send(&backend.CallResourceResponse{
		Status: status,
		Body:   bodyBytes,
	})
}
