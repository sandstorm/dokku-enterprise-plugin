Feature: I can get information about all applications stored in the cloud
  Background:
    Given an empty folder "/tmp/storage" exists
    And an empty folder "/tmp/storage/mybucket" exists
    And the cloud configuration is:
      | localStoragePath | /tmp/storage/                                    |
      | type             | local                                            |
      | storageBucket    | /tmp/storage/mybucket                            |
      | encryptionKey    | foofoofoofooofoofoofoofoofooofoofoofoofoofooofoo |

    And I have an empty Dockerfile application
    And I call dokku "apps:destroy test --force"
    And I call dokku "mariadb:destroy test42 --force"
    And an empty folder "/tmp/test" exists
    And a file "/tmp/test/storage.txt" is created with contents:
    """
    Hallo Welt - we must survive!
    """
    And I call dokku "apps:destroy copy --force"
    And I call dokku "mariadb:destroy copy42 --force"

    And I call dokku "apps:create test"
    And I deploy the application as "test"
    And I call dokku "cloud:backup test"
    And I call dokku "cloud:backup test"
    And I call dokku "cloud:createAppFromCloud copy test"
    And I call dokku "cloud:backup copy"

  Scenario: List all existing applications
    Given I call dokku "cloud:list"
    Then I get back a table with content:
      | NAME | VERSIONS | LATEST       |
      | copy | 1        | copy__.*__.* |
      | test | 2        | test__.*__.* |

  Scenario: List information about versions of an application
    Given I call dokku "cloud:info test"
    Then I get back a table with content:
      | VERSION      | DATE                                | CODE SIZE \(KB\) | DATA SIZE \(KB\) |
      | test__.*__.* | \d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2} | 0\.\d{3}         | \d{1,2}\.\d{3}   |
      | test__.*__.* | \d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2} | 0\.\d{3}         | \d{1,2}\.\d{3}   |