locals {
  have_repositories = length(var.github_repositories) > 0
  plain_oidc_url    = trimprefix(var.github_actions_oidc_url, "https://")
}

resource "aws_iam_openid_connect_provider" "github" {
  count = local.have_repositories ? 1 : 0

  url = var.github_actions_oidc_url
  client_id_list = [
    "sts.amazonaws.com"
  ]

  thumbprint_list = [
    data.tls_certificate.github_actions_oidc_endpoint.certificates.0.sha1_fingerprint
  ]

  tags = var.tags
}

data "aws_iam_policy_document" "federated_assume_policy" {
  count = local.have_repositories ? 1 : 0

  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]
    effect  = "Allow"

    principals {
      type = "Federated"
      identifiers = [
        aws_iam_openid_connect_provider.github[0].arn
      ]
    }

    condition {
      test     = "StringLike"
      variable = "${local.plain_oidc_url}:sub"

      values = [for repo in var.github_repositories : "repo:${repo}:*"]
    }
  }
}

resource "aws_iam_role" "federated_auth_role" {
  for_each = var.role_names

  name        = each.key
  path        = var.role_path
  description = "Federated identity role for GitHub Actions"

  assume_role_policy = data.aws_iam_policy_document.federated_assume_policy[0].json

  tags = var.tags
}
