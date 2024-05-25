# AWS Serverless Static Site Gateway

This project provides a serverless architecture for hosting static websites in S3 using AWS Lambda and API Gateway, without the need for enabling S3 Public Access or using Amazon CloudFront. The setup serves static content stored in an S3 bucket through a Lambda function, which acts as a reverse proxy.

## Prerequisites

Before starting the setup, ensure you have the following:

- **AWS Account**: An active account with AWS that can access Lambda, S3, and API Gateway services.
- **Golang Compilation Environment**: A setup with Go installed, enabling you to compile the source code for deployment.

## Setup

### Step 1: Prepare the Lambda Function

1. **Create the Lambda Function**:

    - Navigate to the AWS Lambda console and create a new Lambda function.

    - Choose the runtime as **Amazon Linux 2 arm64**.

    - Set the handler to bootstrap. This configuration is necessary because you are using a custom runtime environment.

2. **Compile the Code**:
    - Use the make package command in your local Go environment to compile the Lambda function code. This command should create a function.zip file in the build directory.

3. **Upload the Compiled Function**:
    - Upload the function.zip file from the build directory to your Lambda function.

4. **Configure Environment Variables**:
    - Set necessary environment variables such BUCKET_NAME to specify the S3 bucket and AWS region your Lambda will interact with.

5. **Set IAM Role Permissions**:
    - Ensure the IAM role associated with your Lambda function has the necessary permissions to access S3. This typically includes actions like s3:GetObject within the policy attached to the role.

### Step 2: S3 Bucket Configuration

- **Create an S3 Bucket**: If you don’t already have one, create a new S3 bucket where your static files and SPA will reside.

- **Upload Files**: Populate your bucket with the static and SPA files like index.html, JavaScript, CSS, and images.

### Step 3: Configure API Gateway

1. **Create API Gateway**:
    - Open the AWS API Gateway console and create a new REST API.

2. **Configure GET Method for Root**:

    - Create a new resource corresponding to the root path /.

    - Attach a GET method to this resource that triggers your Lambda function. This setup handles requests directed at your root URL.

3. **Setup {proxy+} with ANY Method**:
    - Add a {proxy+} resource under the root to handle all other paths and HTTP methods. Set this to trigger the same Lambda function.

4. **Enable Binary Support**:
    - Go to the API settings and add */* to the binary media types list to handle all content types as binary data.

5. **Deploy the API**:
    - Deploy your API to make it accessible, selecting or creating a new deployment stage.
6. **Configure Custom Domain Name**
    - Navigate to the “Custom Domain Names” section in the API Gateway console.
    - Click “Create”, and specify your desired domain name.
    - Link the domain name with a certificate from AWS Certificate Manager (ACM). Note: Ensure you have a certificate ready in ACM for the domain you intend to use.
    - Set up a base path mapping that connects your custom domain to the deployed API stage.

After deploying the API Gateway and Lambda function, you can access your SPA and static files via the URL provided by API Gateway. Use this URL to access your site:

```
https://example.com/
```

## Troubleshooting

**Images not displaying**: Verify that binary media types are properly set in API Gateway.