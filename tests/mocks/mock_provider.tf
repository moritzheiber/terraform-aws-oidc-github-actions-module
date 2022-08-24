provider "aws" {
  access_key                  = "mock_access_key" #tfsec:ignore:aws-misc-no-exposing-plaintext-credentials
  region                      = "eu-central-1"
  s3_use_path_style           = true
  secret_key                  = "mock_secret_key" #tfsec:ignore:general-secrets-sensitive-in-attribute
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    iam = "http://localhost:4566"
  }
}
