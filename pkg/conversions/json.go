package conversions

import (
	"encoding/json"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func ToJSON(span opentracing.Span, obj interface{}) string {
	jsonObj, err := json.Marshal(obj)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(
			log.Error(err),
			log.Event("converting object to json"),
			log.Object("object", obj),
		)
		return ""
	}
	return string(jsonObj)
}
