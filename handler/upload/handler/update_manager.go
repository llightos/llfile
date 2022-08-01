package handler

import "os"

func (manager *ManagerUpload) AddEvent(u *UploadEvent) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.m[u.Id] = u
}

func (manager *ManagerUpload) QueryEvent(eventId string) (u *UploadEvent, ok bool) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if val, ok := manager.m[eventId]; ok {
		return val, ok
	} else {
		return nil, false
	}
}

func (manager *ManagerUpload) DeleteEvent(eventId string) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	switch manager.m[eventId].ReadWriter.(type) {
	case *os.File:
		manager.m[eventId].ReadWriter.(*os.File).Close()
	}
	delete(manager.m, eventId)
}
