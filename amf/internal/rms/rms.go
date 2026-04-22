package rms

import "github.com/free5gc/util/fsm"

type CustomizedRMS struct {
	// implement your customized RMS fields here
	subscriptions map[string]Subscription // key: subId
}

func NewRMS(
// implement your customized RMS initialization here
) fsm.RMS {
	return &CustomizedRMS{
		make(map[string]Subscription),
	}
}

func (rms *CustomizedRMS) Add(sub Subscription) {
	// implement your customized RMS add logic here
	rms.subscriptions[subscription.SubID] = sub
}

func (rms *CustomizedRMS) Modify(subid string, sub Subscription) bool {
	// implement your customized RMS modify logic here
	if _, exists := rms.subscriptions[subid]; exists {
		rms.subscriptions[subid] = sub
		return true
	}
	return false
}

func (rms *CustomizedRMS) Query(subid string) (Subscription, bool) {
	// implement your customized RMS query logic here
	sub, exists := rms.subscriptions[subid]
	return sub, exists
}

func (rms *CustomizedRMS) QueryAll() []Subscription {
	// implement your customized RMS query all logic here
	subs := make([]Subscription, 0, len(rms.subscriptions))
	for _, sub := range rms.subscriptions {
		subs = append(subs, sub)
	}
	return subs
}



func (rms *CustomizedRMS) Delete(subid string) bool{
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
