package plugin

import (
	"github.com/hashicorp/nomad/plugins/base"
	"github.com/hashicorp/nomad/plugins/drivers"
	"github.com/hashicorp/nomad/plugins/drivers/fsisolation"
	"github.com/hashicorp/nomad/plugins/shared/hclspec"
)

const (
	// pluginName is the name of the plugin as it will be known in Nomad
	pluginName = "systemd"

	// taskHandleVersion is the version of the task handle encoding
	taskHandleVersion = 1
)

var (
	// pluginVersion is the version reported to Nomad and in the fingerprint.
	// Release builds set it at link time with
	// -ldflags "-X github.com/webitel/nomad-driver-systemd/plugin.pluginVersion=<version>".
	pluginVersion = "dev"

	// pluginInfo is the response returned for the PluginInfo RPC
	pluginInfo = &base.PluginInfoResponse{
		Type:              base.PluginTypeDriver,
		PluginApiVersions: []string{drivers.ApiVersion010},
		PluginVersion:     pluginVersion,
		Name:              pluginName,
	}

	// capabilities is returned by the Capabilities RPC and indicates what
	// optional features this driver supports
	capabilities = &drivers.Capabilities{
		SendSignals: false,
		Exec:        false,
		FSIsolation: fsisolation.None,
		NetIsolationModes: []drivers.NetIsolationMode{
			drivers.NetIsolationModeHost,
		},
		MustInitiateNetwork: false,
		MountConfigs:        drivers.MountConfigSupportNone,
	}
)

var (
	// configSpec is the HCL specification for the driver configuration.
	configSpec = hclspec.NewObject(map[string]*hclspec.Spec{
		"units": hclspec.NewBlock("units", false, hclspec.NewObject(map[string]*hclspec.Spec{
			"allowed": hclspec.NewAttr("allowed", "list(string)", false),
			"denied":  hclspec.NewAttr("denied", "list(string)", false),
		})),
		"pprof_addr": hclspec.NewAttr("pprof_addr", "string", false),
	})

	// taskConfigSpec is the HCL specification for per-task configuration
	taskConfigSpec = hclspec.NewObject(map[string]*hclspec.Spec{
		"unit": hclspec.NewAttr("unit", "string", true),
	})
)
