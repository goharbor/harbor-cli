package main

import (
	"context"

	"dagger/harbor-cli/internal/dagger"
)

// Executes Go tests
func (m *HarborCli) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	err := m.init(ctx, source)
	if err != nil {
		return "", err
	}

	test := dag.Container().
		From("golang:"+m.GoVersion+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+m.GoVersion)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+m.GoVersion)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"go", "test", "-v", "./..."})
	return test.Stdout(ctx)
}

// Executes Go tests and returns TestReport in json file
// TestReport executes Go tests and returns only the JSON report file
func (m *HarborCli) TestReport(ctx context.Context, source *dagger.Directory) (*dagger.File, error) {
	err := m.init(ctx, source)
	if err != nil {
		return nil, err
	}

	reportName := "TestReport.json"
	test := dag.Container().
		From("golang:"+m.GoVersion+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+m.GoVersion)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+m.GoVersion)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec([]string{"go", "install", "gotest.tools/gotestsum@latest"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"gotestsum", "--jsonfile", reportName, "./..."})

	return test.File(reportName), nil
}

// Tests Coverage of code base
func (m *HarborCli) TestCoverage(ctx context.Context, source *dagger.Directory) (*dagger.File, error) {
	err := m.init(ctx, source)
	if err != nil {
		return nil, err
	}
	coverage := "coverage.out"
	test := dag.Container().
		From("golang:"+m.GoVersion+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+m.GoVersion)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+m.GoVersion)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec([]string{"go", "install", "gotest.tools/gotestsum@latest"}).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"gotestsum", "--", "-coverprofile=" + coverage, "./..."})

	return test.File(coverage), nil
}

// TestCoverageReport processes coverage data and returns a formatted markdown report
func (m *HarborCli) TestCoverageReport(ctx context.Context, source *dagger.Directory) (*dagger.File, error) {
	err := m.init(ctx, source)
	if err != nil {
		return nil, err
	}
	coverageFile := "coverage.out"
	reportFile := "coverage-report.md"
	test := dag.Container().
		From("golang:"+m.GoVersion+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+m.GoVersion)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+m.GoVersion)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "add", "--no-cache", "bc"}).
		WithExec([]string{"go", "test", "-coverprofile=" + coverageFile, "./..."})
	return test.WithExec([]string{"sh", "-c", `
        echo "<h2> 📊 Test Coverage Results</h2>" > ` + reportFile + `
        if [ ! -f "` + coverageFile + `" ]; then
            echo "<p>❌ Coverage file not found!</p>" >> ` + reportFile + `
            exit 1
        fi
        total_coverage=$(go tool cover -func=` + coverageFile + ` | grep total: | grep -Eo '[0-9]+\.[0-9]+')
        echo "DEBUG: Total coverage is $total_coverage" >&2
        if (( $(echo "$total_coverage >= 80.0" | bc -l) )); then
            emoji="✅"
        elif (( $(echo "$total_coverage >= 60.0" | bc -l) )); then
            emoji="⚠️"
        else
            emoji="❌"
        fi
		echo "<p><b>Total coverage: $emoji $total_coverage% (Target: 80%)</b></p>" >> ` + reportFile + `
		echo "<details><summary>Detailed package coverage</summary><pre>" >> ` + reportFile + `
        go tool cover -func=` + coverageFile + ` >> ` + reportFile + `
        echo "</pre></details>" >> ` + reportFile + `
        cat ` + reportFile + ` >&2
    `}).File(reportFile), nil
}

// TestIntegration runs integration tests against a local Harbor service
func (m *HarborCli) TestIntegration(ctx context.Context, source *dagger.Directory) (string, error) {
	err := m.init(ctx, source)
	if err != nil {
		return "", err
	}

	// 1. Redis
	redis := dag.Container().
		From("redis:7-alpine").
		WithExposedPort(6379).
		AsService()

	// 2. Postgres
	postgres := dag.Container().
		From("postgres:13-alpine").
		WithEnvVariable("POSTGRES_PASSWORD", "root123").
		WithEnvVariable("POSTGRES_DB", "registry").
		WithExposedPort(5432).
		AsService()

	// 3. Harbor Core
	harborHost := "harbor.local"
	// Minimal app.conf content
	appConf := `
appname = Harbor
runmode = prod
enablegzip = true
[prod]
httpport = 8080
`
	configDir := dag.Directory().WithNewFile("app.conf", appConf)

	harbor := dag.Container().
		From("goharbor/harbor-core:v2.10.0").
		WithMountedDirectory("/etc/core", configDir).
		WithEnvVariable("PORT", "8080").
		WithEnvVariable("CORE_SECRET", "1234567890123456").
		WithEnvVariable("HARBOR_ADMIN_PASSWORD", "Harbor12345").
		WithEnvVariable("POSTGRESQL_HOST", "postgres").
		WithEnvVariable("POSTGRESQL_PORT", "5432").
		WithEnvVariable("POSTGRESQL_USERNAME", "postgres").
		WithEnvVariable("POSTGRESQL_PASSWORD", "root123").
		WithEnvVariable("POSTGRESQL_DATABASE", "registry").
		WithEnvVariable("_REDIS_URL_CORE", "redis://redis:6379").
		WithEnvVariable("EXT_ENDPOINT", "http://"+harborHost+":8080").
		WithServiceBinding("postgres", postgres).
		WithServiceBinding("redis", redis).
		WithExposedPort(8080).
		AsService()

	test := dag.Container().
		From("golang:"+m.GoVersion+"-alpine").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+m.GoVersion)).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-"+m.GoVersion)).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithServiceBinding("redis", redis).
		WithServiceBinding("postgres", postgres).
		WithServiceBinding(harborHost, harbor).
		WithEnvVariable("HARBOR_URL", "http://"+harborHost+":8080").
		WithEnvVariable("HARBOR_USERNAME", "admin").
		WithEnvVariable("HARBOR_PASSWORD", "Harbor12345").
		WithExec([]string{"go", "test", "-v", "./..."})

	return test.Stdout(ctx)
}
