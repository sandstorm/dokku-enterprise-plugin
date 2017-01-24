Feature: I can export the Configuration Manifest to the cloud

  Scenario: Application with a database and storage
    Given an empty folder "/tmp/storage" exists
    And an empty folder "/tmp/storage/mybucket" exists
    And the cloud configuration is:
      | localStoragePath | /tmp/storage/                                    |
      | type             | local                                            |
      | storageBucket    | /tmp/storage/mybucket                            |
      | encryptionKey    | foofoofoofooofoofoofoofoofooofoofoofoofoofooofoo |

    And I have an empty Dockerfile application
    And I call dokku "apps:destroy test --force"
    And an empty folder "/temp/test" exists
    And a file "/temp/test/storage.txt" is created with contents:
    """
    Hallo Welt - we must survive!
    """
    And I call dokku "mariadb:destroy test42 --force"

    And I call dokku "apps:create test"
    And I call dokku "mariadb:create test42"
    And I call dokku "mariadb:link test42 test"
    And I call dokku "config:set test FOO=bar"
    And I call dokku "storage:mount test /tmp/test/:/b"
    And I deploy the application as "test"
    And I execute the following SQL statements on database "test42":
    """
    CREATE TABLE foo (id int, content text);
    INSERT INTO foo VALUES (1, "haha");
    """
    When I call dokku "cloud:backup test"
    Then I expect a file "test__.*-manifest.json.gpg" in folder "/tmp/storage/mybucket"
    And I expect a file "test__.*-persistent-data.tar.gz.gpg" in folder "/tmp/storage/mybucket"
    And I expect a file "test__.*-code.tar.gz.gpg" in folder "/tmp/storage/mybucket"


    Given I call dokku "mariadb:destroy copy42 --force"
    And an empty folder "/tmp/copy" exists
    When I call dokku "cloud:createAppFromCloud copy <lastExportedCloudId>"
    Then I expect a file "/tmp/copy/storage.txt" with contents:
    """
    Hallo Welt - we must survive!
    """
    And the SQL statement "SELECT content FROM foo WHERE id=1" on database "copy42" must return "haha"