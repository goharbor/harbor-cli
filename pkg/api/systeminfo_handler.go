package api

import (
	"context"
	"fmt"
	"io"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/systeminfo"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

type SystemInfoAPI struct {
	Client *systeminfo.Client
}

func NewSystemInfoAPI(transport runtime.ClientTransport, formats strfmt.Registry, authInfo runtime.ClientAuthInfoWriter) *SystemInfoAPI {
	return &SystemInfoAPI{
		Client: systeminfo.New(transport, formats, authInfo),
	}
}

func (api *SystemInfoAPI) GetSystemInfo(ctx context.Context) (*models.GeneralInfo, error) {
	params := systeminfo.NewGetSystemInfoParams()
	resp, err := api.Client.GetSystemInfo(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get system info: %w", err)
	}
	return resp.Payload, nil
}

func (api *SystemInfoAPI) GetVolumes(ctx context.Context) (*systeminfo.GetVolumesOK, error) {
	params := systeminfo.NewGetVolumesParams()
	resp, err := api.Client.GetVolumes(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get volumes info: %w", err)
	}
	return resp, nil
}

func (api *SystemInfoAPI) GetCert(ctx context.Context, writer io.Writer) error {
	if writer == nil {
		return fmt.Errorf("writer cannot be nil")
	}
	params := systeminfo.NewGetCertParams()
	_, err := api.Client.GetCert(ctx, params, writer)
	if err != nil {
		return fmt.Errorf("failed to get certificate: %w", err)
	}
	return nil
}