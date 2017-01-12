Feature: I can export the Configuration Manifest to the cloud

  Scenario: Application with a database
    Given an empty folder "/tmp/storage" exists
    Given an empty folder "/tmp/storage/mybucket" exists
    And the cloud configuration is:
      | localStoragePath | /tmp/storage/                                    |
      | type             | local                                            |
      | storageBucket    | /tmp/storage/mybucket                            |
      | encryptionKey    | foofoofoofooofoofoofoofoofooofoofoofoofoofooofoo |

    And I have an empty Dockerfile application
    And I call dokku "apps:destroy test --force"
    And I call dokku "apps:create test"
    And I call dokku "mariadb:create test42"
    And I call dokku "mariadb:link test42 test"
    And I call dokku "config:set test FOO=bar"
    And I call dokku "storage:mount test /tmp/test/:/b"
    And I deploy the application as "test"

    When I call dokku "cloud:backup test"
    Then I expect a file "test__.*-manifest.json.gpg" in folder "/tmp/storage/mybucket"
    And I expect a file "test__.*-persistent-data.tar.gz.gpg" in folder "/tmp/storage/mybucket"
