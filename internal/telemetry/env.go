// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telemetry

import (
	"context"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func PropagationEnv(ctx context.Context) []string {
	carrier := NewEnvCarrier(false)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return carrier.Env
}

type EnvCarrier struct {
	System bool
	Env    []string
}

func NewEnvCarrier(system bool) *EnvCarrier {
	return &EnvCarrier{
		System: system,
	}
}

var _ propagation.TextMapCarrier = (*EnvCarrier)(nil)

func (c *EnvCarrier) Get(key string) string {
	envName := strings.ToUpper(key)
	for _, env := range c.Env {
		env, val, ok := strings.Cut(env, "=")
		if ok && env == envName {
			return val
		}
	}
	if c.System {
		if envVal := os.Getenv(envName); envVal != "" {
			return envVal
		}
	}
	return ""
}

func (c *EnvCarrier) Set(key, val string) {
	c.Env = append(c.Env, strings.ToUpper(key)+"="+val)
}

func (c *EnvCarrier) Keys() []string {
	keys := make([]string, 0, len(c.Env))
	for _, env := range c.Env {
		env, _, ok := strings.Cut(env, "=")
		if ok {
			keys = append(keys, env)
		}
	}
	return keys
}
