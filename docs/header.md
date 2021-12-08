# Terraform AWS OIDC GitHub Actions Module

A module for creating a federated OIDC provider on AWS for dynamically authenticating and authorizing GitHub Actions workflow runs.

## Setting up the OIDC AWS provider

Add the module to one of your Terraform configurations to create an OIDC provider and one or more roles that can be assumed via the provider. The names and ARNs of the created roles will be provided in the `roles` output of the module. You will need one or more names for GitHub repositories that GitHub Actions should be allowed to assume the roles from in order to configure the module:

```hcl
module "oidc_auth" {
  source = "github.com/moritzheiber/terraform-aws-oidc-github-actions-module"

  github_repositories = toset(["my-org/my-repository"])
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

```yaml
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

You can get the ARN of any of the roles created via the `roles` output of the OIDC module. In keeping with our previous example, the ARN for the `some-role` role would be accessible via `module.oidc_auth.roles["some-role"]`.

You will probably want to add the ARN for `role-to-assume` as [a GitHub Actions secret](https://docs.github.com/en/actions/security-guides/encrypted-secrets) instead of writing it directly into the workflow YAML.

# Terraform module documentation
