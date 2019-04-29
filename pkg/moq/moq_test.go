package moq

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMoq(t *testing.T) {
	m, err := New("testpackages/example", "")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "PersonStore")
	if err != nil {
		t.Errorf("m.Mock: %s", err)
	}
	s := buf.String()
	// assertions of things that should be mentioned
	var strs = []string{
		"package example",
		"type PersonStoreMock struct",
		"CreateFunc func(ctx context.Context, person *Person, confirm bool) error",
		"GetFunc func(ctx context.Context, id string) (*Person, error)",
		"func (mock *PersonStoreMock) Create(ctx context.Context, person *Person, confirm bool) error",
		"func (mock *PersonStoreMock) Get(ctx context.Context, id string) (*Person, error)",
		"panic(\"PersonStoreMock.CreateFunc: method is nil but PersonStore.Create was just called\")",
		"panic(\"PersonStoreMock.GetFunc: method is nil but PersonStore.Get was just called\")",
		"lockPersonStoreMockGet.Lock()",
		"mock.calls.Get = append(mock.calls.Get, callInfo)",
		"lockPersonStoreMockGet.Unlock()",
		"// ID is the id argument value",
	}
	for _, str := range strs {
		if !strings.Contains(s, str) {
			t.Errorf("expected but missing: \"%s\"", str)
		}
	}
}

func TestMoqWithStaticCheck(t *testing.T) {
	m, err := New("testpackages/example", "")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "PersonStore")
	if err != nil {
		t.Errorf("m.Mock: %s", err)
	}
	s := buf.String()
	// assertions of things that should be mentioned
	var strs = []string{
		"package example",
		"var _ PersonStore = &PersonStoreMock{}",
		"type PersonStoreMock struct",
		"CreateFunc func(ctx context.Context, person *Person, confirm bool) error",
		"GetFunc func(ctx context.Context, id string) (*Person, error)",
		"func (mock *PersonStoreMock) Create(ctx context.Context, person *Person, confirm bool) error",
		"func (mock *PersonStoreMock) Get(ctx context.Context, id string) (*Person, error)",
		"panic(\"PersonStoreMock.CreateFunc: method is nil but PersonStore.Create was just called\")",
		"panic(\"PersonStoreMock.GetFunc: method is nil but PersonStore.Get was just called\")",
		"lockPersonStoreMockGet.Lock()",
		"mock.calls.Get = append(mock.calls.Get, callInfo)",
		"lockPersonStoreMockGet.Unlock()",
		"// ID is the id argument value",
	}
	for _, str := range strs {
		if !strings.Contains(s, str) {
			t.Errorf("expected but missing: \"%s\"", str)
		}
	}
}

func TestMoqExplicitPackage(t *testing.T) {
	m, err := New("testpackages/example", "different")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "PersonStore")
	if err != nil {
		t.Errorf("m.Mock: %s", err)
	}
	s := buf.String()
	// assertions of things that should be mentioned
	var strs = []string{
		"package different",
		"type PersonStoreMock struct",
		"CreateFunc func(ctx context.Context, person *example.Person, confirm bool) error",
		"GetFunc func(ctx context.Context, id string) (*example.Person, error)",
		"func (mock *PersonStoreMock) Create(ctx context.Context, person *example.Person, confirm bool) error",
		"func (mock *PersonStoreMock) Get(ctx context.Context, id string) (*example.Person, error)",
	}
	for _, str := range strs {
		if !strings.Contains(s, str) {
			t.Errorf("expected but missing: \"%s\"", str)
		}
	}
}

func TestMoqExplicitPackageWithStaticCheck(t *testing.T) {
	m, err := New("testpackages/example", "different")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "PersonStore")
	if err != nil {
		t.Errorf("m.Mock: %s", err)
	}
	s := buf.String()
	// assertions of things that should be mentioned
	var strs = []string{
		"package different",
		"var _ example.PersonStore = &PersonStoreMock{}",
		"type PersonStoreMock struct",
		"CreateFunc func(ctx context.Context, person *example.Person, confirm bool) error",
		"GetFunc func(ctx context.Context, id string) (*example.Person, error)",
		"func (mock *PersonStoreMock) Create(ctx context.Context, person *example.Person, confirm bool) error",
		"func (mock *PersonStoreMock) Get(ctx context.Context, id string) (*example.Person, error)",
	}
	for _, str := range strs {
		if !strings.Contains(s, str) {
			t.Errorf("expected but missing: \"%s\"", str)
		}
	}
}

