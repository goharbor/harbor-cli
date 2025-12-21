package main

import (
	"context"
	"errors"
	"log"

	"dagger/harbor-cli/internal/dagger"
)

const (
	harborAdminUser     = "admin"
	harborAdminPassword = "Harbor12345"

	harborImageTag = "satellite"

	postgresImage = "registry.goharbor.io/dockerhub/goharbor/harbor-db:dev"
	redisImage    = "registry.goharbor.io/dockerhub/goharbor/redis-photon:dev"
	coreImage     = "registry.goharbor.io/harbor-next/harbor-core:" + harborImageTag

	configDirPath = "./test/e2e/testconfig/config/"

	postgresPort  = 5432
	redisPort     = 6379
	corePort      = 8080
	coreDebugPort = 4001
)

func (m *HarborCli) HarborTest(ctx context.Context) (string, error) {
	core := m.setupHarborRegistry(ctx)

	// Create instance for the HarborCLI to run tests in
	test := dag.Container().
		From("golang:"+GO_VERSION+"-alpine").
		WithServiceBinding("core", core).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src").
		WithEnvVariable("TEST_HARBOR_URL", "core:8080").
		WithEnvVariable("TEST_HARBOR_USERNAME", harborAdminUser).
		WithEnvVariable("TEST_HARBOR_PASSWORD", harborAdminPassword).
		WithExec([]string{"go", "test", "-v", "./..."})

	return test.Stdout(ctx)
}

// Returns container running harbor registry with all services running
func (m *HarborCli) setupHarborRegistry(ctx context.Context) *dagger.Service {
	log.Println("setting up harbor registry environment...")

	if err := m.startPostgresql(ctx); err != nil {
		requireNoExecError(err, "start postgresql")
	}
	log.Println("postgresql service started")

	if err := m.startRedis(ctx); err != nil {
		requireNoExecError(err, "start redis")
	}
	log.Println("redis service started")

	core, err := m.startCore(ctx)
	if err != nil {
		requireNoExecError(err, "start core service")
	}
	log.Println("core service started")

	return core
}

func (m *HarborCli) startPostgresql(ctx context.Context) error {
	_, err := dag.Container().
		From(postgresImage).
		WithExposedPort(postgresPort).
		AsService().
		WithHostname("postgresql").
		Start(ctx)

	return err
}

func (m *HarborCli) startRedis(ctx context.Context) error {
	_, err := dag.Container().
		From(redisImage).
		WithExposedPort(redisPort).
		AsService().
		WithHostname("redis").
		Start(ctx)

	return err
}

func (m *HarborCli) startCore(ctx context.Context) (*dagger.Service, error) {

	coreConfig := m.Source.File(configDirPath + "core/app.conf")
	envCoreFile := m.Source.File(configDirPath + "core/env")
	runScript := m.Source.File(configDirPath + "run_env.sh")
	privatekey := m.Source.File(configDirPath + "core/private_key.pem")

	return dag.Container().
		From(coreImage).
		WithMountedFile("/etc/core/app.conf", coreConfig).
		WithMountedFile("/etc/core/private_key.pem", privatekey).
		WithMountedFile("/envFile", envCoreFile).
		WithMountedFile("/run_script", runScript).
		WithExec([]string{"chmod", "+x", "/run_script"}).
		WithExposedPort(corePort, dagger.ContainerWithExposedPortOpts{ExperimentalSkipHealthcheck: true}).
		WithExposedPort(coreDebugPort, dagger.ContainerWithExposedPortOpts{ExperimentalSkipHealthcheck: true}).
		WithEntrypoint([]string{"/run_script", "/core"}).
		AsService().
		WithHostname("core").
		Start(ctx)
}

func requireNoExecError(err error, step string) {
	var e *dagger.ExecError
	if errors.As(err, &e) {
		log.Fatalf("failed to %s (exec error): %s", step, err)
	} else {
		log.Fatalf("failed to %s (unexpected error): %s", step, err)
	}
}
