Feature: additional nginx configuration should be taken from /app/nginx.conf.d

  In order to be able to deploy custom nginx config as part of an application,
  we want to be able to place it inside the git repository

  Scenario: Application with custom nginx.conf.d
    Given I have an empty node.js application
    When I create the file "nginx.conf.d/test.conf" with the following contents:
      """
        location /custom-config {
          return 200 'this is custom config';
        }
      """
    And I deploy the application as "test"
    When I call the URL "/custom-config" of the "test" application
    Then the response should contain "this is custom config"

    When I remove the file "nginx.conf.d/test.conf"
    And I deploy the application as "test"
    When I call the URL "/custom-config" of the "test" application
    Then the response should not contain "this is custom config"