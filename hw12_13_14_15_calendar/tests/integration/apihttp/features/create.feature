# file: create.feature
Feature: create event
    Scenario: simple creation request
        When I send create request with title "title test create 1" description "descr test create 1" startAt "2020-01-02T15:00:00Z" endAt "2020-01-02T16:00:00Z" notifyAt "2020-01-02T14:00:00Z"
        Then the response should be without errors
        And the response has correct event id
    Scenario: creation request with constraint violation
        Given there are events:
            | title               | description         | during                                    | notify_at | user_id |
            | title test create 2 | descr test create 2 | 2020-02-02T15:00:00Z 2020-02-02T16:00:00Z |           | 1       |
        When I send create request with title "title test 2-1" description "descr test 2-1" startAt "2020-02-02T14:30:00Z" endAt "2020-02-02T15:30:00Z" notifyAt ""
        Then the response should has error text "time is busy"
