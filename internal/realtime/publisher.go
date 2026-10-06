package realtime

import "context"

type Publisher struct {
	manager *Manager
}

func NewPublisher(manager *Manager) *Publisher {
	return &Publisher{
		manager: manager,
	}
}

func (p *Publisher) Publish(
	ctx context.Context,
	userID string,
	event Event,
) error {
	return p.manager.SendEvent(ctx, userID, event)
}
