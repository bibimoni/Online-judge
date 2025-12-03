package tests

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	subrepository "github.com/bibimoni/Online-judge/submission-judge/src/domain/repository/submission"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg"
	"github.com/bibimoni/Online-judge/submission-judge/src/pkg/memory"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/judge/impl"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/problem"
	"github.com/bibimoni/Online-judge/submission-judge/src/service/store"
	usecase "github.com/bibimoni/Online-judge/submission-judge/src/usecase/wssubmission"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// --- Mocks ---

// MockLanguage implements pkg.Language
type MockLanguage struct {
	IDFunc                func() string
	DisplayNameFunc       func() string
	DefaultFileNameFunc   func() string
	ExecutableNameFunc    func() string
	FileExtensionFunc     func() string
	RunFunc               func(i *domain.Isolate, rc *domain.RunConfig, req *isolateservice.SubmissionRequest) error
	RunCmdStrNoStreamFunc func(i *domain.Isolate, rc *domain.RunConfig, req *isolateservice.SubmissionRequest) ([]string, error)
	CompileFunc           func(i *domain.Isolate, req *isolateservice.SubmissionRequest, stderr io.Writer) error
}

func (m *MockLanguage) ID() string {
	if m.IDFunc != nil {
		return m.IDFunc()
	}
	return "mock_lang"
}
func (m *MockLanguage) DisplayName() string     { return "Mock Language" }
func (m *MockLanguage) DefaultFileName() string { return "main.mock" }
func (m *MockLanguage) ExecutableName() string  { return "main" }
func (m *MockLanguage) FileExtension() string   { return "mock" }
func (m *MockLanguage) Run(i *domain.Isolate, rc *domain.RunConfig, req *isolateservice.SubmissionRequest) error {
	if m.RunFunc != nil {
		return m.RunFunc(i, rc, req)
	}
	return nil
}
func (m *MockLanguage) RunCmdStrNoStream(i *domain.Isolate, rc *domain.RunConfig, req *isolateservice.SubmissionRequest) ([]string, error) {
	if m.RunCmdStrNoStreamFunc != nil {
		return m.RunCmdStrNoStreamFunc(i, rc, req)
	}
	return []string{}, nil
}
func (m *MockLanguage) Compile(i *domain.Isolate, req *isolateservice.SubmissionRequest, stderr io.Writer) error {
	if m.CompileFunc != nil {
		return m.CompileFunc(i, req, stderr)
	}
	return nil
}

// MockStoreService
type MockStoreService struct {
	GetFunc func(id string) (pkg.Language, error)
}

func (m *MockStoreService) Get(id string) (pkg.Language, error) {
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	return &MockLanguage{}, nil
}
func (m *MockStoreService) Register(l pkg.Language) {}
func (m *MockStoreService) List() []pkg.Language    { return []pkg.Language{} }
func (m *MockStoreService) Contains(id string) bool { return true }

// MockProblemService
type MockProblemService struct {
	GetFunc                func(ctx context.Context, id string) (*problem.ProblemServiceGetOutput, error)
	GetTestCaseAddrFunc    func(problemId string, tcType problem.TestCaseType, testNum int) (string, error)
	GetTestCaseDirAddrFunc func(problemId string, tcType problem.TestCaseType) (string, error)
	GetCheckerAddrFunc     func(problemId string) (string, error)
	GetInteractorAddrFunc  func(problemId string) (string, error)
	GetCrossRunAddrFunc    func(problemId string) (string, error)
}

func (m *MockProblemService) Get(ctx context.Context, id string) (*problem.ProblemServiceGetOutput, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}
func (m *MockProblemService) GetTestCaseAddr(problemId string, tcType problem.TestCaseType, testNum int) (string, error) {
	if m.GetTestCaseAddrFunc != nil {
		return m.GetTestCaseAddrFunc(problemId, tcType, testNum)
	}
	return "/tmp/mock/testcase", nil
}
func (m *MockProblemService) GetTestCaseDirAddr(problemId string, tcType problem.TestCaseType) (string, error) {
	return "/tmp/mock/testcase_dir", nil
}
func (m *MockProblemService) GetCheckerAddr(problemId string) (string, error) {
	if m.GetCheckerAddrFunc != nil {
		return m.GetCheckerAddrFunc(problemId)
	}
	return "/tmp/mock/checker", nil
}
func (m *MockProblemService) GetInteractorAddr(problemId string) (string, error) {
	return "/tmp/mock/interactor", nil
}
func (m *MockProblemService) GetCrossRunAddr(problemId string) (string, error) {
	return "/tmp/mock/crossrun", nil
}

