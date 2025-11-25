package impl

import (
	"context"
	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	repository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/submission"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem"
	usecasews "github.com/bibimoni/Online-judge/submission-judge/src/usecase/wssubmission"
	"go.mongodb.org/mongo-driver/v2/bson"
	"testing"
)

type MockRedisRepo struct {
	Queue []*isolateservice.SubmissionRequest
}

func (m *MockRedisRepo) PulishSubmission(ctx context.Context, res usecasews.WSSubmissionResponse) error {
	return nil
}
func (m *MockRedisRepo) Subscribe(ctx context.Context, channelId string) (<-chan *usecasews.WSSubmissionResponse, error) {
	return nil, nil
}
func (m *MockRedisRepo) GetChannelString(problemId, username, submissionId string) string {
	return ""
}
func (m *MockRedisRepo) PushSubmissionJob(ctx context.Context, req *isolateservice.SubmissionRequest) error {
	m.Queue = append(m.Queue, req)
	return nil
}
func (m *MockRedisRepo) PopSubmissionJob(ctx context.Context) (*isolateservice.SubmissionRequest, error) {
	if len(m.Queue) == 0 {
		return nil, context.DeadlineExceeded // Simulate empty
	}
	req := m.Queue[0]
	m.Queue = m.Queue[1:]
	return req, nil
}

type MockPoolService struct{}

func (m *MockPoolService) Get() (*domain.Isolate, error) {
	return &domain.Isolate{ID: 1, Logger: nil}, nil // Return dummy isolate
}
func (m *MockPoolService) Put(i *domain.Isolate) {}
func (m *MockPoolService) Len() int              { return 1 }

type MockProblemService struct{}

func (m *MockProblemService) Get(ctx context.Context, id string) (*problem.ProblemServiceGetOutput, error) {
	return &problem.ProblemServiceGetOutput{TimeLimit: 1000, MemoryLimit: 1024}, nil
}
func (m *MockProblemService) GetTestCaseAddr(problemId string, tcType problem.TestCaseType, testNum int) (string, error) {
	return "/tmp/testcase", nil
}
func (m *MockProblemService) GetTestCaseDirAddr(problemId string, tcType problem.TestCaseType) (string, error) {
	return "/tmp/testcase_dir", nil
}
func (m *MockProblemService) GetCheckerAddr(problemId string) (string, error) {
	return "/tmp/checker", nil
}
func (m *MockProblemService) GetInteractorAddr(problemId string) (string, error) {
	return "/tmp/interactor", nil
}
func (m *MockProblemService) GetCrossRunAddr(problemId string) (string, error) {
	return "/tmp/crossrun", nil
}

type MockEvalRepo struct{}

func (m *MockEvalRepo) CreateEval(ctx context.Context, submissionId string, timeLimit int, memoryLimit memory.Memory, testNum int) (string, error) {
	return "eval_123", nil
}
func (m *MockEvalRepo) GetEval(ctx context.Context, evalId string) (*domain.EvaluationResult, error) {
	return &domain.EvaluationResult{}, nil
}
func (m *MockEvalRepo) GetEvalBson(ctx context.Context, evalId bson.ObjectID) (*domain.EvaluationResult, error) {
	return &domain.EvaluationResult{}, nil
}
func (m *MockEvalRepo) GetEvalBySubmissionId(ctx context.Context, submissionId bson.ObjectID) (*domain.EvaluationResult, error) {
	return &domain.EvaluationResult{}, nil
}
func (m *MockEvalRepo) UpdateVerdict(ctx context.Context, evalId string, vert domain.Verdict) error {
	return nil
}
func (m *MockEvalRepo) UpdateCase(ctx context.Context, evalId string, verdictCase domain.Verdict, cpuTimeCase float64, memoryUsageCase memory.Memory, outputCase string, pointsCase int, cpuTime float64, memoryUsage memory.Memory, nsucess int) error {
	return nil
}
func (m *MockEvalRepo) UpdateFinal(ctx context.Context, evalId string, verdict domain.Verdict, cpuTime float64, memoryUsage memory.Memory, nsucess int, points int, message string) error {
	return nil
}

type MockCheckerService struct{}

func (m *MockCheckerService) RunChecker(checkerAddr string, inputAddr string, outputAddr string, answerAddr string) (domain.Verdict, int, string, error) {
	return domain.ACCEPTED, 100, "ok", nil
}

