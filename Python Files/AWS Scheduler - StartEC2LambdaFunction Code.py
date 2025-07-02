import boto3

def lambda_handler(event, context):
    ec2 = boto3.client('ec2')
    sns = boto3.client('sns')
    
    instance_id = 'i-0202015e2be260647'
    topic_arn = 'arn:aws:sns:us-east-2:212460019729:OHSTE001A-ALS001A-Schedule-Notifications'
    
    try:
        print(f"Attempting to start EC2 instance: {instance_id}")
        ec2.start_instances(InstanceIds=[instance_id])
        message = f"EC2 instance {instance_id} started successfully."
        print(message)

        sns.publish(
            TopicArn=topic_arn,
            Message=message,
            Subject="EC2 Scheduler - Start Success"
        )
        print("SNS notification sent.")
        print("EC2 start and notification complete.")  # Added line
        
        return {
            'statusCode': 200,
            'body': message
        }

    except Exception as e:
        message = f"Error starting EC2 instance {instance_id}: {str(e)}"
        print(message)

        sns.publish(
            TopicArn=topic_arn,
            Message=message,
            Subject="EC2 Scheduler - Start FAILED"
        )
        return {
            'statusCode': 500,
            'body': message
        }