package ydb

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/ydb-platform/ydb-go-sdk/v3/pkg/xerrors"
)

type Driver struct {
	core.Base
}

// DSN format: https://github.com/ydb-platform/ydb-go-sdk/blob/a804c31be0d3c44dfd7b21ed49d863619217b11d/connection.go#L339
func (d *Driver) Parse(driverName, dataSourceName string) (*core.Uri, error) {
	info := &core.Uri{DbType: "ydb"}

	uri, err := url.Parse(dataSourceName)
	if err != nil {
		return nil, xerrors.WithStackTrace(fmt.Errorf("failed on parse data source %v", dataSourceName))
	}

	const (
		secure   = "grpcs"
		insecure = "grpc"
	)

	if uri.Scheme != secure && uri.Scheme != insecure {
		return nil, xerrors.WithStackTrace(fmt.Errorf("unsupported scheme %v", uri.Scheme))
	}

	info.Host = uri.Host
	if spl := strings.Split(uri.Host, ":"); len(spl) > 1 {
		info.Host = spl[0]
		info.Port = spl[1]
	}

	info.DbName = uri.Path
	if info.DbName == "" {
		return nil, xerrors.WithStackTrace(errors.New("database path can not be empty"))
	}

	if uri.User != nil {
		info.Passwd, _ = uri.User.Password()
		info.User = uri.User.Username()
	}

	return info, nil
}
