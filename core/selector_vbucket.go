package core

import (
	"github.com/nicholaskh/mysql-cluster/config"
)

const (
	ServerActive      = "active"  // fully oprational
	ServerDead        = "dead"    // fully non-operational
	ServerPending     = "pending" // blocks clients, receives replicas
	ServerReplicating = "replica" // dead to clients, receive replicas, transfer vbuckets from one server to another
)

// vBucket-aware mysql cluster client selector.
//
// The vBucket mechanism provides a layer of indirection between the hashing algorithm
// and the server responsible for a given key.
//
// This indirection is useful in managing the orderly transition from one cluster configuration to
// another, whether the transition was planned (e.g. adding new servers to a cluster) or unexpected (e.g. a server failure)
//
// Every key belongs to a vBucket, which maps to a server instance
// The number of vBuckets is a fixed number
//
// key ---------------> server
// h(key) -> vBucket -> server
//
// servers = ['server1:11211', 'server2:11211', 'server3:11211']
// vbuckets = [0, 0, 1, 1, 2, 2]
// server_for_key(key) = servers[vbuckets[hash(key) % vbuckets.length]]
//
// how to add a new server:
// push config to all clients -> to make the new server useful, transfer vbuckets from one server to another, set
// them to ServerPending state on the receiving server
//
// The vBucket-Server map is updated internally: transmitted from server to all cluster participants:
// servers, clients and proxies

type VbucketSelector struct {
	serverPool ServerPool
	config     *config.VbucketShardingConfig
}

func newVbucketSelector(config *config.VbucketShardingConfig, serverPool ServerPool) *VbucketSelector {
	this := new(VbucketSelector)
	this.config = config
	this.serverPool = serverPool

	return this
}

// for 3 servers, vbuckets = [0, 1, 2, 0, 1, 2 ...]
// server = (bucket -1) % servers.length
func (this *VbucketSelector) PickServer(pool string, hintId int, sql string) (server *mysql, ex error) {
	if _, exists := this.serverPool[pool]; !exists {
		ex = ErrInvalidPool
		return
	}

	serverId := (hintId%this.config.VbucketBaseNumber - 1) % len(this.serverPool[pool])
	server, _ = this.serverPool[pool][serverId]
	return
}
