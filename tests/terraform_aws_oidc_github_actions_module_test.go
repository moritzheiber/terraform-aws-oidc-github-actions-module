package test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	awssdk "github.com/aws/aws-sdk-go/aws"
)

const allowedRegion = "eu-central-1"
const terraformDir = "../"

func TestGitHubActionsOidcModuleGoodInput(t *testing.T) {
	t.Parallel()

	awsRegion := aws.GetRandomStableRegion(t, []string{allowedRegion}, nil)
	TestKey := "Test"
	TestValue := "Test"
	TestRolePath := "/test/"

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDir,
		Vars: map[string]interface{}{
			"github_repositories": []string{
				"moritzheiber/terraform-aws-oidc-github-actions-module",
				"moritzheiber/some-other-repository",
			},
			"role_names": []string{
				"some-role",
				"another-role",
			},
			"tags": map[string]interface{}{
				TestKey: TestValue,
			},
			"role_path": TestRolePath,
		},
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	})

	t.Run("happy_path", func(t *testing.T) {
		defer terraform.Destroy(t, terraformOptions)
		terraform.InitAndApply(t, terraformOptions)

		arns := terraform.OutputMap(t, terraformOptions, "roles")
		session, _ := session.NewSession(&awssdk.Config{
			Region:           awssdk.String(awsRegion),
			Credentials:      credentials.NewStaticCredentials("test", "test", ""),
			S3ForcePathStyle: awssdk.Bool(true),
			Endpoint:         awssdk.String("http://localhost:4566"),
		})

		client := iam.New(session)

		for name, arn := range arns {
			role, err := client.GetRole(&iam.GetRoleInput{
				RoleName: &name,
			})
			assert.NoError(t, err)
			assert.Equal(t, role.Role.Arn, &arn)
			assert.Equal(t, role.Role.RoleName, &name)
			assert.Equal(t, role.Role.Tags[0].Key, &TestKey)
			assert.Equal(t, role.Role.Tags[0].Value, &TestValue)
			assert.Equal(t, role.Role.Path, &TestRolePath)
		}
	})

	t.Run("without_input", func(t *testing.T) {
		localOptions := terraformOptions
		localOptions.Vars = map[string]interface{}{}

		defer terraform.Destroy(t, localOptions)
		terraform.InitAndApply(t, localOptions)

		arns := terraform.OutputMap(t, terraformOptions, "roles")

		assert.Empty(t, arns)
	})
}

func TestGitHubActionsOidcModuleBadInput(t *testing.T) {
	t.Parallel()

	awsRegion := aws.GetRandomStableRegion(t, []string{allowedRegion}, nil)

	t.Run("bad_https_github_url", func(t *testing.T) {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: terraformDir,
			Vars: map[string]interface{}{
				"github_repositories": []string{
					"https://github.com/moritzheiber/terraform-aws-oidc-github-actions-module",
				},
			},
			EnvVars: map[string]string{
				"AWS_DEFAULT_REGION": awsRegion,
			},
		})

		defer terraform.DestroyE(t, terraformOptions)

		_, err := terraform.InitAndApplyE(t, terraformOptions)
		assert.Error(t, err)
	})

	t.Run("bad_path_variable_suffix", func(t *testing.T) {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: terraformDir,
			Vars: map[string]interface{}{
				"role_path": "/fubar",
			},
			EnvVars: map[string]string{
				"AWS_DEFAULT_REGION": awsRegion,
			},
		})

		defer terraform.DestroyE(t, terraformOptions)

		_, err := terraform.InitAndApplyE(t, terraformOptions)
		assert.Error(t, err)
	})

	t.Run("bad_path_variable_prefix", func(t *testing.T) {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: terraformDir,
			Vars: map[string]interface{}{
				"role_path": "fubar/",
			},
			EnvVars: map[string]string{
				"AWS_DEFAULT_REGION": awsRegion,
			},
		})

		defer terraform.DestroyE(t, terraformOptions)

		_, err := terraform.InitAndApplyE(t, terraformOptions)
		assert.Error(t, err)
	})

	t.Run("bad_path_variable_prefix_suffix", func(t *testing.T) {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: terraformDir,
			Vars: map[string]interface{}{
				"role_path": "fubar",
			},
			EnvVars: map[string]string{
				"AWS_DEFAULT_REGION": awsRegion,
			},
		})

		defer terraform.DestroyE(t, terraformOptions)

		_, err := terraform.InitAndApplyE(t, terraformOptions)
		assert.Error(t, err)
	})

	t.Run("bad_oidc_url", func(t *testing.T) {
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: terraformDir,
			Vars: map[string]interface{}{
				"github_actions_oidc_url": "http://vstoken.actions.githubusercontent.com",
			},
			EnvVars: map[string]string{
				"AWS_DEFAULT_REGION": awsRegion,
			},
		})

		defer terraform.DestroyE(t, terraformOptions)

		_, err := terraform.InitAndApplyE(t, terraformOptions)
		assert.Error(t, err)
	})
}
