## API Endpoints

### Submission

#### /api/v1/submission/submit

* Required authentication: `true`
* Content-Type: application/json
* Param: `none`
* Description: Only supports `ICPC` as a `submission_type`.
<!--     * Submit user's code. if `submission_type` is `CUSTOM` then `problem_id` can be null, otherwise `problem_id` must be present. -->
```jsonld
POST /api/v1/submission/submit
{
    "problem_id": "123a" || null, // string
    "code": "int main() {}", // string, required
    "language": "c++17", // string, required
    "submission_type", // string, required
}
```

* Response: Returns the `submission_id` if submission has been queued.
```jsonld
{
    "data": {
        "message": "Submit successfully!",
        "id": "68b5715e7144975ef3d8987b"
    },
    "success": true
}
```

#### /api/v1/submission/view/:submission_id

* Required authentication: `not nessecary`
* Content-Type: application/json
* Param: `submission_id` (type `string`)
* Description:
    * View a submission.
    *  Result access resources: 

<!-- |  | Code & Sample cases | Hidden cases preview |
| -------- | -------- | ------ |
| Submission's owner `&&` ongoing contest  | Yes | No |
| Others `&&` ongoing contest | No | No |
| _ `&&` `!`ongoing contest | Yes | Yes| 
 -->
```jsonld
GET /api/v1/submission/view/:submission_id
```
<!-- {
    "username": "distiled" || null, // string
} -->
* Return value:
```jsonld
{
    {
	"data": {
		"problem_id": "445985",
		"verdict": "ACCEPTED",
		"verdict_case": [
			"ACCEPTED",
			"ACCEPTED",
			"ACCEPTED",
			"ACCEPTED",
			"ACCEPTED",
			"ACCEPTED",
			"ACCEPTED",
			"ACCEPTED",
			"ACCEPTED"
		],
		"cpu_time": 0.099,
		"cpu_time_case": [
			0.029,
			0.028,
			0.016,
			0.026,
			0.02,
			0.022,
			0.02,
			0.051,
			0.099
		],
		"memory_usage": "77176832B",
		"memory_usage_case": [
			"46714880B",
			"46727168B",
			"46673920B",
			"52293632B",
			"50909184B",
			"50888704B",
			"50634752B",
			"53420032B",
			"77176832B"
		],
		"n_success": 9,
		"outputs": [
			"ok 6 numbers",
			"ok 10 numbers",
			"ok 10 numbers",
			"ok 100 numbers",
			"ok 100 numbers",
			"ok 100 numbers",
			"ok 100 numbers",
			"ok 1000 numbers",
			"ok 1 number(s): \"2314375\""
		],
		"points": 1,
		"points_case": [
			1,
			1,
			1,
			1,
			1,
			1,
			1,
			1,
			1
		],
		"message": "Accepted",
		"n_cases": 9,
		"tl": 1000,
		"ml": "268435456B",
		"username": "john_doe",
		"timestamp": "2025-09-05T17:56:23.779Z",
		"type": "ICPC",
		"language": "PyPy3",
		"source_code": "# Author: distiled (converted from C++)\n\ndef main():\n    import sys\n    input = sys.stdin.readline  # fast input\n\n    tt = int(input())\n    for _ in range(tt):\n        n, r, c = map(int, input().split())\n        h = list(map(int, input().split()))\n        w = list(map(int, input().split()))\n\n        ans = 0\n        for i in range(n):\n            ans += ((w[i] + r - 1) // r) * ((h[i] + c - 1) // c)\n        print(ans)\n\n\nif __name__ == \"__main__\":\n    main()",
		"eval_status": "FINISHED"
	},
	"success": true
    }
}
```

<!-- #### /api/v1/submission/contest/:contest_id
* Required Authentication: `No`
* Param: `contest_id` (type `string`)
* Description: Return all the `submission_id` and related information of a contest.
```jsonld
GET /api/v1/submission/contest/:contest_id
{
    "results": [ // required (can be empty)
        "submission_id": "123456", // required
        "message": "Accepted", // required
        "cpu_time": 1420, // required(ms)
        "memory_usage": 512512, // required(KB)
        "language": "C++ 17", // required
        "index": "A", // problem index in contest
        "name": "Data miner", // required, problem's name
        "created_at": "9/7/2025...", // required
    ]
}
```
 -->
