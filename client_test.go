package wxhelper

import (
	"context"
	"testing"
)

func TestClient_CheckLogin(t *testing.T) {
	messageReceiver, err := NewHttpMessageRetriever(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	msg, err := messageReceiver.RetrieveMessage()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			select {
			case m := <-msg:
				t.Logf("received message: %+v", m)
			}
		}
	}()
	select {}
}
