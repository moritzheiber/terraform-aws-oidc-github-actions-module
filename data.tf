data "tls_certificate" "github_actions_oidc_endpoint" {
  url = var.github_actions_oidc_url
}
