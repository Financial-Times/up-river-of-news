
#up-slack-integration
https://hooks.slack.com/services/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

You have two options for sending data to the Webhook URL above:
Send a JSON string as the payload parameter in a POST request
Send a JSON string as the body of a POST request
For a simple message, your JSON payload could contain a text property at minimum. This is the text that will be posted to the channel.
A simple example:
payload={"text": "This is a line of text in a channel.\nAnd this is another line of text."}
This will be displayed in the channel as:


To display a richly-formatted message attachment in Slack, you can use the same JSON payload as above, but add in an attachments array. Each element of this array is a hash containing the following parameters:
{
	"fallback": "Required text summary of the attachment that is shown by clients that understand attachments but choose not to show them.",

	"text": "Optional text that should appear within the attachment",
	"pretext": "Optional text that should appear above the formatted data",

	"color": "#36a64f", // Can either be one of 'good', 'warning', 'danger', or any hex color code

	// Fields are displayed in a table on the message
	"fields": [
		{
			"title": "Required Field Title", // The title may not contain markup and will be escaped for you
			"value": "Text value of the field. May contain standard message markup and must be escaped as normal. May be multi-line.",
			"short": false // Optional flag indicating whether the `value` is short enough to be displayed side-by-side with other values
		}
	]
}				
Please note that the fallback field is required, and is displayed whenever message attachments cannot be shown (ie. mobile notifications, desktop notifications, IRC).

Example
This example will post a detailed message attachment as though it were sent from a task management service:
{
   "attachments":[
      {
         "fallback":"New open task [Urgent]: <http://url_to_task|Test out Slack message attachments>",
         "pretext":"New open task [Urgent]: <http://url_to_task|Test out Slack message attachments>",
         "color":"#D00000",
         "fields":[
            {
               "title":"Notes",
               "value":"This is much easier than I thought it would be.",
               "short":false
            }
         ]
      }
   ]
}				
