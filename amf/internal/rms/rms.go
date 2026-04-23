package rms

import (
	"sync"
	"fmt"
	"github.com/free5gc/util/fsm"
)
type Subscription struct {
	SubId     string `json:"subId"`
	UeId      string `json:"ueId"`
	NotifyUri string `json:"notifyUri"`
}
type CustomizedRMS struct {
	// implement your customized RMS fields here
	mu            sync.RWMutex
	subscriptions map[string]Subscription // key: subId
	nextID int
}

func NewRMS(
// implement your customized RMS initialization here
) *CustomizedRMS {
	return &CustomizedRMS{
		subscriptions: make(map[string]Subscription),
		nextID: 1,
	}
}

func (rms *CustomizedRMS) Add(sub Subscription) Subscription{
	rms.mu.Lock()
	defer rms.mu.Unlock()
	// implement your customized RMS add logic here
	sub.SubId = fmt.Sprintf("sub-%03d", rms.nextID)
    rms.nextID++
    rms.subscriptions[sub.SubId] = sub
    return sub
}

func (rms *CustomizedRMS) Modify(subid string, sub Subscription) bool {
	rms.mu.Lock()
	defer rms.mu.Unlock()
	// implement your customized RMS modify logic here
	if _, exists := rms.subscriptions[subid]; exists {
		sub.SubId = subid
		rms.subscriptions[subid] = sub
		return true
	}
	return false
}
func (rms *CustomizedRMS) Put(subid string, sub Subscription) (Subscription, bool) {
	rms.mu.Lock()
	defer rms.mu.Unlock()
	// implement your customized RMS put logic here
	_, existed := rms.subscriptions[subid]
	sub.SubId = subid
	rms.subscriptions[subid] = sub
	return sub, existed
}
func (rms *CustomizedRMS) Query(subid string) (Subscription, bool) {
	rms.mu.Lock()
	defer rms.mu.Unlock()
	// implement your customized RMS query logic here
	sub, exists := rms.subscriptions[subid]
	return sub, exists
}

func (rms *CustomizedRMS) QueryAll() []Subscription {
	rms.mu.Lock()
	defer rms.mu.Unlock()
	// implement your customized RMS query all logic here
	subs := make([]Subscription, 0, len(rms.subscriptions))
	for _, sub := range rms.subscriptions {
		subs = append(subs, sub)
	}
	return subs
}



func (rms *CustomizedRMS) Delete(subid string) bool{
	rms.mu.Lock()
	defer rms.mu.Unlock()
	// implement your customized RMS remove logic here
	if _, exists := rms.subscriptions[subid]; exists {
		delete(rms.subscriptions, subid)
		return true
	}
	return false
}
func (rms *CustomizedRMS) HandleEvent(state *fsm.State, event fsm.EventType, args fsm.ArgsType, trans fsm.Transition) {
	// implement your customized RMS logic here
	
}
