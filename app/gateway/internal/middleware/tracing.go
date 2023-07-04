package middleware

import (
	"bufio"
	"errors"
	"github.com/cherish-chat/imcloudx-server/app/gateway/internal/svc"
	"github.com/cherish-chat/imcloudx-server/common"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"net"
	"net/http"
)

const serviceName = "gateway"

func Tracing(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	tracer := otel.Tracer(common.TraceName)
	propagator := otel.GetTextMapPropagator()
	return func(c *gin.Context) {
		spanName := c.Request.URL.Path
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		spanCtx, span := tracer.Start(
			ctx,
			spanName,
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(
				serviceName, spanName, c.Request)...),
		)
		defer span.End()
		propagator.Inject(spanCtx, propagation.HeaderCarrier(c.Writer.Header()))
		trw := &WithCodeResponseWriter{ResponseWriter: c.Writer, Code: 200}
		c.Writer = trw
		c.Next()
		span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(trw.Code)...)
		span.SetStatus(semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(
			trw.Code, oteltrace.SpanKindServer))
	}
}

// A WithCodeResponseWriter is a helper to delay sealing a http.ResponseWriter on writing code.
type WithCodeResponseWriter struct {
	gin.ResponseWriter
	Code int
}

// Flush flushes the response writer.
func (w *WithCodeResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Header returns the http header.
func (w *WithCodeResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Hijack implements the http.Hijacker interface.
// This expands the Response to fulfill http.Hijacker if the underlying http.ResponseWriter supports it.
func (w *WithCodeResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacked, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacked.Hijack()
	}

	return nil, nil, errors.New("server doesn't support hijacking")
}

// Write writes bytes into w.
func (w *WithCodeResponseWriter) Write(bytes []byte) (int, error) {
	return w.ResponseWriter.Write(bytes)
}

// WriteHeader writes code into w, and not sealing the writer.
func (w *WithCodeResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.Code = code
}
