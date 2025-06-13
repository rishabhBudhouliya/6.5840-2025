package kvsrv

import (
	"log"
	"sync"

	"6.5840/kvsrv1/rpc"
	"6.5840/labrpc"
	tester "6.5840/tester1"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type KVServer struct {
	mu sync.Mutex
	m  map[string]*Value // key value string
}

type Value struct {
	value   string
	version rpc.Tversion
}

func MakeKVServer() *KVServer {
	kv := &KVServer{}
	kv.m = make(map[string]*Value)
	return kv
}

// Get returns the value and version for args.Key, if args.Key
// exists. Otherwise, Get returns ErrNoKey.
func (kv *KVServer) Get(args *rpc.GetArgs, reply *rpc.GetReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	if kv.m[args.Key] != nil {
		reply.Value = kv.m[args.Key].value
		reply.Version = rpc.Tversion(kv.m[args.Key].version)
	} else {
		reply.Err = rpc.ErrNoKey
		return
	}
	reply.Err = rpc.OK
}

// Update the value for a key if args.Version matches the version of
// the key on the server. If versions don't match, return ErrVersion.
// If the key doesn't exist, Put installs the value if the
// args.Version is 0, and returns ErrNoKey otherwise.
func (kv *KVServer) Put(args *rpc.PutArgs, reply *rpc.PutReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	v := kv.m[args.Key]
	if args.Version > 0 && v == nil {
		reply.Err = rpc.ErrNoKey
		return
	}
	if v == nil && args.Version == 0 {
		kv.m[args.Key] = &Value{value: args.Value, version: 1}
	} else if args.Version == v.version {
		v.value = args.Value
		v.version = v.version + 1
		kv.m[args.Key] = v
	} else if args.Version != v.version {
		reply.Err = rpc.ErrVersion
		return
	}
	reply.Err = rpc.OK
}

// You can ignore Kill() for this lab
func (kv *KVServer) Kill() {
}

// You can ignore all arguments; they are for replicated KVservers
func StartKVServer(ends []*labrpc.ClientEnd, gid tester.Tgid, srv int, persister *tester.Persister) []tester.IService {
	kv := MakeKVServer()
	return []tester.IService{kv}
}
