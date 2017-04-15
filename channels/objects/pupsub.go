package objects

// Subscriber -
type (

	// SubscriberEntries -
	SubscriberEntries []SubscriberEntry

	// SubscriberEntry -
	SubscriberEntry struct {
		Channel FleaCommandChannel
	}

	// PubSubController -
	PubSubController struct {
		SubscriptionEntries map[string]SubscriberEntries
	}
)

// Notify -
func (p *PubSubController) Notify(topic string, command FleaCommand) {
	for _, se := range p.SubscriptionEntries[topic] {
		se.Channel <- command
	}
}

// Register -
func (p *PubSubController) Register(topic string, channel FleaCommandChannel) {
	r := append(
		p.SubscriptionEntries[topic],
		SubscriberEntry{channel})
	p.SubscriptionEntries[topic] = r
}
