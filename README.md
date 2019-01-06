# Lambda S3 Bucket Policy for Cloudflare IP Address

AWS Lambda to update S3 bucket policy IP addresses to cloudflare.

#### Cloudwatch schedule

Create a schedule trigger with the following 'cron(0 17 ? * 1 *)' which runs everything week on Sunday at 17:00 GMT.

#### IAM Policy

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "lambda:CreateFunction",
                "lambda:UpdateFunctionCode",
                "lambda:ListFunctions",
                "lambda:InvokeFunction",
                "lambda:ListVersionsByFunction",
                "lambda:GetFunction",
                "lambda:UpdateFunctionConfiguration",
                "lambda:DeleteFunction",
                "iam:GetRole",
                "iam:PassRole"
            ],
            "Resource": "*"
        }
    ]
}
```