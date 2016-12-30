Feature: I can export the Configuration Manifest -- part CONFIG vars

  Scenario: Application with a config option
    Given I have an empty Dockerfile application
    And I call dokku "apps:destroy test --force"
    And I call dokku "apps:create test"
    And I call dokku "config:set test FOO=bar"

    When I call dokku "manifest:export test"
    Then I get back a JSON object with the following structure:
      | version             | equals   | 1    |
      | appName             | equals   | test |
      | manifest.config.FOO | equals   | bar  |
      | errors              | is empty |      |

  Scenario: default config is not serialized
    Given I have an empty Dockerfile application
    And I call dokku "apps:destroy test --force"
    And I call dokku "apps:create test"
    And I call dokku "mariadb:create test42"
    And I call dokku "mariadb:link test42 test"
    And I call dokku "config:set test FOO=bar"
    And I deploy the application as "test"

    When I call dokku "manifest:export test"
    Then I get back a JSON object with the following structure:
      | version             | equals   | 1    |
      | appName             | equals   | test |
      | manifest.config.FOO | equals   | bar  |
      | manifest.config     | count    | 1    |
      | errors              | is empty |      |