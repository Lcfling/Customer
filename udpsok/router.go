package udpsok

import "sync"

type Route struct {

	mutex        sync.Mutex
	clients      map[int64]*Client
}
