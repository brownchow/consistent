package consistent

import (
	"testing"
)

func TestBasicOp(t *testing.T) {
	ch := NewConsistentHashing()
	ch.Add("server1")
	// map 自动扩容？？ 可很多key我也没添加啊
	t.Logf("server hashes: %d, circleRing size: %d", ch.hashSortedKeys, len(ch.circleRing))
	for serverHash, server := range ch.circleRing {
		t.Logf("serverHash: %d, server: %s", serverHash, server)
	}

	ch.Add("server2")
	t.Logf("server hashes: %d, circleRing size: %d", ch.hashSortedKeys, len(ch.circleRing))
	for serverHash, server := range ch.circleRing {
		t.Logf("serverHash: %d, server: %s", serverHash, server)
	}

	ch.Add("server3")
	t.Logf("server hashes: %d, circleRing size: %d", ch.hashSortedKeys, len(ch.circleRing))
	for serverHash, server := range ch.circleRing {
		t.Logf("serverHash: %d, server: %s", serverHash, server)
	}

	t.Logf("noedes are %s \n", ch.ListNodes())
	targetObj := []string{"client1", "client2", "client3", "client4", "client5", "client6"}
	for _, client := range targetObj {
		server, _ := ch.Get(client)
		t.Logf("client: %s in server: %s\n", client, server)
	}

	t.Logf("------------------------------------------")
	ch.Add("server4")
	ch.Add("serevr5")
	for _, client := range targetObj {
		server, err := ch.Get(client)
		if err == nil {
			t.Logf("client: %s in server : %s \n", client, server)
		}
	}
	t.Logf("\n")
}
