package dao

import (
	"context"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestResolveAccountID(t *testing.T) {
	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar?readPreference=primary&ssl=false"))
	if err != nil {
		log.Fatalf("cannot connect mongodb: %v", err)
	}

	m := NewMongo(mc.Database("coolcar"))

	id, err := m.ResolveAccountID(c, "123")
	if err != nil {
		log.Fatalf("faild resolve account id for 123: %v", err)
	} else {
		want := "628e3c48b556eba50da86d57"
		if id != want {
			t.Errorf("resolve account id: want: %q; got: %q", want, id)
		}
	}
}
