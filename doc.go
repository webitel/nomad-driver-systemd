// Command nomad-driver-systemd serves the systemd task driver to a Nomad
// client.
//
// It is not run directly: Nomad launches it and speaks to it over the plugin
// protocol. Everything it can be configured with comes from the plugin block of
// the Nomad client configuration; see the plugin package.
package main
