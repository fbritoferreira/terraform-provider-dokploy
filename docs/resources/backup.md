---
page_title: "dokploy_backup Resource - dokploy"
subcategory: ""
description: |-
  Manages automated database backups in Dokploy.
---

# dokploy_backup (Resource)

Manages automated database backups in Dokploy. Backups are scheduled using cron expressions and stored in configured destinations (S3-compatible storage).

## Example Usage

### Basic Daily Backup

```terraform
resource "dokploy_backup" "daily" {
  database_id    = dokploy_database.postgres.id
  destination_id = dokploy_destination.s3.id
  schedule       = "0 2 * * *"  # Daily at 2 AM
  prefix         = "daily"
  enabled        = true
}
```

### Hourly Backup

```terraform
resource "dokploy_backup" "hourly" {
  database_id    = dokploy_database.postgres.id
  destination_id = dokploy_destination.s3.id
  schedule       = "0 * * * *"  # Every hour
  prefix         = "hourly"
  enabled        = true
}
```

### Weekly Backup

```terraform
resource "dokploy_backup" "weekly" {
  database_id    = dokploy_database.postgres.id
  destination_id = dokploy_destination.s3.id
  schedule       = "0 3 * * 0"  # Every Sunday at 3 AM
  prefix         = "weekly"
  enabled        = true
}
```

### Multiple Backup Schedules

```terraform
# Destination for all backups
resource "dokploy_destination" "backups" {
  name              = "backup-storage"
  storage_provider  = "s3"
  access_key        = var.s3_access_key
  secret_access_key = var.s3_secret_key
  bucket            = "db-backups"
  region            = "us-east-1"
  endpoint          = "https://s3.amazonaws.com"
}

# Production database
resource "dokploy_database" "prod_db" {
  project_id     = dokploy_project.main.id
  environment_id = dokploy_environment.prod.id
  name           = "production-db"
  type           = "postgres"
  password       = var.db_password
  version        = "16"
}

# Hourly backup for quick recovery
resource "dokploy_backup" "hourly" {
  database_id    = dokploy_database.prod_db.id
  destination_id = dokploy_destination.backups.id
  schedule       = "0 * * * *"
  prefix         = "hourly"
  enabled        = true
}

# Daily backup for standard recovery
resource "dokploy_backup" "daily" {
  database_id    = dokploy_database.prod_db.id
  destination_id = dokploy_destination.backups.id
  schedule       = "0 2 * * *"
  prefix         = "daily"
  enabled        = true
}

# Weekly backup for long-term retention
resource "dokploy_backup" "weekly" {
  database_id    = dokploy_database.prod_db.id
  destination_id = dokploy_destination.backups.id
  schedule       = "0 3 * * 0"
  prefix         = "weekly"
  enabled        = true
}
```

### Backup for MySQL Database

```terraform
resource "dokploy_database" "mysql" {
  project_id     = dokploy_project.main.id
  environment_id = dokploy_environment.prod.id
  name           = "mysql-db"
  type           = "mysql"
  password       = var.mysql_password
  version        = "8"
}

resource "dokploy_backup" "mysql_backup" {
  database_id    = dokploy_database.mysql.id
  destination_id = dokploy_destination.s3.id
  schedule       = "0 4 * * *"  # Daily at 4 AM
  prefix         = "mysql-daily"
  enabled        = true
}
```

### Backup for MongoDB

```terraform
resource "dokploy_database" "mongo" {
  project_id     = dokploy_project.main.id
  environment_id = dokploy_environment.prod.id
  name           = "mongo-db"
  type           = "mongo"
  password       = var.mongo_password
  version        = "7"
}

resource "dokploy_backup" "mongo_backup" {
  database_id    = dokploy_database.mongo.id
  destination_id = dokploy_destination.s3.id
  schedule       = "0 5 * * *"  # Daily at 5 AM
  prefix         = "mongo-daily"
  enabled        = true
}
```

### Backup for MariaDB

