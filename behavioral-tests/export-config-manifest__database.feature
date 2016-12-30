Feature: I can export the Configuration Manifest -- part DATABASE

  In order to create reproducible installations across instances of services, I want to be able to export
  all configuration to a "manifest".

  Scenario: Application with a single connected database
    Given I have an empty Dockerfile application
    And I call dokku "apps:destroy test --force"
    And I call dokku "apps:create test"
    And I call dokku "mariadb:create test42"
    And I call dokku "mariadb:link test42 test"
#    And I deploy the application as "test"

    When I call dokku "manifest:export test"
    Then I get back a JSON object with the following structure:
      | version              | equals   | 1           |
      | appName              | equals   | test        |
      | manifest.mariadb.[0] | equals   | [appName]42 |
      | errors               | is empty |             |

  Scenario: Application with two connected databases
    Given I have an empty Dockerfile application
    And I call dokku "apps:destroy test --force"
    And I call dokku "apps:create test"
    And I call dokku "mariadb:create test42"
    And I call dokku "mariadb:link test42 test"
    And I call dokku "mariadb:create testabc"
    And I call dokku "mariadb:link testabc test"
#    And I deploy the application as "test"

    When I call dokku "manifest:export test"
    Then I get back a JSON object with the following structure:
      | version              | equals   | 1            |
      | appName              | equals   | test         |
      | manifest.mariadb.[0] | equals   | [appName]42  |
      | manifest.mariadb.[1] | equals   | [appName]abc |
      | errors               | is empty |              |

  Scenario: If a database does not contain the application name in its name, an error is added to the list.
    Given I have an empty Dockerfile application
    And I call dokku "apps:destroy test --force"
    And I call dokku "apps:create test"
    And I call dokku "mariadb:create aa42"
    And I call dokku "mariadb:link aa42 test"
    And I call dokku "mariadb:create abc"
    And I call dokku "mariadb:link abc test"
#    And I deploy the application as "test"

    When I call dokku "manifest:export test"
    Then I get back a JSON object with the following structure:
      | version              | equals   | 1                                                                              |
      | appName              | equals   | test                                                                           |
      | manifest.mariadb.[0] | equals   | aa42                                                                           |
      | manifest.mariadb.[1] | equals   | abc                                                                            |
      | errors.[0]           | contains | mariadb.DATABASE_URL: did not find application name 'test' inside string: aa42 |
      | errors.[1]           | contains | did not find application name 'test' inside string: abc                        |