Feature: I can import the Configuration Manifest

  Scenario: Error when application already exists
    Given I call dokku "apps:destroy test --force"
    And I call dokku "apps:create test"

    When I call dokku "manifest:import test" with payload:
    """
    {
      "manifest": {
      }
    }
    """
    Then I get back a message "ERROR: Application 'test' already exists!"

  Scenario: Error when manifest could not be parsed
    Given I call dokku "apps:destroy test --force"
    When I call dokku "manifest:import test" with payload:
    """
    {
      "manifest: {
      }
    }
    """
    Then I get back a message "ERROR: JSON could not be parsed"

  Scenario: When manifest contains errors, an error is shown.
    Given I call dokku "apps:destroy test --force"
    When I call dokku "manifest:import test" with payload:
    """
    {
      "manifest": {
      },
      "errors": [
        "Some error"
      ]
    }
    """
    Then I get back a message "ERROR: The manifest had errors; which means that the manifest is NOT fully self-contained"
    And I get back a message "Some error"


  Scenario: Application is created from manifest
    Given I call dokku "apps:destroy testXXX --force"
    When I call dokku "manifest:import testXXX" with payload:
    """
    {
      "manifest": {
      }
    }
    """
    And I call dokku "apps"
    Then I get back a message "testXXX"

  Scenario: Database is created inside a manifest
    Given I call dokku "apps:destroy test --force"
    Given I call dokku "mariadb:destroy testX --force"
    Given I call dokku "mariadb:destroy testfoo --force"

    When I call dokku "manifest:import test" with payload:
    """
    {
      "manifest": {
        "mariadb": [
          "[appName]X",
          "[appName]foo"
        ]
      }
    }
    """
    And I call dokku "mariadb:list"
    Then I get back a message "testX"
    Then I get back a message "testfoo"
    When I call dokku "config test"
    Then I get back a message "dokku-mariadb-testX"
    Then I get back a message "dokku-mariadb-testfoo"

  Scenario: Persistent volume is created inside a manifest
    Given I call dokku "apps:destroy test --force"

    When I call dokku "manifest:import test" with payload:
    """
    {
      "manifest": {
        "dockerOptions": {
          "deploy": [
            "-v /tmp/[appName]/HUHU:/app/x"
          ],
          "run": [
            "-v /tmp/[appName]/BAR:/app/x"
          ],
          "build": [
            "-v /tmp/[appName]/BUILD:/app/x"
          ]
        }
      }
    }
    """
    And I call dokku "docker-options test"
    Then I get back a message "/tmp/test/HUHU"
    Then I get back a message "/tmp/test/BAR"
    Then I get back a message "/tmp/test/BUILD"


  Scenario: configuration option is created inside a manifest
    Given I call dokku "apps:destroy test --force"

    When I call dokku "manifest:import test" with payload:
    """
    {
      "manifest": {
        "config": {
          "FOO": "BAR[appName] 'Bla'"
        }
      }
    }
    """
    And I call dokku "config test"
    Then I get back a message "BARtest 'Bla'"
