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

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// LiveSpanProcessor is a SpanProcessor whose OnStart calls OnEnd on the
// underlying SpanProcessor in order to send live telemetry.
type LiveSpanProcessor struct {
	sdktrace.SpanProcessor
}

func NewLiveSpanProcessor(exp sdktrace.SpanExporter) *LiveSpanProcessor {
	return &LiveSpanProcessor{
		SpanProcessor: sdktrace.NewBatchSpanProcessor(
			// NOTE: span heartbeating is handled by the Cloud exporter
			exp,
			sdktrace.WithBatchTimeout(NearlyImmediate),
		),
	}
}

func (p *LiveSpanProcessor) OnStart(ctx context.Context, span sdktrace.ReadWriteSpan) {
	// Send a read-only snapshot of the live span downstream so it can be
	// filtered out by FilterLiveSpansExporter. Otherwise the span can complete
	// before being exported, resulting in two completed spans being sent, which
	// will confuse traditional OpenTelemetry services.
	p.SpanProcessor.OnEnd(SnapshotSpan(span))
}
