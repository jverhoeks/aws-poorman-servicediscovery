# Poor man's Service Discovery

Updates a Route53 based on the local ipv4.
Works on EC2 and ECS


## Running

Run:

`poorman-sd -r test.domain.com. -z AWSZONEID`

With CNAME

`poorman-sd -r test.domain.com. -z AWSZONEID -t CNAME`

Provide IP, don't use metadeta

`poorman-sd -r test.domain.com. -z AWSZONEID -v 1.2.3.4`


## Using

Include the command in the user-data or the dockerentrypoint.sh to update the dns record of the service.


## Building

Install AWS-SDK:  `go get -u github.com/aws/aws-sdk-go`

Run: `./build.sh`

The script wil build linux and darwin executables.
If UPX is installed it'll compress the executables to redux the size.

## AWS IAM Profile

The following permission is required to update the record
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "route53:ChangeResourceRecordSets"
            ],
            "Resource": "arn:aws:route53:::hostedzone/AWSZONEID"
        }
    ]
}
```
