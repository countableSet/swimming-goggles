# swimming-goggles 

AWS lambda to update S3 bucket policy with Cloudflare IP addresses

#### Cloudwatch schedule

Create a schedule trigger with the following 'cron(0 17 ? * 1 *)' which runs everything week on Sunday at 17:00 GMT.
Input should look something like `{"buckets":["logs.countableset.com"]}` for the event json.

#### IAM Policies

Role policy for lambda:
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "s3:PutBucketPolicy",
                "s3:DeleteBucketPolicy",
                "s3:GetBucketPolicy"
            ],
            "Resource": "arn:aws:s3:::logs.countableset.com"
        }
    ]
}
```

User policy for lambda access:
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

#### Environment Variables

```
export AWS_REGION="us-west-2"
export AWS_LAMBDA_S3_BUCKET_POLICY_ROLE_ARN="arn:aws:iam::XXYYZZ:role/lambda-s3-bucket-policy"
```