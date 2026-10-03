job "systemd-example" {
  datacenters = ["dc1"]
  type        = "service"

  group "g" {
    task "unit-task" {
      driver = "systemd"

      config {
        unit = "my.service" # replace with your unit name, e.g., "nginx.service"
      }
    }
  }
}
