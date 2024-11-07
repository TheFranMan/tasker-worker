# Tasker Worker
This application processes task requests.

These requests are broken up into one or more steps, with each step containing one or more jobs. Each step is processed synchronously, while its jobs are processed asynchronously. Once a step's jobs have completed successfully, the next step's jobs are processed. This approach allows concurrent processing of tasks (jobs), all while retaining a strict ordering of when certain jobs are to be processed.

An overview of the lifecycle of a request:
   * A request is placed in the `requests` table. It is given the status `Queuing`.
   * Once the request is picked up, it's status is updated to `Inprogress`, and it's first step jobs are placed into the `jobs` table with the status `new`.
   * The `new` jobs are picked up and placed in the status `inprogress`. The job's callback is triggered with the values from the request's params and extras column.
   * The job can return successfully or with an error:
      * No error returned, the job is marked as `completed`. The requests extras value may be updated based on what was returned from the callback.
      * Failure error returned, the job's status is updated to `failure` and it's error recorded. This is regarded as an unrecoverable error and as such the request cannot be completed. No futher steps/jobs are processed and the request is marked as `Failed`.
      * Retry error returned, the job's status is updated to `retry` and error its recorded. While an error was returned, it is regarded as a recoverable error and the job should be retried. A new job is inserted with the same values (except created date) with a `new` status. This will get picked up on the next run, all while leaving a history of attempts in the table.
   * Once all steps's jobs are no longer in `new` or `inprogress` status, we can calculate if we can go to the next step.
      * All jobs have a `completed` status, the next step's jobs are inserted.
      * If any job has a `failure` status, mark the request as `Failed`

Several cron jobs are used to pick up the requests and jobs for processing.

## Dependencies
   * **MySQL**: Persist the requested tasks.

## Run

In the project root, first create a `.env` file with the below env vars, then run `docker compose up -d`.

| Name                        | Required           | Default      | Description                                    |
| --------------------------- | ------------------ | ------------ | ---------------------------------------------- |
| PORT                        | :white_check_mark: | n/a          | The server port                                |
| ENV                         | :x:                | n/a          | The enviroment the app is running in           |
| DB_USER                     | :white_check_mark: | n/a          | The MySQL user                                 |
| DB_PASS                     | :white_check_mark: | n/a          | The MySQL password                             |
| DB_HOST                     | :white_check_mark: | n/a          | The MySQL host                                 |
| DB_PORT                     | :white_check_mark: | n/a          | The MySQL port                                 |
| DB_NAME                     | :white_check_mark: | n/a          | The MySQL database name                        |
| WRK_ENABLED                 | :x:                | true         | Flag to enable the workers                     |
| WRK_JOB_CRON                | :x:                | */15 * * * * | Cronjob to run the job worker                  |
| WRK_REQUEST_NEW_CRON        | :x:                | */5 * * * *  | Cronjob to run the new requests worker         |
| WRK_REQUEST_INPROGRESS_CRON | :x:                | */1 * * * *  | Cronjob to run the in-progress requests worker |
| URL_SERVICE1                | :white_check_mark: | n/a          | URL to service 1                               |
| URL_SERVICE2                | :white_check_mark: | n/a          | URL to service 2                               |
| URL_SERVICE3                | :white_check_mark: | n/a          | URL to service 3                               |



To run the project separate from Docker, and so avaoiding the Docker build step, from the project root follow the below steps:
   1. `docker compose up -p`.
   2. `docker compose stop worker`.
   3. `go get && go build && ./worker`

## Testing
Tests can be run from the root dir using: `go test -v ./... --tags=integration`. The `integration` tag runs the integration tests using Docker. Remove this if you only want to run unit tests.

## Routes
   * `GET /heartbeat`

     #### Description:

     Simple endpoint used to return a 200 response to check for application health.

     #### Responses:
     * 200

   * `GET /metrics`

     #### Description:

     Used by the Prometheus scrapper to gather site statistics

     #### Responses:
     * 200

        Body:
        ```
       # HELP go_gc_duration_seconds A summary of the wall-time pause (stop-the-world) duration in garbage collection cycles.
       # TYPE go_gc_duration_seconds summary
       go_gc_duration_seconds{quantile="0"} 0.000521397
       go_gc_duration_seconds{quantile="0.25"} 0.000521397
       go_gc_duration_seconds{quantile="0.5"} 0.000521397
       go_gc_duration_seconds{quantile="0.75"} 0.000521397
       go_gc_duration_seconds{quantile="1"} 0.000521397
       ...
       ```