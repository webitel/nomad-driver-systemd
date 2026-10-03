package plugin

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/nomad/helper/pluginutils/hclutils"
)

func TestConfig_Decode(t *testing.T) {
	cases := []struct {
		name string
		hcl  string
		want Config
	}{
		{
			name: "empty config leaves every unit eligible",
			hcl:  `config {}`,
			want: Config{},
		},
		{
			name: "units block with both lists",
			hcl: `config {
				units {
					allowed = ["^app-.*\\.service$"]
					denied  = ["^nomad\\.service$", "^sshd?\\.service$"]
				}
				pprof_addr = "127.0.0.1:6061"
			}`,
			want: Config{
				Units: UnitsConfig{
					Allowed: []string{`^app-.*\.service$`},
					Denied:  []string{`^nomad\.service$`, `^sshd?\.service$`},
				},
				PprofAddr: "127.0.0.1:6061",
			},
		},
		{
			name: "units block with only denied",
			hcl: `config {
				units {
					denied = ["^nomad\\.service$"]
				}
			}`,
			want: Config{
				Units: UnitsConfig{Denied: []string{`^nomad\.service$`}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got Config
			hclutils.NewConfigParser(configSpec).ParseHCL(t, tc.hcl, &got)

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("decoded config = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestTaskConfig_Validate(t *testing.T) {
	cases := []struct {
		name    string
		unit    string
		wantErr bool
	}{
		// Accepted: the unit types systemd defines, plus templated instances.
		{name: "service", unit: "nginx.service"},
		{name: "service with dash", unit: "my-app.service"},
		{name: "service with underscore", unit: "my_app.service"},
		{name: "templated instance", unit: "app@1.service"},
		{name: "timer", unit: "backup.timer"},
		{name: "socket", unit: "app.socket"},
		{name: "mount", unit: "var-lib-data.mount"},
		{name: "automount", unit: "data.automount"},
		{name: "swap", unit: "swapfile.swap"},
		{name: "target", unit: "multi-user.target"},
		{name: "path", unit: "watched.path"},
		{name: "slice", unit: "system.slice"},
		{name: "scope", unit: "session-1.scope"},
		{name: "device", unit: "dev-sda.device"},

		// Rejected: anything that isn't a syntactically valid unit name.
		{name: "empty", unit: "", wantErr: true},
		{name: "no suffix", unit: "nginx", wantErr: true},
		{name: "unknown suffix", unit: "nginx.exe", wantErr: true},
		{name: "trailing space", unit: "nginx.service ", wantErr: true},
		{name: "leading space", unit: " nginx.service", wantErr: true},
		{name: "trailing newline", unit: "nginx.service\n", wantErr: true},
		{name: "embedded newline", unit: "nginx\n.service", wantErr: true},
		{name: "path traversal", unit: "../etc/passwd.service", wantErr: true},
		{name: "absolute path", unit: "/etc/systemd/system/x.service", wantErr: true},
		{name: "shell injection attempt", unit: "nginx.service; rm -rf /", wantErr: true},
		{name: "suffix only", unit: ".service", wantErr: true},
		{name: "over 255 chars", unit: strings.Repeat("a", 250) + ".service", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &TaskConfig{Unit: tc.unit}

			if err := cfg.Validate(); (err != nil) != tc.wantErr {
				t.Fatalf("Validate(%q) error = %v, wantErr %v", tc.unit, err, tc.wantErr)
			}
		})
	}
}
