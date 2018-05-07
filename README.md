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


#E CS/Docker Integration

Find the entrypoint location in the original Dockerfile or just inspect. In this case /usr/local/bin/docker-entrypoint.sh

## 1. Create Custom Entrypoint.sh
```bash
#!/bin/bash

[[ -n ${SD_DNSRECORD} && -n ${SD_DNSZONE} ]] && /usr/local/bin/poorman-sd -r "${SD_DNSRECORD}" -z "${SD_DNSZONE}"

/usr/local/bin/docker-entrypoint.sh "$*"
```

2. Add our tool and custom entrypoint
# Add tool for service discovery
```
COPY poorman-sd /usr/local/bin
RUN chmod 755 /usr/local/bin/poorman-sd

# add Custom entrypoint
COPY customdocker-entrypoint.sh /usr/local/bin
RUN chmod 755 /usr/local/bin/customdocker-entrypoint.sh

# point entrypoint to our entrypoint
ENTRYPOINT ["/usr/local/bin/customdocker-entrypoint.sh"]
```

3. Add the variable to the ECS or docker task definition
```yaml
  "environment": [
      {
        "name": "SD_DNSRECORD",
        "value": "test.myzone.test"
      },
      {
        "name": "SD_DNSZONE",
        "value": "AWSZONEID32"
      }
  ]
```