type MockInteractorService struct{}

func (m *MockInteractorService) RunInteractor(crossRunAddr, interactorAddr, inputAddr, outputAddr, answerAddr, reportAddr string, isolateStr []string) (domain.Verdict, int, string, error) {
	return domain.ACCEPTED, 100, "ok", nil
}

type MockSourceCodeRepo struct{}

func (m *MockSourceCodeRepo) CreateSourcecode(ctx context.Context, sourcecode string, languageId string, submissionId string) (string, error) {
	return "sc_123", nil
}
func (m *MockSourceCodeRepo) GetSourcecode(ctx context.Context, id string) (*domain.SourceCode, error) {
	return &domain.SourceCode{}, nil
}
func (m *MockSourceCodeRepo) GetSourcecodeBson(ctx context.Context, bid bson.ObjectID) (*domain.SourceCode, error) {
	return &domain.SourceCode{}, nil
}
func (m *MockSourceCodeRepo) GetSourceBySubmissionId(ctx context.Context, submissionId bson.ObjectID) (*domain.SourceCode, error) {
	return &domain.SourceCode{}, nil
}

type MockSubmissionRepo struct{}

func (m *MockSubmissionRepo) CreateSubmission(ctx context.Context, params repository.CreateSubmissionInput) (string, error) {
	return "sub_123", nil
}
func (m *MockSubmissionRepo) FindSubmission(ctx context.Context, submissionId string) (*domain.Submission, error) {
	return &domain.Submission{}, nil
}
func (m *MockSubmissionRepo) FindAllProblemSubmissionIds(ctx context.Context, problemId string) ([]string, error) {
	return []string{}, nil
}

type MockIsolateService struct{}

func (m *MockIsolateService) NewIsolate(id int) (*domain.Isolate, error) {
	return &domain.Isolate{ID: id}, nil
}
func (m *MockIsolateService) Cleanup(i *domain.Isolate) error { return nil }
func (m *MockIsolateService) Init(i *domain.Isolate) error    { return nil }
func (m *MockIsolateService) Run(i *domain.Isolate, rc domain.RunConfig, req *isolateservice.SubmissionRequest, toRun string, toRunArgs ...string) error {
	return nil
}
func (m *MockIsolateService) RunBinary(i *domain.Isolate, rc domain.RunConfig, req *isolateservice.SubmissionRequest, exeName string) error {
	return nil
}
func (m *MockIsolateService) RunCmdStrNoStream(i *domain.Isolate, rc domain.RunConfig, req *isolateservice.SubmissionRequest, toRun string, toRunArgs ...string) ([]string, error) {
	return []string{}, nil
}

func TestJudgeFlow(t *testing.T) {
	mockRedis := &MockRedisRepo{Queue: make([]*isolateservice.SubmissionRequest, 0)}
	mockPool := &MockPoolService{}
	mockProblem := &MockProblemService{}
	mockEval := &MockEvalRepo{}
	mockChecker := &MockCheckerService{}
	mockInteractor := &MockInteractorService{}
	mockSub := &MockSubmissionRepo{}
	mockSource := &MockSourceCodeRepo{}
	mockIsolate := &MockIsolateService{}
	js := NewJudgeServiceImpl(
		mockPool,
		mockProblem,
		mockEval,
		mockChecker,
		mockInteractor,
		mockRedis,
		mockSub,
		mockSource,
		mockIsolate,
	)
	req := &isolateservice.SubmissionRequest{
		SubmissionId: "sub_123",
		ProblemId:    "prob_1",
		LanguageId:   "cpp",
		IService:     mockIsolate, // Inject mock
	}
	problemInfo := &problem.ProblemServiceGetOutput{
		TimeLimit: 1000,
	}
	ctx := context.Background()
	err := js.Judge(ctx, req, problemInfo)
	if err != nil {
		t.Fatalf("Judge() returned error: %v", err)
	}
	if len(mockRedis.Queue) != 1 {
		t.Errorf("Expected 1 job in Redis queue, got %d", len(mockRedis.Queue))
	} else {
		t.Log("Success: Job pushed to Redis")
	}
	t.Log("Verification Complete: JudgeService correctly integrates with Redis Queue.")
}
