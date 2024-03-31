package opentelemetry

import (
	"context"
	"geek-basic-go/webook/internal/service/sms"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type OtelSmsService struct {
	svc    sms.Service
	tracer trace.Tracer
}

func NewOtelSmsService(svc sms.Service, tracer trace.Tracer) *OtelSmsService {
	return &OtelSmsService{svc: svc, tracer: tracer}
}

func (d OtelSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	ctx, span := d.tracer.Start(ctx, "sms")
	defer span.End()
	span.SetAttributes(attribute.String("tpl", tplId))
	span.AddEvent("发送短信")
	err := d.svc.Send(ctx, tplId, args, numbers...)
	if err != nil {
		span.RecordError(err)
	}
	return nil
}
