Feature: I can export the Configuration Manifest, containing the full config for an application with all DBs etc

  In order to create reproducible installations across instances of services, I want to be able to export
  all configuration to a "manifest".

  Scenario: Application with connected database
    Given I have an empty Dockerfile application
    And I call dokku "apps:destroy test --force"
    And I call dokku "apps:create test"
    And I call dokku "mariadb:create test42"
    And I call dokku "mariadb:link test42 test"
    And I deploy the application as "test"

    When I call dokku "manifest:export test"
    Then I get back a JSON object with the following structure:
      | version          | equals   | 1           |
      | appName          | equals   | test        |
      | manifest.mariadb | equals   | [appName]42 |
      | manifest.errors  | is empty |             |