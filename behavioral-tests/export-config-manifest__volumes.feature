Feature: I can export the Configuration Manifest -- part VOLUMES/Docker Options

  Scenario: Application with a volume
    Given I have an empty Dockerfile application
    And I call dokku "apps:destroy test --force"
    And I call dokku "apps:create test"
    And I call dokku "storage:mount test /tmp/test/:/b"

    When I call dokku "manifest:export test"
    Then I get back a JSON object with the following structure:
      | version                           | equals   | 1                     |
      | appName                           | equals   | test                  |
      | manifest.dockerOptions.deploy.[0] | equals   | -v /tmp/[appName]/:/b |
      | manifest.dockerOptions.run.[0]    | equals   | -v /tmp/[appName]/:/b |
      | errors                            | is empty |                       |

  Scenario: Application with a custom docker option
    Given I have an empty Dockerfile application
    And I call dokku "apps:destroy test --force"
    And I call dokku "apps:create test"
    And I call dokku "docker-options:add test deploy --restart=on-failure:15"
    And I call dokku "docker-options:add test build --foobar"
    And I call dokku "docker-options:add test run --foobar2"

    When I call dokku "manifest:export test"
    Then I get back a JSON object with the following structure:
      | version                           | equals   | 1                       |
      | appName                           | equals   | test                    |
      | manifest.dockerOptions.deploy.[0] | equals   | --restart=on-failure:15 |
      | manifest.dockerOptions.build.[0]  | equals   | --foobar                |
      | manifest.dockerOptions.run.[0]    | equals   | --foobar2               |
      | errors                            | is empty |                         |

  Scenario: If a storage does not contain the application name in its path, an error is added to the list.
    Given I have an empty Dockerfile application
    And I call dokku "apps:destroy test --force"
    And I call dokku "apps:create test"
    And I call dokku "storage:mount test /tmp/foobar/:/b"

    When I call dokku "manifest:export test"
    Then I get back a JSON object with the following structure:
      | version                           | equals | 1                                                                                             |
      | appName                           | equals | test                                                                                          |
      | manifest.dockerOptions.deploy.[0] | equals | -v /tmp/foobar/:/b                                                                            |
      | manifest.dockerOptions.run.[0]    | equals | -v /tmp/foobar/:/b                                                                            |
      | errors.[0]                        | equals | dockerOptions.deploy: did not find application name 'test' inside string: -v /tmp/foobar/:/b. |
      | errors.[1]                        | equals | dockerOptions.run: did not find application name 'test' inside string: -v /tmp/foobar/:/b.    |