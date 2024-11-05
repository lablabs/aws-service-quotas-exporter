# aws-service-quotas-exporter

## Description

AWS service quotas exporter exposes actual quotas for your AWS accounts and allow you to scrape actual
usage of AWS resources. Base on those two types of data, you can easily build
alert rules to prevent case when you are not able to provision another AWS resources due to reach the limit of AWS quota

## Usage & Configuration

This exporter allows you to scrape metrics via two following approach

### Service quota API

[Viewing service quotas](https://docs.aws.amazon.com/servicequotas/latest/userguide/gs-request-quota.html)

Usually if you are interested to see current value of specific quota metric, you have to use following AWS CLI command:

```bash
aws service-quotas get-service-quota --service-code ec2 --quota-code L-0263D0A3
# or
aws service-quotas get-aws-default-service-quota --service-code route53 --quota-code L-E209CC9F --region us-east-1
```

This can be configured via `quotas` section of configuration file:

Example:
```yaml
quotas:
  # Quota for EC2/ EIPs
  - serviceCode: "ec2"
    quotaCode: "L-0263D0A3"
  # Quota for number of records in route53 zone
  - serviceCode: "route53"
    quotaCode: "L-E209CC9F"
    region: "us-east-1"
    default: true
```

### Usage of AWS resources exports via bash scraping

The pity of AWS service-quota API is that, for most of the resources there is no easy way how to show actual usage of concrete
AWS resources. You have to do it via aws cli/sdk. This exporter try to solve this issue via simple bash script scheduler which
is capable to run specific set of bash scripts and extract actual usage of AWS resources.

Example:

You would like to see usage of AWS EC2 EIPs in your account, cli example:
```bash
aws ec2 describe-addresses --query 'length(Addresses[])'
```
This can be transformed to metric via config:
```yaml
metrics:
    # Name of exporter metric
  - name: "ec2_elastic_ips_usage"
    # Help message for exporter metric
    help: "Usage of ec2 elastic ips"
    # Script used for scraping. It can be inline script or file path to script
    script: "echo \"quota_code=L-0263D0A3,$(aws ec2 describe-addresses --query \'length(Addresses[])\')\""
```

Example above will produce following prometheus metric export:

```
# HELP quota_exporter_ec2_elastic_ips_usage Usage of ec2 elastic ips
# TYPE quota_exporter_ec2_elastic_ips_usage gauge
quota_exporter_ec2_elastic_ips_usage{quota_code="L-0263D0A3"} 4
```

There is requirement regarding format output. For every unique
combination of labels, it has to be one line stdout of your bash script in following csv format

`lable_name_a=label_value,label_name_b=label_value_b,label_name_c=label_value_c,value_for_metric`

## Development

This application requires Go1.22, or later. For common tasks, you can use predefined tasks
via [https://taskfile.dev/](https://taskfile.dev/)

```bash
task --list-all
```

### Docker

#### Build
```
docker build -t ghcr.io/lablabs/aws-service-quotas-exporter:latest .
```
#### Run
Authenticating with AWS credentials:

```
docker run --rm -p 8080:8080 \
    -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
    -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
    -e AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
     ghcr.io/lablabs/aws-service-quotas-exporter:latest
```

Access help:
```
docker run --rm -p 8080:8080 -i ghcr.io/lablabs/aws-service-quotas-exporter:latest --help
```

## Contributing and reporting issues
Feel free to create an issue in this repository if you have questions, suggestions or feature requests.

### Validation, linters and pull-requests

We want to provide high quality code and modules. For this reason we are using
several [pre-commit hooks](.pre-commit-config.yaml) and
[GitHub Actions workflow](.github/workflows/lint.yml). A pull-request to the
master branch will trigger these validations and lints automatically. Please
check your code before you will create pull-requests. See
[pre-commit documentation](https://pre-commit.com/) and
[GitHub Actions documentation](https://docs.github.com/en/actions) for further
details.

## License
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

See [LICENSE](LICENSE) for full details.

    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
