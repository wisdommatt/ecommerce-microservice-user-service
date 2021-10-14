package password

import (
	"context"
	"errors"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(ctx context.Context, passwordStr string, cost int) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "hash-password")
	defer span.Finish()
	ext.Component.Set(span, "bcrypt")
	span.SetTag("param.cost", cost)
	span.SetTag("time", time.Now())

	if passwordStr == "" {
		ext.Error.Set(span, true)
		span.SetTag("param.passwordStr", passwordStr)
		span.LogFields(
			log.String("event", "password hash error"),
			log.String("message", "password is empty"),
		)
		return "", errors.New("Password must not be empty !")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(passwordStr), cost)
	if err != nil {
		ext.Error.Set(span, true)
		span.SetTag("param.passwordStr", passwordStr)
		span.LogFields(
			log.String("event", "password hash error"),
			log.Error(err),
		)
		return "", err
	}
	return string(passwordHash), nil
}
