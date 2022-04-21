package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider() {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("test-service"),
		),
	)
	if err != nil {
		log.Printf("failed to create resource: %s", err)
		return
	}

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns
	conn, err := grpc.DialContext(
		ctx,
		"datadog-agent:30080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Printf("failed to create gRPC connection: %s", err)
		return
	}

	metricClient := otlpmetricgrpc.NewClient(otlpmetricgrpc.WithGRPCConn(conn))
	metricExp, err := otlpmetric.New(ctx, metricClient)
	if err != nil {
		log.Printf("failed to create gRPC connection: %s", err)
		return
	}
	pusher := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			metricExp,
		),
		controller.WithExporter(metricExp),
	)
	global.SetMeterProvider(pusher)
	if err := pusher.Start(ctx); err != nil {
		log.Printf("failed to start the pusher: %s", err)
		return
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Printf("failed to create trace exporter: %s", err)
		return
	}
	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func main() {
	ctx := context.Background()
	meter := global.Meter("test-meter")
	workDuration, _ := meter.SyncInt64().UpDownCounter(
		"otlp_client/work_duration",
		instrument.WithDescription("The duration for a work in milliseconds"),
	)

	log.Printf("Initialize a provider ...")
	initProvider()

	tracer := otel.Tracer("test-tracer")

	// labels represent additional key-value descriptors that can be bound to a
	// metric observer or recorder.
	commonLabels := []attribute.KeyValue{
		attribute.String("labelA", "chocolate"),
		attribute.String("labelB", "raspberry"),
		attribute.String("labelC", "vanilla"),
	}

	// work begins
	ctx, span := tracer.Start(
		ctx,
		"CollectorExporter-Example",
		trace.WithAttributes(commonLabels...),
	)
	defer span.End()
	limit := 100
	for i := 0; i < limit; i++ {
		now := time.Now()
		_, iSpan := tracer.Start(ctx, fmt.Sprintf("Sample-%d", i))
		log.Printf("Doing really hard work (%d / %d)\n", i+1, limit)
		workDuration.Add(ctx, time.Since(now).Milliseconds(), commonLabels...)
		<-time.After(time.Second)
		iSpan.End()
	}

	log.Printf("Done!")
}
