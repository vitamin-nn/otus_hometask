# file: get.feature
Feature: get events
    Scenario: day events
        Given there are events:
            | title        | description  | during                                    | notify_at | user_id |
            | title test 1 | descr test 1 | 2020-01-02T15:00:00Z 2020-01-02T16:00:00Z |           | 1       |
            | title test 2 | descr test 2 | 2020-01-02T18:00:00Z 2020-01-02T19:00:00Z |           | 1       |
            | title test 3 | descr test 3 | 2020-01-03T18:00:00Z 2020-01-03T19:00:00Z |           | 1       |
        When I send get events day request with beginAt "2020-01-02T15:00:00Z"
        Then count event in the response should be 2
    Scenario: week events
        Given there are events:
            | title        | description  | during                                    | notify_at | user_id |
            | title test 1 | descr test 1 | 2020-01-02T15:00:00Z 2020-01-02T16:00:00Z |           | 1       |
            | title test 2 | descr test 2 | 2020-01-02T18:00:00Z 2020-01-02T19:00:00Z |           | 1       |
            | title test 3 | descr test 3 | 2020-01-04T18:00:00Z 2020-01-04T19:00:00Z |           | 1       |
            | title test 4 | descr test 4 | 2020-01-15T18:00:00Z 2020-01-15T19:00:00Z |           | 1       |
        When I send get events week request with beginAt "2019-12-30T15:00:00Z"
        Then count event in the response should be 3
    Scenario: month events
        Given there are events:
            | title        | description  | during                                    | notify_at | user_id |
            | title test 1 | descr test 1 | 2020-01-02T15:00:00Z 2020-01-02T16:00:00Z |           | 1       |
            | title test 2 | descr test 2 | 2020-01-02T18:00:00Z 2020-01-02T19:00:00Z |           | 1       |
            | title test 3 | descr test 3 | 2020-01-04T18:00:00Z 2020-01-04T19:00:00Z |           | 1       |
            | title test 4 | descr test 4 | 2020-01-25T18:00:00Z 2020-01-26T19:00:00Z |           | 1       |
            | title test 5 | descr test 5 | 2020-02-15T18:00:00Z 2020-02-15T19:00:00Z |           | 1       |
        When I send get events month request with beginAt "2020-01-01T15:00:00Z"
        Then count event in the response should be 4