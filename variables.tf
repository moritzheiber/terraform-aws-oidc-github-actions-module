variable "github_actions_oidc_url" {
  type        = string
  description = "The URL to use for the OIDC handshake"
  default     = "https://vstoken.actions.githubusercontent.com"

  validation {
    condition     = substr(var.github_actions_oidc_url, 0, 5) == "https"
    error_message = "The OIDC URL needs to start with https."
  }
}

variable "github_repositories" {
  type        = set(string)
  description = "A list of GitHub repositories the OIDC provider should authenticate against. The format is <org/user>/<repository-name>"
  default     = []

  validation {
    condition     = alltrue([for repo in var.github_repositories : substr(repo, 0, 4) != "http"])
    error_message = "The repositories must not have a http(s):// prefix. The format is <org/user>/<repository-name>."
  }
}

variable "role_names" {
  type        = set(string)
  description = "The set of names for roles that GitHub Actions will be able to assume"
  default     = []
}

variable "role_path" {
  type        = string
  description = "The path the created roles are going to live under"
  default     = "/"

  validation {
    condition     = substr(var.role_path, 0, 1) == "/" && substr(strrev(var.role_path), 0, 1) == "/"
    error_message = "The path for the role must start and end in a slash (/)."
  }
}

variable "tags" {
  type        = map(string)
  description = "A key > value map of tags to associate with the resources that are being created"
  default     = {}
}
