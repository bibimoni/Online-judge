# Problem service

## API Schema
```
- POST /problem/add?problemId=...
- GET /problem/latest-version?problemId=...
- GET /problem/get/${problemId}/${version}/{static files}
- GET /problem/all
Return:
{
	"list": [
		352118,
		440176
	]
}
```

### Problem directory structure
```
+-- problem.json
+-- statement.pdf
+-- checker
+-- tests/
|   +-- input/
|   |   +-- 1
|   |   +-- 2
|   |   +-- ...
|   +-- output/
|   |   +-- 1
|   |   +-- 2
|   |   +-- ...
+-- interactor (interactive problems)
+-- CrossRun.jar (interactive problems)
```
### Sample of a problem.json
```json
{
  "ID": "...",
  "problem-id": 327480,
  "name": "Chơi đùa cùng bộ bài",
  "short-name": "playful-card-game",
  "tags": [
    "dp"
  ],
  "test-num": 103,
  "time-limit": 1000,
  "memory-limit": 268435456
}
```

## Current problems
- Only accept English statement (could use Vietnamese inside)