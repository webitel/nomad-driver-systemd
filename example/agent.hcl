log_level = "TRACE"

plugin "nomad-driver-systemd" {
  config {
    # Which units tasks may manage. Without this block every unit on the host
    # is eligible, so any job submitter can take over arbitrary host units.
    #
    # Patterns are unanchored, as in Go's regexp.MatchString: "nginx\\.service"
    # also matches "not-my-nginx.service". Anchor them with ^ and $.
    units {
      # Regexes matching the units tasks may manage. If the list is non-empty
      # the driver becomes an allowlist: a unit matching none of the patterns
      # is rejected.
      allowed = [
        "^nginx\\.service$",
        "^app-[a-z0-9-]+\\.service$",
      ]

      # Regexes matching units tasks may never manage. A match here always
      # wins, even over allowed.
      denied = [
        "^(nomad|consul|vault|ssh|sshd|systemd-.*)\\.service$",
        "\\.(mount|swap|device)$",
      ]
    }

    # Serve net/http/pprof on this address. Empty or unset disables it.
    # The profiles expose the memory of a process that talks to systemd over
    # DBus, so keep the address on loopback.
    # pprof_addr = "127.0.0.1:6061"
  }
}
