---
page_title: "dokploy_application Resource - dokploy"
subcategory: ""
description: |-
  Manages a Dokploy application. Applications can be deployed from various sources including Docker images, Git repositories, and GitHub.
---

# dokploy_application (Resource)

Manages a Dokploy application. Applications can be deployed from various sources including Docker images, Git repositories (public or private), and GitHub repositories. This resource supports various build types and deployment configurations.

## Example Usage

### Docker Image Deployment

```terraform
resource "dokploy_project" "myproject" {
  name = "My Application"
}

resource "dokploy_application" "nginx" {
  name         = "nginx-app"
  project_id   = dokploy_project.myproject.id
  source_type  = "docker"
  docker_image = "nginx:alpine"
  
  deploy_on_create = true
}
```

### Docker Image with Registry Authentication

```terraform
resource "dokploy_application" "private_app" {
  name         = "private-app"
  project_id   = dokploy_project.myproject.id
  source_type  = "docker"
  docker_image = "registry.example.com/myorg/myapp:latest"
  registry_url = "registry.example.com"
  username     = "deploy"
  password     = var.registry_password
  
  deploy_on_create = true
}
```

### Git Repository with Dockerfile

```terraform
resource "dokploy_application" "api" {
  name              = "api"
  project_id        = dokploy_project.myproject.id
  source_type       = "git"
  custom_git_url    = "https://github.com/myorg/api.git"
  custom_git_branch = "main"
  build_type        = "dockerfile"
  dockerfile_path   = "Dockerfile"
  
  auto_deploy      = true
  deploy_on_create = true
}
```

### Private Git Repository with SSH Key

```terraform
resource "dokploy_ssh_key" "deploy_key" {
  name        = "api-deploy-key"
  description = "Deploy key for API repository"
  public_key  = file("~/.ssh/deploy_rsa.pub")
  private_key = file("~/.ssh/deploy_rsa")
}

resource "dokploy_application" "private_api" {
  name                  = "private-api"
  project_id            = dokploy_project.myproject.id
  source_type           = "git"
  custom_git_url        = "git@github.com:myorg/private-api.git"
  custom_git_branch     = "main"
  custom_git_ssh_key_id = dokploy_ssh_key.deploy_key.id
  build_type            = "dockerfile"
  dockerfile_path       = "Dockerfile"
  
  deploy_on_create = true
}
```

### Nixpacks Build

```terraform
resource "dokploy_application" "nodejs_app" {
  name              = "nodejs-app"
  project_id        = dokploy_project.myproject.id
  source_type       = "git"
  custom_git_url    = "https://github.com/myorg/nodejs-app.git"
  custom_git_branch = "main"
  build_type        = "nixpacks"
  
  auto_deploy      = true
  deploy_on_create = true
}
```

### Static Site

```terraform
resource "dokploy_application" "website" {
  name              = "website"
  project_id        = dokploy_project.myproject.id
  source_type       = "git"
  custom_git_url    = "https://github.com/myorg/website.git"
  custom_git_branch = "main"
  build_type        = "static"
  publish_directory = "dist"
  
  deploy_on_create = true
}
```

### Multi-stage Dockerfile Build

```terraform
resource "dokploy_application" "app" {
  name               = "app"
  project_id         = dokploy_project.myproject.id
  source_type        = "git"
  custom_git_url     = "https://github.com/myorg/app.git"
  custom_git_branch  = "main"
  build_type         = "dockerfile"
  dockerfile_path    = "Dockerfile"
  docker_build_stage = "production"
  docker_context_path = "."
  
  build_args = <<-EOT
    NODE_ENV=production
    API_URL=https://api.example.com
  EOT
  
  deploy_on_create = true
}
```

### Application with Resource Limits

```terraform
resource "dokploy_application" "limited_app" {
  name         = "limited-app"
  project_id   = dokploy_project.myproject.id
  source_type  = "docker"
  docker_image = "myorg/myapp:latest"
  
  # CPU limits (in millicores, 1000 = 1 CPU)
  cpu_limit       = 2000  # 2 CPUs max
  cpu_reservation = 500   # 0.5 CPU guaranteed
  
  # Memory limits (in bytes)
  memory_limit       = 536870912  # 512MB max
  memory_reservation = 268435456  # 256MB guaranteed
  
  replicas = 2
  
  deploy_on_create = true
}
```

### Application with Custom Command

```terraform
resource "dokploy_application" "worker" {
  name         = "worker"
  project_id   = dokploy_project.myproject.id
  source_type  = "docker"
  docker_image = "myorg/myapp:latest"
  command      = "npm run worker"
  
  deploy_on_create = true
}
```

### Application in Environment

```terraform
resource "dokploy_project" "myproject" {
  name = "My Project"
}

resource "dokploy_environment" "staging" {
  name       = "Staging"
  project_id = dokploy_project.myproject.id
}

resource "dokploy_application" "staging_app" {
  name           = "app-staging"
  project_id     = dokploy_project.myproject.id
  environment_id = dokploy_environment.staging.id
  source_type    = "docker"
  docker_image   = "myorg/myapp:staging"
  
  deploy_on_create = true
}
```

