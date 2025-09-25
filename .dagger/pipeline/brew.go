package pipeline

import (
	"context"
	"fmt"
	"strings"

	"dagger/harbor-cli/internal/dagger"
)

func (s *Pipeline) BrewFormula(ctx context.Context, dist *dagger.Directory) (*dagger.Directory, error) {
	archs := []string{"amd64", "arm64"}
	shaMap := map[string]string{}

	for _, arch := range archs {
		tarPath := fmt.Sprintf("/archive/harbor-cli_%s_darwin_%s.tar.gz", s.appVersion, arch)

		shasum, err := s.dag.Container().
			From("alpine").
			WithMountedFile("./bin.tar.gz", dist.File(tarPath)).
			WithWorkdir("/").
			WithExec([]string{"sha256sum", "./bin.tar.gz"}).
			Stdout(ctx)
		if err != nil {
			return nil, err
		}

		shaMap[arch] = strings.Split(shasum, " ")[0] // taking only the shasum
	}

	formulaContent := FormulaTemplate(s.appVersion, shaMap)

	// Creating Formula.rb
	dir := s.dag.Directory().WithNewFile("Formula.rb", formulaContent)

	dist = dist.WithDirectory("brew", dir)

	return dist, nil
}

func FormulaTemplate(appVer string, shaMap map[string]string) string {
	content := fmt.Sprintf(`class harbor-cli < Formula
  desc "Harbor CLI - A command-line interface for CNCF Harbor, the cloud native registry!"
  homepage "https://github.com/goharbor/harbor-cli"
  license "Apache 2.0" 
  version "%s" 

  if Hardware::CPU.intel?
    url "https://github.com/goharbor/harbor-cli/releases/download/%s/harbor-cli_%s_darwin_amd64.tar.gz"
    sha256 "%s"
  end

  if Hardware::CPU.arm?
    url "https://github.com/goharbor/harbor-cli/releases/download/%s/harbor-cli_%s_darwin_arm64.tar.gz"
    sha256 "%s"
  end

  def install
    bin.install "harbor-cli"
  end

  test do
    system "#{bin}/harbor-cli", "--version"
  end
end
`, appVer,
		appVer, appVer, shaMap["amd64"],
		appVer, appVer, shaMap["arm64"],
	)

	return content
}