// TestVeradicArguments tests to ensure variadic work as
// expected.
// see https://github.com/matryer/moq/issues/5
func TestVariadicArguments(t *testing.T) {
	m, err := New("testpackages/variadic", "")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "Greeter")
	if err != nil {
		t.Errorf("m.Mock: %s", err)
	}
	s := buf.String()
	// assertions of things that should be mentioned
	var strs = []string{
		"package variadic",
		"type GreeterMock struct",
		"GreetFunc func(ctx context.Context, names ...string) string",
		"return mock.GreetFunc(ctx, names...)",
	}
	for _, str := range strs {
		if !strings.Contains(s, str) {
			t.Errorf("expected but missing: \"%s\"", str)
		}
	}
}

func TestNothingToReturn(t *testing.T) {
	m, err := New("testpackages/example", "")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "PersonStore")
	if err != nil {
		t.Errorf("m.Mock: %s", err)
	}
	s := buf.String()
	if strings.Contains(s, `return mock.ClearCacheFunc(id)`) {
		t.Errorf("should not have return for items that have no return arguments")
	}
	// assertions of things that should be mentioned
	var strs = []string{
		"mock.ClearCacheFunc(id)",
	}
	for _, str := range strs {
		if !strings.Contains(s, str) {
			t.Errorf("expected but missing: \"%s\"", str)
		}
	}
}

func TestChannelNames(t *testing.T) {
	m, err := New("testpackages/channels", "")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "Queuer")
	if err != nil {
		t.Errorf("m.Mock: %s", err)
	}
	s := buf.String()
	var strs = []string{
		"func (mock *QueuerMock) Sub(topic string) (<-chan Queue, error)",
	}
	for _, str := range strs {
		if !strings.Contains(s, str) {
			t.Errorf("expected but missing: \"%s\"", str)
		}
	}
}

func TestImports(t *testing.T) {
	m, err := New("testpackages/imports/two", "")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "DoSomething")
	if err != nil {
		t.Errorf("m.Mock: %s", err)
	}
	s := buf.String()
	var strs = []string{
		`	"sync"`,
		`	"github.com/matryer/moq/pkg/moq/testpackages/imports/one"`,
	}
	for _, str := range strs {
		if !strings.Contains(s, str) {
			t.Errorf("expected but missing: \"%s\"", str)
		}
		if len(strings.Split(s, str)) > 2 {
			t.Errorf("more than one: \"%s\"", str)
		}
	}
}

func TestTemplateFuncs(t *testing.T) {
	fn := templateFuncs["Exported"].(func(string) string)
	if fn("var") != "Var" {
		t.Errorf("exported didn't work: %s", fn("var"))
	}
}

func TestVendoredPackages(t *testing.T) {
	m, err := New("testpackages/vendoring/user", "")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "Service")
	if err != nil {
		t.Errorf("mock error: %s", err)
	}
	s := buf.String()
	// assertions of things that should be mentioned
	var strs = []string{
		`"github.com/matryer/somerepo"`,
	}
	for _, str := range strs {
		if !strings.Contains(s, str) {
			t.Errorf("expected but missing: \"%s\"", str)
		}
	}
}

func TestVendoredBuildConstraints(t *testing.T) {
	m, err := New("testpackages/buildconstraints/user", "")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "Service")
	if err != nil {
		t.Errorf("mock error: %s", err)
	}
	s := buf.String()
	// assertions of things that should be mentioned
	var strs = []string{
		`"github.com/matryer/buildconstraints"`,
	}
	for _, str := range strs {
		if !strings.Contains(s, str) {
			t.Errorf("expected but missing: \"%s\"", str)
		}
	}
}