### Deploy to Remote Server

```terraform
resource "dokploy_server" "worker" {
  name         = "worker-server"
  ip_address   = "192.168.1.100"
  port         = 22
  username     = "root"
  ssh_key_id   = dokploy_ssh_key.server_key.id
  enable_docker = true
}

resource "dokploy_application" "remote_app" {
  name         = "remote-app"
  project_id   = dokploy_project.myproject.id
  server_id    = dokploy_server.worker.id
  source_type  = "docker"
  docker_image = "nginx:alpine"
  
  deploy_on_create = true
}
```

### Full Application Stack

```terraform
resource "dokploy_project" "webapp" {
  name        = "Web Application"
  description = "Full stack web application"
}

resource "dokploy_database" "postgres" {
  name              = "postgres"
  project_id        = dokploy_project.webapp.id
  type              = "postgres"
  docker_image      = "postgres:16"
  database_name     = "webapp"
  database_user     = "webapp"
  database_password = "secretpassword"
}

resource "dokploy_redis" "cache" {
  name         = "redis-cache"
  project_id   = dokploy_project.webapp.id
  docker_image = "redis:7-alpine"
}

resource "dokploy_application" "api" {
  name              = "api"
  project_id        = dokploy_project.webapp.id
  description       = "Backend API"
  source_type       = "git"
  custom_git_url    = "https://github.com/myorg/api.git"
  custom_git_branch = "main"
  build_type        = "dockerfile"
  
  deploy_on_create = true
}

resource "dokploy_environment_variables" "api_env" {
  application_id = dokploy_application.api.id
  
  variables = {
    DATABASE_URL = "postgresql://webapp:secretpassword@${dokploy_database.postgres.app_name}:5432/webapp"
    REDIS_URL    = "redis://${dokploy_redis.cache.app_name}:6379"
    NODE_ENV     = "production"
  }
}

resource "dokploy_domain" "api" {
  application_id   = dokploy_application.api.id
  host             = "api.example.com"
  port             = 3000
  https            = true
  certificate_type = "letsencrypt"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the application.
- `project_id` (String) The project ID this application belongs to.

### Optional

- `app_name` (String) The app name used for Docker container naming. Auto-generated if not specified.
- `auto_deploy` (Boolean) Enable automatic deployment on Git push.
- `branch` (String) Branch to deploy from.
- `build_args` (String) Build arguments in KEY=VALUE format, one per line.
- `build_path` (String) Build path within the repository for GitHub source.
- `build_type` (String) Build type: dockerfile, heroku_buildpacks, paketo_buildpacks, nixpacks, static, or railpack.
- `command` (String) Custom command to run.
- `cpu_limit` (Number) CPU limit (in millicores, e.g., 1000 = 1 CPU).
- `cpu_reservation` (Number) CPU reservation (in millicores).
- `custom_git_branch` (String) Branch to use for custom Git repository.
- `custom_git_build_path` (String) Build path within the custom Git repository.
- `custom_git_ssh_key_id` (String) SSH key ID for accessing the custom Git repository.
- `custom_git_url` (String) Custom Git repository URL (for source_type 'git').
- `deploy_on_create` (Boolean) Trigger a deployment after creating the application.
- `description` (String) A description of the application.
- `docker_build_stage` (String) Target stage for multi-stage Docker builds.
- `docker_context_path` (String) Docker build context path.
- `docker_image` (String) Docker image to use (for source_type 'docker').
- `dockerfile_path` (String) Path to the Dockerfile.
- `enable_submodules` (Boolean) Enable Git submodules support.
- `env` (String) Environment variables in KEY=VALUE format, one per line.
- `environment_id` (String) The environment ID this application belongs to.
- `github_id` (String) GitHub App installation ID. Required for GitHub source type.
- `memory_limit` (Number) Memory limit in bytes.
- `memory_reservation` (Number) Memory reservation in bytes.
- `owner` (String) Repository owner/organization for GitHub source.
- `password` (String, Sensitive) Password for Docker registry authentication.
- `publish_directory` (String) Publish directory for static builds.
- `registry_url` (String) Docker registry URL.
- `replicas` (Number) Number of replicas to run.
- `repository` (String) Repository name for GitHub source (e.g., 'my-repo').
- `server_id` (String) Server ID to deploy the application to. If not specified, deploys to the default server.
- `source_type` (String) The source type for the application: github, gitlab, bitbucket, git, docker, or drop.
- `trigger_type` (String) Trigger type for deployments: 'push' (default) or 'tag'.
- `username` (String) Username for Docker registry authentication.

### Read-Only

- `id` (String) The unique identifier of the application.

## Import

Import is supported using the following syntax:

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used, for example:

```shell
# Applications can be imported using their ID
terraform import dokploy_application.myapp "application-id-123"
```
