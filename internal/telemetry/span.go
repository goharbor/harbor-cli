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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Encapsulate can be applied to a span to indicate that this span should
// collapse its children by default.
func Encapsulate() trace.SpanStartOption {
	return trace.WithAttributes(attribute.Bool(UIEncapsulateAttr, true))
}

// Encapsulated can be applied to a child span to indicate that it should be
// collapsed by default.
func Encapsulated() trace.SpanStartOption {
	return trace.WithAttributes(attribute.Bool(UIEncapsulatedAttr, true))
}

// Internal can be applied to a span to indicate that this span should not be
// shown to the user by default.
func Internal() trace.SpanStartOption {
	return trace.WithAttributes(attribute.Bool(UIInternalAttr, true))
}

// Passthrough can be applied to a span to cause the UI to skip over it and
// show its children instead.
func Passthrough() trace.SpanStartOption {
	return trace.WithAttributes(attribute.Bool(UIPassthroughAttr, true))
}

// End is a helper to end a span with an error if the function returns an error.
//
// It is optimized for use as a defer one-liner with a function that has a
// named error return value, conventionally `rerr`.
//
//	defer telemetry.End(span, func() error { return rerr })
func End(span trace.Span, fn func() error) {
	if err := fn(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}
