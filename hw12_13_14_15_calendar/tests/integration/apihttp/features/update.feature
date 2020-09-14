# file: update.feature
Feature: update event
    Scenario: simple update request
        Given there are events:
            | id | title               | description         | during                                    | notify_at | user_id |
            | 1  | title test update 1 | descr test update 1 | 2020-02-02T15:00:00Z 2020-02-02T16:00:00Z |           | 1       |
        When I send update request with eventID "1" title "title test updated" description "descr test updated" startAt "2020-01-02T15:00:00Z" endAt "2020-01-02T16:00:00Z" notifyAt "2020-01-02T14:00:00Z"
        Then the response should be without errors
        And the response should has title "title test updated"
    Scenario: update request with constraint violation
        Given there are events:
            | id | title               | description         | during                                    | notify_at | user_id |
            | 1  | title test update 1 | descr test update 1 | 2020-02-02T15:00:00Z 2020-02-02T16:00:00Z |           | 1       |
            | 2  | title test update 2 | descr test update 2 | 2020-02-02T17:00:00Z 2020-02-02T18:00:00Z |           | 1       |
        When I send update request with eventID "1" title "title test 2-1" description "descr test 2-1" startAt "2020-02-02T16:30:00Z" endAt "2020-02-02T17:30:00Z" notifyAt ""
        Then the response should has error text "time is busy"