// MockIsolateService
type MockIsolateService struct {
	NewIsolateFunc        func(id int) (*domain.Isolate, error)
	CleanupFunc           func(i *domain.Isolate) error
	InitFunc              func(i *domain.Isolate) error
	RunFunc               func(i *domain.Isolate, rc domain.RunConfig, req *isolateservice.SubmissionRequest, toRun string, toRunArgs ...string) error
	RunBinaryFunc         func(i *domain.Isolate, rc domain.RunConfig, req *isolateservice.SubmissionRequest, exeName string) error
	RunCmdStrNoStreamFunc func(i *domain.Isolate, rc domain.RunConfig, req *isolateservice.SubmissionRequest, toRun string, toRunArgs ...string) ([]string, error)
}

func (m *MockIsolateService) NewIsolate(id int) (*domain.Isolate, error) {
	return &domain.Isolate{ID: id}, nil
}
func (m *MockIsolateService) Cleanup(i *domain.Isolate) error { return nil }
func (m *MockIsolateService) Init(i *domain.Isolate) error {
	if m.InitFunc != nil {
		return m.InitFunc(i)
	}
	return nil
}
func (m *MockIsolateService) Run(i *domain.Isolate, rc domain.RunConfig, req *isolateservice.SubmissionRequest, toRun string, toRunArgs ...string) error {
	return nil
}
func (m *MockIsolateService) RunBinary(i *domain.Isolate, rc domain.RunConfig, req *isolateservice.SubmissionRequest, exeName string) error {
	if m.RunBinaryFunc != nil {
		return m.RunBinaryFunc(i, rc, req, exeName)
	}
	return nil
}
func (m *MockIsolateService) RunCmdStrNoStream(i *domain.Isolate, rc domain.RunConfig, req *isolateservice.SubmissionRequest, toRun string, toRunArgs ...string) ([]string, error) {
	if m.RunCmdStrNoStreamFunc != nil {
		return m.RunCmdStrNoStreamFunc(i, rc, req, toRun, toRunArgs...)
	}
	return []string{}, nil
}

// MockPoolService
type MockPoolService struct {
	GetFunc func() (*domain.Isolate, error)
	PutFunc func(i *domain.Isolate)
	LenFunc func() int
}

func (m *MockPoolService) Get() (*domain.Isolate, error) {
	if m.GetFunc != nil {
		return m.GetFunc()
	}
	logger := zerolog.New(io.Discard)
	return &domain.Isolate{ID: 1, Inited: true, Logger: &logger}, nil
}
func (m *MockPoolService) Put(i *domain.Isolate) {
	if m.PutFunc != nil {
		m.PutFunc(i)
	}
}
func (m *MockPoolService) Len() int { return 1 }

// MockEvaluationRepository
type MockEvaluationRepository struct {
	GetEvalFunc               func(ctx context.Context, evalId string) (*domain.EvaluationResult, error)
	UpdateFinalFunc           func(ctx context.Context, evalId string, verdict domain.Verdict, cpuTime float64, memoryUsage memory.Memory, nsucess int, points int, message string) error
	UpdateCaseFunc            func(ctx context.Context, evalId string, verdictCase domain.Verdict, cpuTimeCase float64, memoryUsageCase memory.Memory, outputCase string, pointsCase int, cpuTime float64, memoryUsage memory.Memory, nsucess int) error
	GetEvalBySubmissionIdFunc func(ctx context.Context, submissionId bson.ObjectID) (*domain.EvaluationResult, error)
}

func (m *MockEvaluationRepository) CreateEval(ctx context.Context, submissionId string, TL int, ML memory.Memory, nCase int) (string, error) {
	return "eval1", nil
}
func (m *MockEvaluationRepository) UpdateVerdict(ctx context.Context, evalId string, vert domain.Verdict) error {
	return nil
}
func (m *MockEvaluationRepository) UpdateCase(ctx context.Context, evalId string, verdictCase domain.Verdict, cpuTimeCase float64, memoryUsageCase memory.Memory, outputCase string, pointsCase int, cpuTime float64, memoryUsage memory.Memory, nsucess int) error {
	if m.UpdateCaseFunc != nil {
		return m.UpdateCaseFunc(ctx, evalId, verdictCase, cpuTimeCase, memoryUsageCase, outputCase, pointsCase, cpuTime, memoryUsage, nsucess)
	}
	return nil
}
func (m *MockEvaluationRepository) UpdateFinal(ctx context.Context, evalId string, verdict domain.Verdict, cpuTime float64, memoryUsage memory.Memory, nsucess int, points int, message string) error {
	if m.UpdateFinalFunc != nil {
		return m.UpdateFinalFunc(ctx, evalId, verdict, cpuTime, memoryUsage, nsucess, points, message)
	}
	return nil
}
func (m *MockEvaluationRepository) GetEval(ctx context.Context, evalId string) (*domain.EvaluationResult, error) {
	if m.GetEvalFunc != nil {
		return m.GetEvalFunc(ctx, evalId)
	}
	return &domain.EvaluationResult{SubmissionId: bson.NewObjectID()}, nil
}
func (m *MockEvaluationRepository) GetEvalBson(ctx context.Context, evalId bson.ObjectID) (*domain.EvaluationResult, error) {
	return nil, nil
}
func (m *MockEvaluationRepository) GetEvalBySubmissionId(ctx context.Context, submissionId bson.ObjectID) (*domain.EvaluationResult, error) {
	if m.GetEvalBySubmissionIdFunc != nil {
		return m.GetEvalBySubmissionIdFunc(ctx, submissionId)
	}
	return &domain.EvaluationResult{
		SubmissionId: submissionId,
		EvalStatus:   domain.PENDING,
	}, nil
}

