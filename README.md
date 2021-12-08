<!-- BEGIN_TF_DOCS -->
# Terraform AWS OIDC GitHub Actions Module

A module for creating a federated OIDC provider on AWS for dynamically authenticating and authorizing GitHub Actions workflow runs.

## Setting up the OIDC AWS provider

Add the module to one of your Terraform configurations to create an OIDC provider and one or more roles that can be assumed via the provider. The names and ARNs of the created roles will be provided in the `roles` output of the module. You will need one or more names for GitHub repositories that GitHub Actions should be allowed to assume the roles from in order to configure the module:

```hcl
module "oidc_auth" {
  source = "github.com/moritzheiber/terraform-aws-oidc-github-actions-module"

  github_repositories = toset(["my-org/my-repository])
  role_names          = toset(["some-role"])
}

output "github_actions_roles" {
    value = module.oidc_auth.roles
}
```

Continue with assigning permissions to these roles:

```hcl
resource "aws_iam_policy" "policy" {
  name        = "some_policy"
  path        = "/"
  description = "Some policy"

  policy = jsonencode({
    Statement = [
      {
        Action = [
          "ec2:Describe*",
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "attachment" {
  for_each = module.oidc_auth.roles
  role       = each.key
  policy_arn = aws_iam_policy.policy.arn
}
```

## Setting up GitHub Actions

[AWS provides a "native" GitHub Actions action](https://github.com/aws-actions/configure-aws-credentials) to enable you to use the configured OIDC provider. Just add the following two bits to any job for a GitHub repository you passed under `github_repository` to the module previously and you should be good to go:

```
jobs:
    some-job:
        # [...]
        permissions:
          id-token: write
          contents: read
        # [...]
        steps:
            # [...]
            - uses: aws-actions/configure-aws-credentials@v1
                with:
                  role-to-assume: <ARN-of-the-one-or-any-of-the-roles-created-by-the-module>
                  aws-region: <some-region>
            # [...]
            # Any step beyond the last one now has access to your AWS account with the permissions
            # you assigned via the policy associated with the role you want to assume
```

\_Note: You can get the ARN of any of the roles created via the `roles` output of the OIDC module. In keeping with our previous example, the ARN for the `some-role` role would be accessible via `module.oidc_auth.roles["some-role"]`.\_

You will probably want to add the ARN for `role-to-assume` as [a GitHub Actions secret](https://docs.github.com/en/actions/security-guides/encrypted-secrets) instead of writing it directly into the workflow YAML.

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
| <a name="input_github_actions_oidc_url"></a> [github\_actions\_oidc\_url](#input\_github\_actions\_oidc\_url) | The URL to use for the OIDC handshake | `string` | `"https://token.actions.githubusercontent.com"` | no |
| <a name="input_github_repositories"></a> [github\_repositories](#input\_github\_repositories) | A list of GitHub repositories the OIDC provider should authenticate against. The format is <org/user>/<repository-name> | `set(string)` | `[]` | no |
| <a name="input_role_names"></a> [role\_names](#input\_role\_names) | The set of names for roles that GitHub Actions will be able to assume | `set(string)` | `[]` | no |
| <a name="input_role_path"></a> [role\_path](#input\_role\_path) | The path the created roles are going to live under | `string` | `"/"` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | A key > value map of tags to associate with the resources that are being created | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_roles"></a> [roles](#output\_roles) | The names and ARNs of the roles that were created |
<!-- END_TF_DOCS -->