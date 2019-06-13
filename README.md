# insight_lambda
Follow the instructions below to send logs stored on AWS S3 to InsightOPS.

All source code and dependencies can be found on the [insight_lambda Github page](https://github.com/Tweddle-SE-Team/insight_lambda).

###### Example use cases:
* Forwarding AWS ELB/ALB/CloudFront logs
  * (make sure to set ELB/ALB/CloudFront to write logs every 5 minutes)
  * When forwarding these logs, the script will format the log lines according to JSON spec to make them easier to analyze

## Obtain log token
1. Log in to your InsightOPS account

2. Add a new [token based log](https://insightops.help.rapid7.com/docs/token-tcp)

## Deploy the script to AWS Lambda using AWS Console
1. Create a new Lambda function

2. Choose the Python blueprint for S3 objects

   ![Choose Blueprint](https://raw.githubusercontent.com/logentries/le_lambda/master/doc/step2.png)

3. Configure triggers:
   * Choose the bucket log files are being stored in
   * Set event type "Object Created (All)"
   * Tick "Enable Trigger" checkbox

4. Configure function:
   * Give your function a name
   * Set runtime to Go1.x

5. Upload function code:
   * Choose "Upload a .ZIP file" in "Code entry type" dropdown and upload the archive from Github Releases

6. Set Environment Variables:
   * Token value should match UUID provided by InsightOPS UI or API
   * Region should be that of your InsightOPS account - currently only ```eu```

   | Key       | Value      |
   |-----------|------------|
   | INSIGHT_API_REGION | <region> |
   | INSIGHT_API_KEY | token uuid |

7. Lambda function handler and role
   * Change the "Handler" value to ```insight_lambda```
   * Choose "Create a new role from template" from dropdown and give it a name below.
   * Leave "Policy templates" to pre-populated value

8. Advanced settings:
   * Set memory limit to a high enough value to facilitate log parsing and sending - adjust to your needs
   * Set timeout to a high enough value to facilitate log parsing and sending - adjust to your needs
   * Leave VPC value to "No VPC" as the script only needs S3 access
     * If you choose to use VPC, please consult [Amazon Documentation](http://docs.aws.amazon.com/lambda/latest/dg/vpc.html)

9. Enable function:
   * Click "Create function"

## Gotchas:
   * The "Test" button execution in AWS Lambda will **ALWAYS** fail as the trigger is not provided by the built in test function. In order to verify, upload a sample file to source bucket