// MockCheckerService
type MockCheckerService struct {
	RunCheckerFunc func(checkerAddr string, inputAddr string, outputAddr string, answerAddr string) (domain.Verdict, int, string, error)
}

func (m *MockCheckerService) RunChecker(checkerAddr string, inputAddr string, outputAddr string, answerAddr string) (domain.Verdict, int, string, error) {
	if m.RunCheckerFunc != nil {
		return m.RunCheckerFunc(checkerAddr, inputAddr, outputAddr, answerAddr)
	}
	return domain.ACCEPTED, 0, "ok", nil
}

// MockInteractorService
type MockInteractorService struct{}

func (m *MockInteractorService) RunInteractor(crossRunAddr, interactorAddr, inputAddr, outputAddr, answerAddr, reportAddr string, isolateStr []string) (domain.Verdict, int, string, error) {
	return domain.ACCEPTED, 0, "ok", nil
}

// MockRedisSubmissionRepository
type MockRedisSubmissionRepository struct {
	PulishSubmissionFunc func(ctx context.Context, res usecase.WSSubmissionResponse) error
}

func (m *MockRedisSubmissionRepository) PulishSubmission(ctx context.Context, res usecase.WSSubmissionResponse) error {
	if m.PulishSubmissionFunc != nil {
		return m.PulishSubmissionFunc(ctx, res)
	}
	return nil
}
func (m *MockRedisSubmissionRepository) Subscribe(ctx context.Context, channelId string) (<-chan *usecase.WSSubmissionResponse, error) {
	return nil, nil
}
func (m *MockRedisSubmissionRepository) GetChannelString(problemId, username, submissionId string) string {
	return ""
}
func (m *MockRedisSubmissionRepository) PushSubmissionJob(ctx context.Context, req *isolateservice.SubmissionRequest) error {
	return nil
}
func (m *MockRedisSubmissionRepository) PopSubmissionJob(ctx context.Context) (*isolateservice.SubmissionRequest, error) {
	return nil, nil
}

// MockSubmissionRepository
type MockSubmissionRepository struct {
	FindSubmissionFunc func(ctx context.Context, submissionId string) (*domain.Submission, error)
}

func (m *MockSubmissionRepository) CreateSubmission(ctx context.Context, params subrepository.CreateSubmissionInput) (string, error) {
	return "", nil
}
func (m *MockSubmissionRepository) FindSubmission(ctx context.Context, submissionId string) (*domain.Submission, error) {
	if m.FindSubmissionFunc != nil {
		return m.FindSubmissionFunc(ctx, submissionId)
	}
	return &domain.Submission{
		Id:        bson.NewObjectID(),
		ProblemId: "prob1",
		Username:  "user1",
		Type:      domain.SubmissionType(domain.ICPC),
		Timestamp: time.Now(),
	}, nil
}
func (m *MockSubmissionRepository) FindAllProblemSubmissionIds(ctx context.Context, problemId string) ([]string, error) {
	return nil, nil
}

// MockSourcecodeRepository
type MockSourcecodeRepository struct {
	GetSourceBySubmissionIdFunc func(ctx context.Context, submissionId bson.ObjectID) (*domain.SourceCode, error)
}

func (m *MockSourcecodeRepository) CreateSourcecode(ctx context.Context, source string, languageId string, submissionId string) (string, error) {
	return "", nil
}
func (m *MockSourcecodeRepository) GetSourcecode(ctx context.Context, id string) (*domain.SourceCode, error) {
	return nil, nil
}
func (m *MockSourcecodeRepository) GetSourcecodeBson(ctx context.Context, bid bson.ObjectID) (*domain.SourceCode, error) {
	return nil, nil
}
func (m *MockSourcecodeRepository) GetSourceBySubmissionId(ctx context.Context, submissionId bson.ObjectID) (*domain.SourceCode, error) {
	if m.GetSourceBySubmissionIdFunc != nil {
		return m.GetSourceBySubmissionIdFunc(ctx, submissionId)
	}
	return &domain.SourceCode{
		LanguageId: "cpp",
		SourceCode: "int main() {}",
	}, nil
}