#### /api/v1/submission/problem/view/:problem_id
* Required Authentication: `No`
* Param: `problem_id` (type `string`)
* Description: Return all the `submission_id` and related information of a problem.
```jsonld
GET /api/v1/submission/problem/:problem_id
{
    {
	"data": {
		"Submissions": [
			{
				"username": "john_doe",
				"submission_id": "68b5715e7144975ef3d8987b",
				"problem_id": "445985",
				"timestamp": "2025-09-01T10:11:42.383Z",
				"language": "C++14",
				"verdict": "ACCEPTED",
				"verdict_case": [
					"ACCEPTED",
					"ACCEPTED",
					"ACCEPTED",
					"ACCEPTED",
					"ACCEPTED",
					"ACCEPTED",
					"ACCEPTED",
					"ACCEPTED",
					"ACCEPTED"
				],
				"cpu_time": 0.014,
				"cpu_time_case": [
					0.004,
					0.001,
					0.001,
					0.001,
					0.001,
					0.001,
					0.001,
					0.009,
					0.014
				],
				"memory_usage": "4431872B",
				"memory_usage_case": [
					"1286144B",
					"1286144B",
					"1286144B",
					"1286144B",
					"1286144B",
					"1286144B",
					"1286144B",
					"1286144B",
					"4431872B"
				],
				"n_success": 9,
				"outputs": [
					"ok 6 numbers",
					"ok 10 numbers",
					"ok 10 numbers",
					"ok 100 numbers",
					"ok 100 numbers",
					"ok 100 numbers",
					"ok 100 numbers",
					"ok 1000 numbers",
					"ok 1 number(s): \"2314375\""
				],
				"points": 1,
				"points_case": [
					1,
					1,
					1,
					1,
					1,
					1,
					1,
					1,
					1
				],
				"message": "Accepted",
				"eval_status": "FINISHED"
			}
		]
	},
	"success": true
    }
}
```

#### /api/v1/submission/ws?username=:username?problem_id=:problem_id?submission_id=:submission_id
* Required Authentication: `No`
* Param: `problem_id`, `username`, `submission_id` (type `string`). The params is not required but will default to all (`*`).
* Description: Establish a connection to server and receive updates about all the submissions describe in the query
* Example of an information.
```jsonld
{
	"username": "john_doe",
	"submission_id": "68a63ef599f00b4987bb9617",
	"problem_id": "445985",
	"timestamp": "2025-08-20T21:32:37.905Z",
	"language": "C++14",
	"verdict_case": [
		"ACCEPTED"
	],
	"cpu_time": 0.013,
	"cpu_time_case": [
		0.013
	],
	"memory_usage": "1317011456B",
	"memory_usage_case": [
		"1286144B"
	],
	"n_success": 1,
	"outputs": [
		"ok 6 numbers"
	],
	"points_case": [
		1
	],
	"eval_status": "JUDGING"
}```

# Adding a Horizontal Judge Node

This workflow describes how to scale your judging capacity by adding a new `submission-judge` instance on a **different computer**.

## Prerequisites

1.  **Docker & Docker Compose** installed on the new machine.
2.  **Network Access**: The new machine must be able to reach:
    *   **Redis**: The main Redis instance (e.g., `redis-submission-judge` on the main server).
    *   **Problem Service**: The Problem Service API (e.g., `problem` on the main server).
    *   **MongoDB**: (Optional, if the judge writes directly to DB, otherwise it just needs Redis). *Note: Current implementation writes to DB.*

## Step 1: Prepare Environment Variables

Create a `.env` file on the new machine. You need to point the services to the **IP Address** of your main server (where Redis and Problem Service are running).

```bash
# Main Server IP (Replace with actual IP)
MAIN_SERVER_IP=192.168.1.100

# Service Configuration
SUBMISSION_HOST=0.0.0.0
SUBMISSION_PORT=8000
SUBMISSION_LOG_LEVEL=debug

# Redis Connection (Critical for Job Queue)
SUBMISSION_REDIS_URI=${MAIN_SERVER_IP}:6379
SUBMISSION_REDIS_PASSWORD=root

# MongoDB Connection (For writing results)
SUBMISSION_MONGODB_URI=mongodb://${MAIN_SERVER_IP}:27017/submissionjudgedb
SUBMISSION_MONGODB_DATABASE_NAME=submissionjudgedb

# Problem Service (For downloading test cases)
# Note: Ensure this port is exposed on the main server
PROBLEM_HOST=${MAIN_SERVER_IP}
PROBLEM_PORT=3000

# Judge Configuration
SUBMISSION_NUMBER_OF_JUDGE=5
SUBMISSION_JUDGE_ID_OFFSET=0 
# Tip: If you run multiple nodes, you might want to adjust ID_OFFSET to avoid ID collisions if logs/metrics use it, 
# though for stateless workers it matters less.

# Paths
# The judge will use this local cache for downloaded test cases
SUBMISSION_JUDGE_PROBLEM_DIR=/app/cache/problems
```

## Step 2: Create Docker Compose

Create a `docker-compose.yml` for the judge node:

```yaml
services:
  submission-judge-node:
    image: your-docker-registry/submission-judge:latest # Or build locally
    container_name: submission-judge-node
    restart: always
    env_file:
      - .env
    ports:
      - "8000:8000"
    volumes:
      # Local cache for test cases (Hybrid Approach)
      - ./judge-cache:/app/cache/problems
    privileged: true # Required for Isolate sandbox
```

## Step 3: Run the Judge

1.  Start the container:
    ```bash
    docker-compose up -d
    ```

2.  **Verify**:
    *   Check logs: `docker logs -f submission-judge-node`
    *   You should see it connecting to Redis and waiting for jobs.
    *   When a submission comes in, if it doesn't have the test case locally, it will log that it's downloading from `http://${MAIN_SERVER_IP}:3000/...`.

## Troubleshooting

*   **Connection Refused**: Ensure ports `6379` (Redis), `27017` (Mongo), and `3000` (Problem) are exposed on the Main Server and not blocked by firewall.
*   **Download Failures**: Check if `PROBLEM_HOST` and `PROBLEM_PORT` are correct and reachable from inside the container.
