# aws-sinkdrain

aws sink drain

## How to use

### Pre-req

1. Set `AWS_ACCOUNT` and `AWS_REGION` in `env/**/xyz.env`.
2. Source the env file.
   - `source env/xyz.env`
3. Make sure AWS credentials are configured using any [credential chain](https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/credentials.html#credentials-default)
   - `aws sts get-caller-identity`
4. Make sure `docker` is running.
   - `docker info`

### Creating / Updating

1. Launch the infra
   - `make launch`

### Delete

1. Delete the infra
   - `make destroy`

### Check

Once deployed, run `POST` call against api gateway endpoint at route `/sinkdrain`.

ex) `https://<something>.execute-api.<something>.amazonaws.com/sinkdrain`

Wait couple minutes for CloudWatch logs to sync, then check CloudWatch console in AWS UI.

ex) `/aws/lambda/aws-sinkdrain` Log group
