# file: sender.feature
Feature: send notify about event
    Scenario: simple creation request
        Given I start consuming
        When I send create request with title "title test notify 1" description "descr test notify 1" startAt "2020-01-02T15:00:00Z" endAt "2020-01-02T16:00:00Z" notifyAt "2020-01-02T14:00:00Z"
        Then I should receive message with title "title test notify 1"
        And I stop consuming
