// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package declarative

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"go.yaml.in/yaml/v4"
)

// Format identifies a supported declarative document encoding.
type Format string

const (
	// FormatYAML encodes a document as YAML.
	FormatYAML Format = "yaml"
	// FormatJSON encodes a document as JSON.
	FormatJSON Format = "json"
)

// ParseFormat validates a user-supplied output format.
func ParseFormat(value string) (Format, error) {
	switch strings.ToLower(value) {
	case "yaml", "yml":
		return FormatYAML, nil
	case "json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unsupported format %q; expected yaml or json", value)
	}
}

// FormatForFile infers a document encoding from a filename.
func FormatForFile(filename string) (Format, error) {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".yaml", ".yml":
		return FormatYAML, nil
	case ".json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unsupported file extension %q; expected .yaml, .yml, or .json", filepath.Ext(filename))
	}
}

// Decode reads and strictly validates one configuration document.
func Decode(r io.Reader, format Format) (*Configuration, error) {
	configuration, err := decode(r, format)
	if err != nil {
		return nil, err
	}
	if err := configuration.Validate(); err != nil {
		return nil, err
	}
	configuration.Normalize()
	return configuration, nil
}

func decode(r io.Reader, format Format) (*Configuration, error) {
	var configuration Configuration
	switch format {
	case FormatYAML:
		decoder := yaml.NewDecoder(r)
		decoder.KnownFields(true)
		if err := decoder.Decode(&configuration); err != nil {
			return nil, fmt.Errorf("decode YAML configuration: %w", err)
		}
		if err := rejectAdditionalDocument(decoder); err != nil {
			return nil, fmt.Errorf("decode YAML configuration: %w", err)
		}
	case FormatJSON:
		decoder := json.NewDecoder(r)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&configuration); err != nil {
			return nil, fmt.Errorf("decode JSON configuration: %w", err)
		}
		if err := rejectAdditionalDocument(decoder); err != nil {
			return nil, fmt.Errorf("decode JSON configuration: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format %q", format)
	}
	return &configuration, nil
}

func rejectAdditionalDocument(decoder interface{ Decode(any) error }) error {
	var extra any
	err := decoder.Decode(&extra)
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("decode trailing content: %w", err)
	}
	return errors.New("multiple documents in one file are not supported")
}

// ReadFile loads a configuration, inferring its encoding from the extension.
func ReadFile(filename string) (*Configuration, error) {
	format, err := FormatForFile(filename)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read configuration: %w", err)
	}
	return Decode(bytes.NewReader(data), format)
}

// ReadPath loads one configuration file or overlays all configuration files in
// a directory recursively. Directory files are applied in relative-path order.
func ReadPath(path string) (*Configuration, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("inspect configuration path: %w", err)
	}
	if !info.IsDir() {
		return ReadFile(path)
	}

	filenames, err := configurationFiles(path)
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, fmt.Errorf("configuration directory %q contains no YAML or JSON files", path)
	}
	configurations := make([]*Configuration, 0, len(filenames))
	for _, filename := range filenames {
		configuration, readErr := readFragment(filename)
		if readErr != nil {
			return nil, readErr
		}
		if validateErr := configuration.validate(false); validateErr != nil {
			return nil, fmt.Errorf("validate configuration %q: %w", filename, validateErr)
		}
		configurations = append(configurations, configuration)
	}
	configuration, err := Merge(configurations...)
	if err != nil {
		return nil, fmt.Errorf("merge configuration directory %q: %w", path, err)
	}
	return configuration, nil
}

func readFragment(filename string) (*Configuration, error) {
	format, err := FormatForFile(filename)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read configuration %q: %w", filename, err)
	}
	configuration, err := decode(bytes.NewReader(data), format)
	if err != nil {
		return nil, fmt.Errorf("read configuration %q: %w", filename, err)
	}
	return configuration, nil
}

func configurationFiles(root string) ([]string, error) {
	var filenames []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path != root && entry.IsDir() && strings.HasPrefix(entry.Name(), ".") {
			return filepath.SkipDir
		}
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			return nil
		}
		switch strings.ToLower(filepath.Ext(entry.Name())) {
		case ".yaml", ".yml", ".json":
			filenames = append(filenames, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk configuration directory: %w", err)
	}
	slices.Sort(filenames)
	return filenames, nil
}

// Encode writes a normalized configuration in the requested encoding.
func Encode(w io.Writer, configuration *Configuration, format Format) error {
	if err := configuration.Validate(); err != nil {
		return err
	}
	configuration.Normalize()
	switch format {
	case FormatYAML:
		encoder := yaml.NewEncoder(w)
		encoder.SetIndent(2)
		if err := encoder.Encode(configuration); err != nil {
			return fmt.Errorf("encode YAML configuration: %w", err)
		}
		return encoder.Close()
	case FormatJSON:
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(configuration); err != nil {
			return fmt.Errorf("encode JSON configuration: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

// WriteFile atomically writes a configuration with owner-only permissions.
func WriteFile(filename string, configuration *Configuration, format Format) error {
	var data bytes.Buffer
	if err := Encode(&data, configuration, format); err != nil {
		return err
	}

	directory := filepath.Dir(filename)
	temporary, err := os.CreateTemp(directory, ".harbor-configuration-*")
	if err != nil {
		return fmt.Errorf("create temporary configuration: %w", err)
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)

	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return fmt.Errorf("secure temporary configuration: %w", err)
	}
	if _, err := temporary.Write(data.Bytes()); err != nil {
		temporary.Close()
		return fmt.Errorf("write temporary configuration: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary configuration: %w", err)
	}
	if err := os.Rename(temporaryName, filename); err != nil {
		return fmt.Errorf("replace configuration file: %w", err)
	}
	return nil
}
