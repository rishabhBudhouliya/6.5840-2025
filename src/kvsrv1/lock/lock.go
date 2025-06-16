package lock

import (
	"fmt"

	"6.5840/kvsrv1/rpc"
	kvtest "6.5840/kvtest1"
)

type Lock struct {
	// IKVClerk is a go interface for k/v clerks: the interface hides
	// the specific Clerk type of ck but promises that ck supports
	// Put and Get.  The tester passes the clerk in when calling
	// MakeLock().
	ck  kvtest.IKVClerk
	id  string
	key string
	// You may add code here
}

// The tester calls MakeLock() and passes in a k/v clerk; your code can
// perform a Put or Get by calling lk.ck.Put() or lk.ck.Get().
//
// Use l as the key to store the "lock state" (you would have to decide
// precisely what the lock state is).
// assume lock state is ready or in use
func MakeLock(ck kvtest.IKVClerk, l string) *Lock {
	lk := &Lock{ck: ck}
	lk.id = kvtest.RandValue(8)
	lk.key = l
	// randKey := l + kvtest.RandValue(8)
	// lk.ck.Put(randKey, "ready", rpc.Tversion(0))
	// lk.key = randKey
	return lk
}

/*
key -> value,version
key is a constant, value is client's id and version is internally managed by the server
once acquired, we reject any other put because the value is not theirs
once released, we remove the client id from the value
*/
func (lk *Lock) Acquire() {
	for {
		value, version, getErr := lk.ck.Get(lk.key)
		// lock free to be acquired
		if getErr == rpc.OK && value == "" {
			putErr := lk.ck.Put(lk.key, lk.id, version)
			if putErr == rpc.ErrVersion {
				continue
			} else {
				break
			}
		}
		if getErr == rpc.OK && value != lk.id {
			continue
		}
		if getErr == rpc.ErrNoKey || value == lk.id {
			putErr := lk.ck.Put(lk.key, lk.id, version)
			if putErr == rpc.ErrVersion {
				continue
			} else {
				break
			}
		}
	}
}

func (lk *Lock) Release() {
	_, version, getErr := lk.ck.Get(lk.key)
	if getErr == rpc.OK {
		lk.ck.Put(lk.key, "", version)
	} else {
		fmt.Print("lock can't be released")
	}
}

/*
for {
		value, version, err := lk.ck.Get(lk.key)
		if err == rpc.ErrNoKey {
			lk.ck.Put(lk.key, strconv.Itoa(0), rpc.Tversion(0))
			return
		}
		v, _ := strconv.Atoi(value)
		if v != int(version) {
			continue
		} else {
			resErr2 := lk.ck.Put(lk.key, value, version)
			if resErr2 != nil {

			}
			return
		}
	}
*/

/*
lk.mu.Lock()
	defer lk.mu.Unlock()
	_, version, err := lk.ck.Get(lk.key)
	if err == rpc.OK {
		lk.ck.Put(lk.key, strconv.Itoa(int(version)+1), version)
	}
*/
