package kvsrv

import (
	"time"

	"6.5840/kvsrv1/rpc"
	kvtest "6.5840/kvtest1"
	tester "6.5840/tester1"
)

type Clerk struct {
	clnt   *tester.Clnt
	server string
}

func MakeClerk(clnt *tester.Clnt, server string) kvtest.IKVClerk {
	ck := &Clerk{clnt: clnt, server: server}
	// You may add code here.
	return ck
}

// Get fetches the current value and version for a key.  It returns
// ErrNoKey if the key does not exist. It keeps trying forever in the
// face of all other errors.
//
// You can send an RPC with code like this:
// ok := ck.clnt.Call(ck.server, "KVServer.Get", &args, &reply)
//
// The types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. Additionally, reply must be passed as a pointer.
func (ck *Clerk) Get(key string) (string, rpc.Tversion, rpc.Err) {
	// You will have to modify this function.
	args := rpc.GetArgs{Key: key}
	reply := rpc.GetReply{}
	result := ck.clnt.Call(ck.server, "KVServer.Get", &args, &reply)
	if !result {
		retry := ck.RetryGet(key, &reply)
		if retry {
			return reply.Value, reply.Version, reply.Err
		}
	}
	return reply.Value, reply.Version, reply.Err
}

// Put updates key with value only if the version in the
// request matches the version of the key at the server.  If the
// versions numbers don't match, the server should return
// ErrVersion.  If Put receives an ErrVersion on its first RPC, Put
// should return ErrVersion, since the Put was definitely not
// performed at the server. If the server returns ErrVersion on a
// resend RPC, then Put must return ErrMaybe to the application, since
// its earlier RPC might have been processed by the server successfully
// but the response was lost, and the Clerk doesn't know if
// the Put was performed or not.
//
// You can send an RPC with code like this:
// ok := ck.clnt.Call(ck.server, "KVServer.Put", &args, &reply)
//
// The types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. Additionally, reply must be passed as a pointer.
func (ck *Clerk) Put(key, value string, version rpc.Tversion) rpc.Err {
	// You will have to modify this function.
	args := rpc.PutArgs{Key: key, Value: value, Version: version}
	reply := rpc.PutReply{}
	result := ck.clnt.Call(ck.server, "KVServer.Put", &args, &reply)
	// call failed, we want to retry put
	if !result {
		// retry Put until we get a result = true with either rpc ok or err
		return ck.RetryPut(key, value, version)
	}
	return reply.Err
}

func (ck *Clerk) RetryGet(key string, reply *rpc.GetReply) bool {
	for {
		args := rpc.GetArgs{Key: key}
		result := ck.clnt.Call(ck.server, "KVServer.Get", &args, &reply)
		if result && reply.Err == rpc.OK {
			return result
		}
	}
}

func (ck *Clerk) RetryPut(key string, value string, version rpc.Tversion) rpc.Err {
	for {
		args := rpc.PutArgs{Key: key, Value: value, Version: version}
		reply := rpc.PutReply{}
		result := ck.clnt.Call(ck.server, "KVServer.Put", &args, &reply)
		if result && reply.Err == rpc.ErrVersion {
			return rpc.ErrMaybe
		} else if result {
			return reply.Err
		} else {
			time.Sleep(100 * time.Millisecond)
			continue
		}
	}
}
