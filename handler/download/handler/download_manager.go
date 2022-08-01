package downhandler

import "time"

func (manager *ManagerDownload) AddEvent(u *DownloadEvent) {
	manager.mu.Lock()
	timer := time.NewTimer(5 * time.Minute)
	go func() {
		<-timer.C
		DownloadManager.DeleteEvent(u.Id)
	}()
	defer manager.mu.Unlock()
	manager.m[u.Id] = u
}

func (manager *ManagerDownload) QueryEvent(eventId string) (u *DownloadEvent, ok bool) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if val, ok := manager.m[eventId]; ok {
		return val, ok
	} else {
		return nil, false
	}
}

func (manager *ManagerDownload) DeleteEvent(eventId string) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	delete(manager.m, eventId)
}
