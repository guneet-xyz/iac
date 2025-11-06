# IAC - Infrastructure as Code

## Requirements

- Backups of volume mounts using command
- Backups of volume mounts (automated)
- Health checks using command
- Health checks (automated)
- Monitoring (show resource usage)?? (not really needed but whatever)
- Deployment of all services

## Solutions

### Custom Go Thingy

A custom Go application that performs backups, health checks, and monitoring. This would require significant development effort but would allow for complete customization.

```
iac version
iac status - general info about the system (aggregate view)
iac info all - info about everything
iac info services - info about all services
iac info service <service-name> - info about a specific service
iac info backups - info about all backups
iac info backup <service-name> - info about a specific service backup
iac create backup all - create backups for all services
iac create backup <service-name> - create a backup for a specific service
iac config show - dump current config
iac up - deploy all services
iac up <service-name> - deploy a specific service
iac down - take down all services
iac down <service-name> - take down a specific service
iac healthcheck - perform health checks for all services
iac healthcheck <service-name> - perform a health check for a specific service
iac health - alias for healthcheck
```

#### global options

```
-v, -verbose   Enable verbose output
-dry, -dry-run   Simulate actions without making changes
```

#### config options

Since it's a CLI application, configuration should be done via a config file and use environment variables as overrides.
