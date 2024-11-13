# SNS x Slack Integration

## SNS Topic

Add topic to template

```yml
LogToSlackTopic:
  Type: AWS::SNS::Topic
  Properties:
    TopicName: <service>-slack-integration
```

Add to environment variables
```yml
LOG_TO_SLACK_TOPIC_ARN: !Ref LogToSlackTopic
```

Export topic
```yml
Log<Service>ToSlackTopicArn:
  Description: LogToSlackTopic ARN
  Value: !Ref LogToSlackTopic
  Export:
    Name: Log<Service>ToSlackTopicArn
```

Add policy to lambda Function
```yml
SNSPublishMessagePolicy:
  TopicName: !GetAtt LogToSlackTopic.TopicName
```

## Slack Channel
Create a channel
- Create a `public` channel with format `#logs-service-name`
- Generate an email for the channel with these steps
  - Right click channel
  - Click `View Channel Details`
  - Go to `Integrations` tab
  - Click `Send emails to this channel`
  - Copy the generated email

Export channel email from template
```yml
Parameters:
  SlackChannelEmail:
    Type: String
    Default: logs-<service-name>-aaaaosa4mzmuclukfslegmv3se@scrambledeggs.slack.com

Logs<Service>ChannelEmail:
  Description: logs-<service-name> channel
  Value: !Ref SlackChannelEmail
  Export:
    Name: Logs<Service>ChannelEmail
```

## Lambda Function (Application Level)
Make sure to update the `logs` package from `booky-go-common` to the most recent version. Version hash would be
```
# go.mod

github.com/scrambledeggs/booky-go-common/logs v0.0.0-20241112053933-dab1dc30b54e
```

Add flag to logging. 
```go
logs.Error("note", map[string]string{"naknang": "patatas"}, logs.TO_SLACK)
```

_Note: Flag is available to all log types_

## Subscription
Add subscription resource to subscriptions manager's corresponding template

```yml
LogToSlackLogs<Service>ChannelSubscription:
  Type: AWS::SNS::Subscription
  Properties:
    Protocol: email
    TopicArn:
      Fn::ImportValue: Log<Service>ToSlackTopicArn
    Endpoint:
      Fn::ImportValue: Logs<Service>ChannelEmail
```

Upon deployment, the subscription resource will be created and will send this message to the registered slack channel

![image](https://github.com/user-attachments/assets/70a1d223-0665-4f52-b4d5-67942b6886a1)

Upon clicking the given subscription URL

![image](https://github.com/user-attachments/assets/68f52e19-88c0-456b-b66a-5f8bc7cf2d36)

_Happy Logging!_


