package objects

// Subscriber -
type (

	// SubscriberEntries -
	SubscriberEntries []SubscriberEntry

	// SubscriberEntry -
	SubscriberEntry struct {
		Channel chan bool
	}

	// PubSubController -
	PubSubController struct {
		SubscriptionEntries map[string]SubscriberEntries
	}
)

// Notify -
func (p *PubSubController) Notify(topic string) {
	for _, se := range p.SubscriptionEntries[topic] {
		se.Channel <- true
	}
}

// Register -
func (p *PubSubController) Register(topic string, channel chan bool) {
	r := append(
		p.SubscriptionEntries[topic],
		SubscriberEntry{channel})
	p.SubscriptionEntries[topic] = r
}
