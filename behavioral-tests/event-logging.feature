Feature: event logging

  Lifecycle events such as application start, stop or deploy should be tracked.

  Scenario: Lifecycle event tracking caches event log entries on disk if HTTP server is disabled, and sends it when the server is back online.
    Given I have an empty Dockerfile application
    And the event log is empty
    And the configuration is:
      | apiEndpointUrl | http://10.0.0.1:23232/apiTest |
    And the API delivery http server is disabled
    When I deploy the application as "test"
    Then I expect 1 event log entry on disk

    # when the API server is online, but does not respond with a 2xx status code, the event stays cached
    When the API delivery http server is available at port 23232 for at most 10 seconds and 1 request
    And the API delivery http server always responds with status code 503
    And I call dokku "collectMetrics"
    Then I expect 1 event log entries on disk

    # when the API is back online, the event gets delivered
    When the API delivery http server is available at port 23232 for at most 10 seconds and 1 request
    And I call dokku "collectMetrics"
    Then I expect 0 event log entries on disk
    And the API delivery http server received request 1 with the following JSON at "/apiTest/log":
      | event.application | equals        | test                                                                    |
      | event.serverName  | matches regex | (?:[0-9]{1,3}\.){3}[0-9]{1,3}                                           |
      | event.message     | equals        | Deployment successful! (Image Tag: )                                    |
      | event.timestamp   | is a date     |                                                                         |
      | event.uuid        | matches regex | [0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12} |