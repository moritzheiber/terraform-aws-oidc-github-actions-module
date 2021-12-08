<!-- BEGIN_TF_DOCS -->
# Terraform AWS OIDC GitHub Actions Module

A module for creating a federated OIDC provider on AWS for dynamically authenticating and authorizing GitHub Actions workflow runs.

# Terraform module documentation

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 3.67.0 |
| <a name="requirement_tls"></a> [tls](#requirement\_tls) | ~> 3.1.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 3.67.0 |
| <a name="provider_tls"></a> [tls](#provider\_tls) | 3.1.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_iam_openid_connect_provider.github](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_openid_connect_provider) | resource |
| [aws_iam_role.federated_auth_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_policy_document.federated_assume_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [tls_certificate.github_actions_oidc_endpoint](https://registry.terraform.io/providers/hashicorp/tls/latest/docs/data-sources/certificate) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_github_actions_oidc_url"></a> [github\_actions\_oidc\_url](#input\_github\_actions\_oidc\_url) | The URL to use for the OIDC handshake | `string` | `"https://vstoken.actions.githubusercontent.com"` | no |
| <a name="input_github_repositories"></a> [github\_repositories](#input\_github\_repositories) | A list of GitHub repositories the OIDC provider should authenticate against. The format is <org/user>/<repository-name> | `set(string)` | `[]` | no |
| <a name="input_role_names"></a> [role\_names](#input\_role\_names) | The set of names for roles that GitHub Actions will be able to assume | `set(string)` | `[]` | no |
| <a name="input_role_path"></a> [role\_path](#input\_role\_path) | The path the created roles are going to live under | `string` | `"/"` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | A key > value map of tags to associate with the resources that are being created | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_roles"></a> [roles](#output\_roles) | The names and ARNs of the roles that were created |
<!-- END_TF_DOCS -->