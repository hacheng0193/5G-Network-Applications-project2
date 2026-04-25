package rms

import (
	"sync"
	"fmt"
	"github.com/free5gc/util/fsm"
	amf_context "github.com/free5gc/amf/internal/context"
	"bytes"
	"encoding/json"
	"net/http"
)
type Subscription struct {
	SubId     string `json:"subId"`
	UeId      string `json:"ueId"`
	NotifyUri string `json:"notifyUri"`
}
type UeRMNotif struct {
	SubId     string `json:"subId"`
	UeId      string `json:"ueId"`
	PrevState string `json:"from"`
	CurrState string `json:"to"`
}

type CustomizedRMS struct {
	// implement your customized RMS fields here
	mu            sync.RWMutex
	subscriptions map[string]Subscription // key: subId
	nextID int
}
// RMS is a self-contained module that manages its own subscription data internally
var (
	defaultRMS *CustomizedRMS
	once      sync.Once
)
func NewRMS() *CustomizedRMS {
	once.Do(func() {
		defaultRMS = &CustomizedRMS{
			subscriptions: make(map[string]Subscription),
			nextID:        1,
		}
	})
	return defaultRMS
}

// this will cause blank rms table after recall NewRMS
/*
func NewRMS(
// implement your customized RMS initialization here
) *CustomizedRMS {
	return &CustomizedRMS{
		subscriptions: make(map[string]Subscription),
		nextID: 1,
	}
}
*/
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
	rms.mu.RLock()
	defer rms.mu.RUnlock()
	// implement your customized RMS query logic here
	sub, exists := rms.subscriptions[subid]
	return sub, exists
}

func (rms *CustomizedRMS) QueryAll() []Subscription {
	rms.mu.RLock()
	defer rms.mu.RUnlock()
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

func (rms *CustomizedRMS) QueryByUeId(ueId string) []Subscription{
	rms.mu.RLock()
	defer rms.mu.RUnlock()

	result := make([]Subscription, 0)

	for _, sub := range rms.subscriptions {
		if sub.UeId == ueId {
			result = append(result, sub)
		}
	}

	return result

}


func (rms *CustomizedRMS) HandleEvent(state *fsm.State, event fsm.EventType, args fsm.ArgsType, trans fsm.Transition) {
	// implement your customized RMS logic here
	// fmt.Println("===== CustomizedRMS HandleEvent called =====")

	// Handle the event
	_ = state
	_ = event

	var ueId string

	if v, ok := args["AMF Ue"]; ok {
		if amfUe, ok := v.(*amf_context.AmfUe); ok && amfUe != nil {
			ueId = amfUe.Supi
		}
	}

	if ueId == "" {
		return
	}
	
	// find subs by ueid
	subs := rms.QueryByUeId(ueId)
	if len(subs) == 0 {
		return
	}
	/*
	fmt.Println("UE ID:", ueId)
	fmt.Println("===== Existing Subscriptions =====")
	allSubs := rms.QueryAll()
	for _, sub := range allSubs {
		fmt.Printf("SubId: %s, UeId: %s, NotifyUri: %s\n",
			sub.SubId, sub.UeId, sub.NotifyUri)
	}
	fmt.Println("===== Matched UE Subscriptions =====")
	*/
	/*
	for _, sub := range subs {
		fmt.Printf("SubId: %s, UeId: %s, NotifyUri: %s\n",sub.SubId, sub.UeId, sub.NotifyUri)
	}
	*/
	prevState := fmt.Sprintf("%v", trans.From)
	currState := fmt.Sprintf("%v", trans.To)

	for _, sub := range subs {
		notif := UeRMNotif{
			SubId:     sub.SubId,
			UeId:      sub.UeId,
			PrevState: prevState,
			CurrState: currState,
		}
	
		data, err := json.Marshal(notif)
		if err != nil {
			// fmt.Println("marshal notification failed:", err)
			continue
		}
	
		resp, err := http.DefaultClient.Post(
			sub.NotifyUri,
			"application/json",
			bytes.NewBuffer(data),
		)
		if err != nil {
			// fmt.Println("RMS notify failed:", err)
			continue
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
}
