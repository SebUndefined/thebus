package thebus_test

import (
	"context"
	"testing"

	"github.com/sebundefined/thebus"
)

func BenchmarkBus_Publish_NoSubscriber(b *testing.B) {
	bus, _ := thebus.New()
	b.Cleanup(func() {
		_ = bus.Close()
	})

	payload := make([]byte, 256)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := bus.Publish("benchmark/no-subscriber", payload)
		if err != nil {
			b.Fatalf("publish fail: %s", err)
		}
	}
}

func BenchmarkBus_Publish_CopyOnPublish_NoSubscriber(b *testing.B) {
	bus, _ := thebus.New(thebus.WithCopyOnPublish(true))
	b.Cleanup(func() {
		_ = bus.Close()
	})

	payload := make([]byte, 256)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := bus.Publish("benchmark/no-subscriber", payload)
		if err != nil {
			b.Fatalf("publish fail: %s", err)
		}
	}
}

func BenchmarkBus_Publish_OneSubscriber(b *testing.B) {
	bus, _ := thebus.New()
	b.Cleanup(func() {
		_ = bus.Close()
	})
	ctx, cancelFunc := context.WithCancel(context.Background())
	b.Cleanup(cancelFunc)

	subscriber, err := bus.Subscribe(ctx, "benchmark/one-subscriber",
		thebus.WithBufferSize(1<<16),
		thebus.WithDropIfFull(true))
	if err != nil {
		b.Fatalf("subscribe fail: %s", err)
	}
	go func() {
		for range subscriber.Read() {

		}
	}()
	payload := make([]byte, 256)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := bus.Publish("benchmark/no-subscriber", payload)
		if err != nil {
			b.Fatalf("publish fail: %s", err)
		}
	}
}

func BenchmarkBus_Publish_CopyOnPublish_OneSubscriber(b *testing.B) {
	bus, _ := thebus.New(thebus.WithCopyOnPublish(true))
	b.Cleanup(func() {
		_ = bus.Close()
	})
	ctx, cancelFunc := context.WithCancel(context.Background())
	b.Cleanup(cancelFunc)

	subscriber, err := bus.Subscribe(ctx, "benchmark/one-subscriber",
		thebus.WithBufferSize(1<<16),
		thebus.WithDropIfFull(true))
	if err != nil {
		b.Fatalf("subscribe fail: %s", err)
	}
	go func() {
		for range subscriber.Read() {

		}
	}()
	payload := make([]byte, 256)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := bus.Publish("benchmark/no-subscriber", payload)
		if err != nil {
			b.Fatalf("publish fail: %s", err)
		}
	}
}
