# Online-judge
## Run distributed judge
```bash
docker-compose run --rm \
  -p 8001:8001 \
  -e SUBMISSION_IS_MAIN_JUDGE=false \
  -e SUBMISSION_JUDGE_PROBLEM_DIR=/tmp/cache \ # you can change to whatever you like
  -e SUBMISSION_PORT=8001 \ # submission port
  -e SUBMISSION_REDIS_URI=redissubmissionjudge:6379 \ # the main judge's Redis address
  -e SUBMISSION_REDIS_PASSWORD=root \ # main judge's Redis password
  -e SUBMISSION_MONGODB_URI=mongodb://mongosubmissionjudgedb:27017/submissionjudgedb
 \ # main judge's MongoDB address
  -e PROBLEM_ENDPOINT=http://problem:3000 \ # problem service address
  -v /tmp/judge_cache:/tmp/cache \ # mapped to host machine local address to prevent loss
  submission-judge
```
