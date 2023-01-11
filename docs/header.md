[![Module Releases](https://img.shields.io/badge/dynamic/json?color=%237b42bc&label=Release&query=version&url=https%3A%2F%2Fregistry.terraform.io%2Fv1%2Fmodules%2Fmoritzheiber%2Foidc-github-actions-module%2Faws&logo=terraform&style=for-the-badge)](https://registry.terraform.io/modules/moritzheiber/oidc-github-actions-module/aws/latest) ![Module Downloads](https://img.shields.io/badge/dynamic/json?color=%237b42bc&label=Downloads&query=data.attributes.total&url=https%3A%2F%2Fregistry.terraform.io%2Fv2%2Fmodules%2Fmoritzheiber%2Foidc-github-actions-module%2Faws%2Fdownloads%2Fsummary&logo=terraform&style=for-the-badge)

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

*Note: Usually you will want to use a specific version of the module by using the `version` attribute.*

Continue with assigning permissions to these roles:

```hcl
resource "aws_iam_policy" "policy" {
  name        = "some_policy"
  path        = "/"
  description = "Some policy"

  policy = jsonencode({
    Version = "2012-10-17"
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

You can get the ARN of any of the roles created via the `roles` output of the OIDC module. In keeping with our previous example, the ARN for the `some-role` role would be accessible via `module.oidc_auth.roles["some-role"]`. Be sure that `<some-region>` matches the region you used earlier to provision your Terraform code or otherwise you'll run into authentication issues.

You will probably want to add the ARN for `role-to-assume` as [a GitHub Actions secret](https://docs.github.com/en/actions/security-guides/encrypted-secrets) instead of writing it directly into the workflow YAML.

## Troubleshooting

#### "OpenIDConnect provider's HTTPS certificate doesn't match configured thumbprint"

The AWS OIDC provider _requires_ you to store the HTTP TLS thumbprint of any and all OIDC endpoints it is going to use to identify providers. In this case that's the endpoint for GitHub Actions, with its default URL at `https://token.actions.githubusercontent.com` (see `variables.tf`).

TLS certificates change all the time, which is why, occasionally, you will have workflows fail because the thumbprint you stored in the AWS OIDC provider when running this module initially doesn't match "the reality" anymore (i.e. the endpoints HTTP TLS certificate has changed and therefore also its thumbprint). AWS will refuse to operate against an OIDC provider for which it doesn't have a correct thumbprint stored to prevent malicious actors from spoofing a seemingly valid HTTP TLS certificate to gain access to your AWS Account.

You will have to _manually verify that the new TLS certificate for the GitHub Actions endpoint is valid and can be trusted_ and then re-run this module afterwards.

For convenience, the module automatically fetches the latest thumbprint, and it will explicitely tell you that it's going to change the OIDC provider's thumbprint. After having re-run the module the newer thumbprint should be stored in the OIDC provider's definitions and your workflows should run again.

##### Steps to resolve this error

1. Fetch the current thumbprint from the module's output:

```console
$ terraform output [potential-module-prefix].github_actions_thumbprint
```

The thumbprint will look something like this: `15e29108718111e59b3dad31954647e3c344a231` (it's the `sha1` thumbprint, for the curious)

2. Make sure you are on a trusted network (e.g. no suspicious actors between you and GitHub's infrastructure) and run the following command:

```console
$ openssl s_client -connect token.actions.githubusercontent.com:443 | openssl x509 -noout -fingerprint -sha1
depth=2 C = US, O = DigiCert Inc, OU = www.digicert.com, CN = DigiCert Global Root CA
verify return:1
depth=1 C = US, O = DigiCert Inc, CN = DigiCert TLS RSA SHA256 2020 CA1
verify return:1
depth=0 C = US, ST = California, L = San Francisco, O = "GitHub, Inc.", CN = *.actions.githubusercontent.com
verify return:1
sha1 Fingerprint=15:E2:91:08:71:81:11:E5:9B:3D:AD:31:95:46:47:E3:C3:44:A2:31
```

This is important: If the string from the first step doesn't match the string under `sha1 Fingerprint` (separated by colons) there are two scenarios: either the certificate was rotated (likely) or somebody is trying to highjack your connection to GitHub Actions.

Should the thumbprints match it could've been [a temporary issue on AWS](https://health.aws.amazon.com/health/status) (it's the cloud after all) or [GitHub might be having problems](https://www.githubstatus.com/). It's unlikely to be something else, especially not this module. I would try to re-run the workflow(s), either now or at a later point in time.

3. Verify that you're dealing with a newer certificate:

```console
$ openssl s_client -connect token.actions.githubusercontent.com:443 | openssl x509 -noout -dates
depth=2 C = US, O = DigiCert Inc, OU = www.digicert.com, CN = DigiCert Global Root CA
verify return:1
depth=1 C = US, O = DigiCert Inc, CN = DigiCert TLS RSA SHA256 2020 CA1
verify return:1
depth=0 C = US, ST = California, L = San Francisco, O = "GitHub, Inc.", CN = *.actions.githubusercontent.com
verify return:1
notBefore=Jan 11 00:00:00 2022 GMT
notAfter=Jan 11 23:59:59 2023 GMT
```

As you can see, under `notBefore` and `notAfter` are two dates. It is likely that the `notBefore` date is somewhere in the not too distant past (a day to a week), which would indicate that it was recently rotated.

4. This is where you have to ask yourself whether you trust the newer certificate the endpoint is presenting to you. Because a newer certificate could also mean somebody generated a "fake" certificate recently and is trying to use a [MITM](https://en.wikipedia.org/wiki/Man-in-the-middle_attack) attack to harvest your (temporary) AWS credentials. There are other methods to verify an endpoints authenticity (checking on [OCSP revocations](https://www.certificatetools.com/ocsp-checker), [SSLLabs](https://www.ssllabs.com/ssltest/index.html)), but a more in-depth defense against TLS-based attacks is beyond the scope of this document.

5. Adjust the thumbprint in the AWS OIDC provider configuration

If you're sure that the newer certificate was issued by a trusted authority (GitHub, DigiCert or some other trusted source) you can simply re-run the Terraform code this module is used from to replace the old thumbprint with the newer version. Terraform will automatically fetch the latest thumbprint and add it to your configuration

Afterwards your workflows should run without authentication issues again.

# Terraform module documentation
