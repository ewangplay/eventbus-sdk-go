# EventBus GO SDK Package

```
import "github.com/ewangplay/eventbus-sdk-go"
```

## How to use ?

[Client example](./examples/example01.go)

### HTTP API

- Request URL: `POST http://IP:PORT/v1/event`

- Request Body:

    ```
	{
		"subject": "payment",
		"type": "notifier",
		"retry_count": 100,
		"retry_interval": 2,
		"retry_timeout": 600,
		"retry_policy": 2,
		"target_url": "http://127.0.0.1:8091/v1/test",
		"body": {
			"subject": "bocpoints-704fd8a2-dc25-4076-92df-d434e3d7e31f",
			"eventType": 1,
			"chaincodeID": "50aa2bc9277c23a418f843d817bb545aafc87af86b2dc6bc53421b45e733732b5dc6bd484da05a741b8323aa92646c8510972b35de9f623097b0645db9afda81",
			"chaincodeName": "bocpoints",
			"chaincodeEventName": "PaymentNotify",
			"txID": "704fd8a2-dc25-4076-92df-d434e3d7e31f",
			"txType": 2 
		}
	}
    ```

    - subject: event subject
    - type: event type, supported two event types: 1. notifier, 2. queuer.
    - retry_count: The max retry count
    - retry_interval: The retry interval
    - retry_timeout: The retry timeout 
    - retry_policy: 1: Stop to retry when the max retry count reached; 2: Stop to retry when the retry timeout reached 
    - body: the actual data to be published

- Response Result:

    ```
    {
		"error_code" : 0,
		"message": "error description",
		"event_id": "123",
		"create_time" : "2017-08-02 08:17:35"
    }
    ```

    - error_code: 0 indicates success, none-zero indicates failure.
    - message: Represents the error information when error_code is none-zero
    - event_id: if `error_code` is zero, return the event global unique ID
    - create_time: create time of the event
