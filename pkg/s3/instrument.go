package s3

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/arfan21/project-sprint-shopifyx-api/pkg/s3")
