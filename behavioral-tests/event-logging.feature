Feature: event logging

  Lifecycle events such as application start, stop or deploy should be tracked.

  Scenario: Lifecycle event tracking caches event log entries on disk if HTTP server is disabled, and sends it when the server is back online.
    Given I have an empty node.js application
    And the event log is empty
    And the configuration is:
      | apiEndpointUrl | http://10.0.0.1:23232/apiTest |
    And the API delivery http server is disabled
    When I deploy the application as "test"
    Then I expect 1 event log entry on disk

    # when the API is back online, the event gets delivered
    When the API delivery http server is available at /apiTest
    And I call dokku "collectMetrics"
    Then I expect 0 event log entries on disk
    And the API delivery http server received the following JSON at /event:
      | event.application | equals        | test          |
      | event.timestamp   | is a date     |               |
      | event.uuid        | matches regex | [a-z0-9-]{40} |