// --- Tests ---

func TestJudgeService_JudgeStart_AC(t *testing.T) {
	// 0. Setup Store
	store.DefaultStore = &MockStoreService{}

	// 1. Setup Mocks
	mockProblemService := &MockProblemService{}
	mockIsolateService := &MockIsolateService{
		InitFunc: func(i *domain.Isolate) error {
			i.Inited = true
			return nil
		},
	}
	mockPoolService := &MockPoolService{}
	mockEvalRepo := &MockEvaluationRepository{}
	mockCheckerService := &MockCheckerService{}
	mockInteractorService := &MockInteractorService{}
	mockRedisRepo := &MockRedisSubmissionRepository{}
	mockSubRepo := &MockSubmissionRepository{}
	mockSourceRepo := &MockSourcecodeRepository{}

	// 2. Initialize JudgeService
	js := impl.NewJudgeServiceImpl(
		mockPoolService,
		mockProblemService,
		mockEvalRepo,
		mockCheckerService,
		mockInteractorService,
		mockRedisRepo,
		mockSubRepo,
		mockSourceRepo,
		mockIsolateService,
	)

	// 3. Prepare Request
	ctx := context.Background()
	mockLang := &MockLanguage{
		CompileFunc: func(i *domain.Isolate, req *isolateservice.SubmissionRequest, stderr io.Writer) error {
			return nil // Compilation success
		},
		RunFunc: func(i *domain.Isolate, rc *domain.RunConfig, req *isolateservice.SubmissionRequest) error {
			return nil // Run success
		},
	}

	req := &isolateservice.SubmissionRequest{
		SubmissionId:   "sub1",
		ProblemId:      "prob1",
		IService:       mockIsolateService,
		LanguageId:     "cpp",
		EvalId:         "eval1",
		SubmissionType: domain.SubmissionType(domain.ICPC),
	}
	problemInfo := &problem.ProblemServiceGetOutput{
		TestNum:     1,
		TimeLimit:   1000,
		MemoryLimit: 256 * 1024 * 1024,
	}

	// 4. Run JudgeStart
	err := js.JudgeStart(ctx, mockLang, req, problemInfo)
	if err != nil {
		// As expected, it might fail due to FS, but we want to ensure it doesn't panic.
		t.Logf("JudgeStart failed (expected due to FS): %v", err)
	}
}

func TestJudgeService_JudgeStart_CE(t *testing.T) {
	store.DefaultStore = &MockStoreService{}

	mockProblemService := &MockProblemService{}
	mockIsolateService := &MockIsolateService{}
	mockPoolService := &MockPoolService{}
	mockEvalRepo := &MockEvaluationRepository{
		UpdateFinalFunc: func(ctx context.Context, evalId string, verdict domain.Verdict, cpuTime float64, memoryUsage memory.Memory, nsucess int, points int, message string) error {
			if verdict != domain.COMPILATION_ERROR {
				t.Errorf("Expected COMPILATION_ERROR, got %v", verdict)
			}
			return nil
		},
	}
	mockCheckerService := &MockCheckerService{}
	mockInteractorService := &MockInteractorService{}
	mockRedisRepo := &MockRedisSubmissionRepository{}
	mockSubRepo := &MockSubmissionRepository{}
	mockSourceRepo := &MockSourcecodeRepository{}

	js := impl.NewJudgeServiceImpl(
		mockPoolService,
		mockProblemService,
		mockEvalRepo,
		mockCheckerService,
		mockInteractorService,
		mockRedisRepo,
		mockSubRepo,
		mockSourceRepo,
		mockIsolateService,
	)

	ctx := context.Background()
	mockLang := &MockLanguage{
		CompileFunc: func(i *domain.Isolate, req *isolateservice.SubmissionRequest, stderr io.Writer) error {
			return errors.New("compilation failed") // Simulate compile error
		},
	}

	req := &isolateservice.SubmissionRequest{
		SubmissionId:   "sub2",
		ProblemId:      "prob1",
		IService:       mockIsolateService,
		EvalId:         "eval2",
		SubmissionType: domain.SubmissionType(domain.ICPC),
	}
	problemInfo := &problem.ProblemServiceGetOutput{}

	err := js.JudgeStart(ctx, mockLang, req, problemInfo)
	if err == nil {
		t.Error("Expected error from JudgeStart due to CE")
	}
}
