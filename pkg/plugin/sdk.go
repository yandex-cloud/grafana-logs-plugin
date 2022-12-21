package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/logging/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type sdk struct {
	*ycsdk.SDK
}

func newSDK(endpoint string, keyJson string) (*sdk, error) {
	creds, err := getSDKCreds(keyJson)
	if err != nil {
		return nil, err
	}

	ycSDK, err := buildSDK(endpoint, creds)
	if err != nil {
		return nil, err
	}

	return &sdk{SDK: ycSDK}, nil
}

const maxPageSize = 1000

func (s *sdk) readEntries(ctx context.Context, req loggingRequest, from time.Time, to time.Time) ([]*logging.LogEntry, error) {
	pageSize, pagesCount := optimalPageSize(req.Limit)

	lc := s.LogReading().LogReading()
	criteria := logging.Criteria{
		LogGroupId: req.GroupID,
		Since:      timestamppb.New(from),
		Until:      timestamppb.New(to),
		Levels:     req.sdkLevels(),
		Filter:     req.QueryText,
		PageSize:   pageSize,
	}
	if req.Stream != "" {
		criteria.StreamNames = []string{req.Stream}
	}
	if req.ResourceType != "" {
		criteria.ResourceTypes = []string{req.ResourceType}
	}
	if len(req.ResourceIDs) > 0 {
		criteria.ResourceIds = req.ResourceIDs
	}

	inReq := logging.ReadRequest{
		Selector: &logging.ReadRequest_Criteria{
			Criteria: &criteria,
		},
	}
	var entries []*logging.LogEntry
	for i := 0; i < pagesCount; i++ {
		resp, err := lc.Read(ctx, &inReq)
		if err != nil {
			return entries, fmt.Errorf("read logs: %w", err)
		}
		remains := req.Limit - len(entries)
		if remains > len(resp.Entries) {
			remains = len(resp.Entries)
		}
		entries = append(entries, resp.Entries[:remains]...)
		if len(resp.Entries) < int(pageSize) || resp.NextPageToken == "" {
			break
		}
		inReq.Selector = &logging.ReadRequest_PageToken{PageToken: resp.NextPageToken}
	}
	return entries, nil
}

func (s *sdk) listGroups(ctx context.Context, folderID string) (groups []*logging.LogGroup, err error) {
	req := logging.ListLogGroupsRequest{
		FolderId: folderID,
	}
	for {
		resp, err := s.Logging().LogGroup().List(ctx, &req)
		if err != nil {
			return nil, fmt.Errorf("list groups: %w", err)
		}
		groups = append(groups, resp.Groups...)
		if resp.NextPageToken == "" {
			break
		}
		req.PageToken = resp.NextPageToken
	}
	return groups, nil
}

func (s *sdk) listResources(ctx context.Context, groupID string) ([]*logging.LogGroupResource, error) {
	req := logging.ListResourcesRequest{
		LogGroupId: groupID,
	}
	resp, err := s.Logging().LogGroup().ListResources(ctx, &req)

	if err != nil {
		return nil, fmt.Errorf("list resources: %w", err)
	}
	return resp.Resources, nil
}

func getSDKCreds(keyJson string) (ycsdk.Credentials, error) {
	if keyJson == "" {
		return ycsdk.InstanceServiceAccount(), nil
	}
	var key *iamkey.Key
	if err := json.Unmarshal([]byte(keyJson), &key); err != nil {
		return nil, fmt.Errorf("api key unmarshal: %w", err)
	}
	return ycsdk.ServiceAccountKey(key)
}

func buildSDK(endpoint string, creds ycsdk.Credentials) (*ycsdk.SDK, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ycSDK, err := ycsdk.Build(ctx, ycsdk.Config{
		Credentials: creds,
		Endpoint:    endpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("sdk build: %w", err)
	}
	return ycSDK, nil
}

func optimalPageSize(limit int) (size int64, pages int) {
	switch {
	case limit < 1:
		return maxPageSize, 1
	case limit <= maxPageSize:
		return int64(limit), 1
	}
	pages = limit / maxPageSize
	if limit%maxPageSize != 0 {
		pages++
	}
	pageSize := limit / pages
	if limit%pages != 0 {
		pageSize++
	}
	return int64(pageSize), pages
}

func levelFromSDK(sdkLvl logging.LogLevel_Level) Level {
	switch sdkLvl {
	case logging.LogLevel_TRACE:
		return LevelTrace
	case logging.LogLevel_DEBUG:
		return LevelDebug
	case logging.LogLevel_INFO:
		return LevelInfo
	case logging.LogLevel_WARN:
		return LevelWarn
	case logging.LogLevel_ERROR:
		return LevelError
	case logging.LogLevel_FATAL:
		return LevelFatal
	default:
		return LevelUnknown
	}
}

func (r loggingRequest) sdkLevels() (result []logging.LogLevel_Level) {
	for _, l := range r.Levels {
		var sl logging.LogLevel_Level
		switch l {
		case LevelUnknown:
			sl = logging.LogLevel_LEVEL_UNSPECIFIED
		case LevelTrace:
			sl = logging.LogLevel_TRACE
		case LevelDebug:
			sl = logging.LogLevel_DEBUG
		case LevelInfo:
			sl = logging.LogLevel_INFO
		case LevelWarn:
			sl = logging.LogLevel_WARN
		case LevelError:
			sl = logging.LogLevel_ERROR
		case LevelFatal:
			sl = logging.LogLevel_FATAL
		default:
			panic(fmt.Sprintf("unknown log level %d", l))
		}
		result = append(result, sl)
	}
	return result
}
