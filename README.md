# Online-judge
## Run distributed judge
```bash
docker-compose run --rm \
  -p 8001:8001 \
  -e SUBMISSION_IS_MAIN_JUDGE=true \
  -e SUBMISSION_JUDGE_PROBLEM_DIR=/tmp/cache \ # you can change to whatever you like
  -e SUBMISSION_PORT=8001 \ # submission port
  -e SUBMISSION_REDIS_URI=redissubmissionjudge:6379 \ # the main judge's Redis address
  -e SUBMISSION_REDIS_PASSWORD=root \ # main judge's Redis password
  -v /tmp/judge_cache:/tmp/cache \ # mapped to host machine local address to prevent loss
  submission-judge
```
