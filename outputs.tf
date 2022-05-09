output "roles" {
  description = "The names and ARNs of the roles that were created"
  value       = { for r in var.role_names : r => aws_iam_role.federated_auth_role[r].arn }
}

output "github_actions_thumbprint" {
  description = "The thumbprint of the TLS certificate used for the OIDC endpoint at GitHub Actions"
  value       = local.thumbprint
}
