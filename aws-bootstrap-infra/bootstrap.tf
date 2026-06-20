provider "aws" {
  region = "us-east-1"
}

resource "aws_iam_openid_connect_provider" "github" {
  url             = "https://token.actions.githubusercontent.com" 
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = ["6938fd4d98bab03faadb97b34396831e3780aea1"] 
}

resource "aws_iam_role" "github_actions" {
  name               = "GitHubActionsEKSRole"
  assume_role_policy = file("${path.module}/iam/github-oidc-trust-policy.json")
}


resource "aws_iam_policy" "github_actions_policy" {
  name   = "GitHubActionsEKSPolicy"
  policy = file("${path.module}/iam/github-actions-policy.json")
}

resource "aws_iam_role_policy_attachment" "github_actions_attach" {
  role       = aws_iam_role.github_actions.name
  policy_arn = aws_iam_policy.github_actions_policy.arn
}
