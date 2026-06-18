package nom

import (
	"context"
	"errors"
	"testing"
)

func TestMultiSubscriber_Fanout(t *testing.T) {
	t.Parallel()

	sink1 := &countingSubscriber{}
	sink2 := &countingSubscriber{}
	sink3 := &countingSubscriber{}

	multi := NewMultiSubscriber(sink1, sink2, sink3)

	err := multi.OnEvent(context.Background(), &multiTestEvent{eventType: EventActivityStarted})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sink1.count != 1 || sink2.count != 1 || sink3.count != 1 {
		t.Errorf("counts = %d, %d, %d; want 1, 1, 1", sink1.count, sink2.count, sink3.count)
	}
}

func TestMultiSubscriber_OneErrorsOthersStillCalled(t *testing.T) {
	t.Parallel()

	errSub := &countingSubscriber{err: errors.New("boom")}
	normalSub := &countingSubscriber{}

	multi := NewMultiSubscriber(errSub, normalSub)

	err := multi.OnEvent(context.Background(), &multiTestEvent{eventType: EventActivityStarted})
	if err == nil {
		t.Fatal("expected error from first subscriber")
	}

	if normalSub.count != 1 {
		t.Error("second subscriber should still receive event despite first erroring")
	}
}

func TestMultiSubscriber_NilSkipped(t *testing.T) {
	t.Parallel()

	sink := &countingSubscriber{}
	multi := NewMultiSubscriber(nil, sink, nil)

	err := multi.OnEvent(context.Background(), &multiTestEvent{eventType: EventActivityStarted})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sink.count != 1 {
		t.Error("non-nil subscriber should receive event")
	}
}

func TestMultiSubscriber_Add(t *testing.T) {
	t.Parallel()

	sink1 := &countingSubscriber{}
	multi := NewMultiSubscriber(sink1)

	sink2 := &countingSubscriber{}
	multi.Add(sink2)

	_ = multi.OnEvent(context.Background(), &multiTestEvent{eventType: EventActivityStarted})

	if sink1.count != 1 || sink2.count != 1 {
		t.Errorf("both should receive event after Add")
	}
}

type countingSubscriber struct {
	count int
	err   error
}

func (c *countingSubscriber) OnEvent(_ context.Context, _ Event) error {
	c.count++
	return c.err
}

type multiTestEvent struct {
	eventType string
}

func (e *multiTestEvent) GetEventType() string { return e.eventType }