```terraform
resource "dokploy_database" "mariadb" {
  project_id     = dokploy_project.main.id
  environment_id = dokploy_environment.prod.id
  name           = "mariadb"
  type           = "mariadb"
  password       = var.mariadb_password
  version        = "10"
}

resource "dokploy_backup" "mariadb_backup" {
  database_id    = dokploy_database.mariadb.id
  destination_id = dokploy_destination.s3.id
  schedule       = "30 2 * * *"  # Daily at 2:30 AM
  prefix         = "mariadb-daily"
  enabled        = true
}
```

### Disabled Backup (for maintenance)

```terraform
resource "dokploy_backup" "maintenance" {
  database_id    = dokploy_database.postgres.id
  destination_id = dokploy_destination.s3.id
  schedule       = "0 2 * * *"
  prefix         = "daily"
  enabled        = false  # Temporarily disabled
}
```

### Complete Infrastructure Example

```terraform
# Project setup
resource "dokploy_project" "ecommerce" {
  name        = "ecommerce-platform"
  description = "E-commerce application"
}

resource "dokploy_environment" "production" {
  project_id = dokploy_project.ecommerce.id
  name       = "production"
}

# Backup destination
resource "dokploy_destination" "production_backups" {
  name              = "production-backups"
  storage_provider  = "s3"
  access_key        = var.aws_access_key
  secret_access_key = var.aws_secret_key
  bucket            = "ecommerce-backups"
  region            = "us-east-1"
  endpoint          = "https://s3.amazonaws.com"
}

# Main database
resource "dokploy_database" "main" {
  project_id     = dokploy_project.ecommerce.id
  environment_id = dokploy_environment.production.id
  name           = "ecommerce-db"
  type           = "postgres"
  password       = var.db_password
  version        = "16"
}

# Backup configuration with multiple schedules
resource "dokploy_backup" "main_hourly" {
  database_id    = dokploy_database.main.id
  destination_id = dokploy_destination.production_backups.id
  schedule       = "0 * * * *"
  prefix         = "ecommerce/hourly"
  enabled        = true
}

resource "dokploy_backup" "main_daily" {
  database_id    = dokploy_database.main.id
  destination_id = dokploy_destination.production_backups.id
  schedule       = "0 0 * * *"
  prefix         = "ecommerce/daily"
  enabled        = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `database_id` (String) ID of the database to backup.
- `destination_id` (String) ID of the backup destination.
- `enabled` (Boolean) Whether the backup is enabled.
- `prefix` (String) Prefix for backup files in the destination bucket.
- `schedule` (String) Cron expression for backup schedule.

### Read-Only

- `id` (String) Unique identifier for the backup configuration.

## Import

Import is supported using the following syntax:

```shell
# Backups can be imported using their ID
terraform import dokploy_backup.daily "backup-id-123"
```

## Cron Expression Reference

The `schedule` attribute uses standard cron expressions:

```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of the month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday = 0)
│ │ │ │ │
* * * * *
```

### Common Schedules

| Schedule | Cron Expression | Description |
|----------|-----------------|-------------|
| Every hour | `0 * * * *` | At minute 0 of every hour |
| Daily at midnight | `0 0 * * *` | At 00:00 every day |
| Daily at 2 AM | `0 2 * * *` | At 02:00 every day |
| Every 6 hours | `0 */6 * * *` | At minute 0 every 6 hours |
| Weekly on Sunday | `0 3 * * 0` | At 03:00 on Sunday |
| Monthly | `0 4 1 * *` | At 04:00 on the 1st of every month |
| Every 15 minutes | `*/15 * * * *` | Every 15 minutes |

## Notes

- Backups are stored in the destination bucket with the specified prefix.
- The actual backup file names include timestamps for easy identification.
- Backup retention must be managed separately (e.g., S3 lifecycle policies).
- Ensure the destination has sufficient storage space for your backup frequency.
- Monitor backup success/failure in the Dokploy dashboard.