// TestDotImports tests for https://github.com/matryer/moq/issues/21.
func TestDotImports(t *testing.T) {
	preDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Getwd: %s", err)
	}
	err = os.Chdir("testpackages/dotimport")
	if err != nil {
		t.Errorf("Chdir: %s", err)
	}
	defer func() {
		err := os.Chdir(preDir)
		if err != nil {
			t.Errorf("Chdir back: %s", err)
		}
	}()
	m, err := New(".", "moqtest_test")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "Service")
	if err != nil {
		t.Errorf("mock error: %s", err)
	}
	s := buf.String()
	if strings.Contains(s, `"."`) {
		t.Error("contains invalid dot import")
	}
}

func TestEmptyInterface(t *testing.T) {
	m, err := New("testpackages/emptyinterface", "")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "Empty")
	if err != nil {
		t.Errorf("mock error: %s", err)
	}
	s := buf.String()
	if strings.Contains(s, `"sync"`) {
		t.Error("contains sync import, although this package isn't used")
	}
}

func helperExec(t *testing.T, args []string, dir string) []byte {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Exec: %s; Output:\n%s\n", err, output)
	}
	return output
}

func TestGoGenerateVendoredPackages(t *testing.T) {
	output := helperExec(t, []string{"go", "generate", "./..."}, "testpackages/gogenvendoring")
	if bytes.Contains(output, []byte(`vendor/`)) {
		t.Error("contains vendor directory in import path")
	}
}

func TestImportedPackageWithSameName(t *testing.T) {
	m, err := New("testpackages/samenameimport", "")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "Example")
	if err != nil {
		t.Errorf("mock error: %s", err)
	}
	s := buf.String()
	if !strings.Contains(s, `a samename.A`) {
		t.Error("missing samename.A to address the struct A from the external package samename")
	}
}

func TestParamShadowsImport(t *testing.T) {
	m, err := New("testpackages/shadow", "gen")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}

	genDir := "./testpackages/shadow/gen"
	err = os.MkdirAll(genDir, 0777)
	if err != nil {
		t.Fatalf("mkdir: %s", err)
	}
	defer os.Remove(genDir)

	f, err := os.Create(filepath.Join(genDir, "shadowmock.go"))
	if err != nil {
		t.Fatalf("failed to create shadowmock: %s", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	err = m.Mock(f, "Shadower")
	if err != nil {
		t.Fatalf("mock error: %s", err)
	}

	helperExec(t, []string{"go", "vet", genDir}, "")
}

// TestSliceReturns tests to ensure returning a slice works as expected
func TestSliceReturns(t *testing.T) {
	m, err := New("testpackages/slices", "")
	if err != nil {
		t.Fatalf("moq.New: %s", err)
	}
	var buf bytes.Buffer
	err = m.Mock(&buf, "Search")
	if err != nil {
		t.Fatalf("m.Mock: %s", err)
	}
	s := buf.String()
	// assertions of things that should be mentioned
	var strs = []string{
		"package slices",
		"type SearchMock struct",
		"GetNamesFunc func(ctx context.Context, search string) []string",
		"GetOtherFunc func(ctx context.Context, search string) ([]string, error)",
		"GetMultipleSlicesFunc func(ctx context.Context, search ...string) ([]string, []error)",
		"GetPointerSlicesFunc func(ctx context.Context, search string) ([]*string, error)",
		"MultiSearchFunc func(ctx context.Context, searchTerms ...string) []string",
		"func (mock *SearchMock) GetNames(ctx context.Context, search string) []string",
		"func (mock *SearchMock) GetOther(ctx context.Context, search string) ([]string, error)",
		"func (mock *SearchMock) GetMultipleSlices(ctx context.Context, search ...string) ([]string, []error)",
		"func (mock *SearchMock) GetPointerSlices(ctx context.Context, search string) ([]*string, error)",
		"func (mock *SearchMock) MultiSearch(ctx context.Context, searchTerms ...string) []string",
	}
	for _, str := range strs {
		if !strings.Contains(s, str) {
			t.Errorf("expected but missing: \"%s\"", str)
		}
	}
}
