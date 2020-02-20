package notifications

//  holds a queue of message callback responses to be sent yo users.
type Pubsub interface {
	// Subscribe to incoming callbacks
	Subscribe()
	// Publish incoming callbacks
	Publish()
}
