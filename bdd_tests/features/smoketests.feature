Feature: Basic set of smoke tests


  Scenario: Check if mock service is available
    Given the system is in default state
     When I look for executable file insights-results-aggregator-mock
     Then I should find that file on PATH


  Scenario: Check if cleaner displays help message
    Given the system is in default state
     When I start the mock service with the --help command line flag
     Then I should see help messages displayed on standard output


  Scenario: Check if cleaner displays version info
    Given the system is in default state
     When I start the mock service with the --version command line flag
     Then I should see version info displayed on standard output


  Scenario: Check if cleaner displays authors
    Given the system is in default state
     When I start the mock service with the --authors command line flag
     Then I should see info about authors displayed on standard output
