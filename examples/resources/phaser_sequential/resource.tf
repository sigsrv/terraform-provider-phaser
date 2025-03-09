terraform {
  required_providers {
    phaser = {
      source = "sigsrv/phaser"
    }
  }
}

resource "phaser_sequential" "example" {
  phases = ["prepare", "ready", "running"]
}

output "example_phase" {
  value = phaser_sequential.example.phase
}
