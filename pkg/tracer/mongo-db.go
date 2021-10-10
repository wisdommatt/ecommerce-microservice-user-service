package tracer

import (
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func SetMongoDBSpanComponentTags(span opentracing.Span, collectionName string) {
	ext.DBInstance.Set(span, collectionName)
	ext.DBType.Set(span, "mongodb")
	ext.SpanKindRPCClient.Set(span)
	span.SetTag("time", time.Now())
}
