package service

import (
	"context"

	"github.com/pkg/errors"
)

func (m *Service) SendOrderNotifications(ctx context.Context) error {
	outbox, err := m.LOMSRepo.GetOutbox(ctx)
	if err != nil {
		return err
	}
	var result error
	for _, msg := range outbox {
		err = m.NotificationsSender.SendNotification(ctx, msg)
		if err != nil {
			if result != nil {
				result = errors.WithMessage(result, err.Error())
			} else {
				result = err
			}
		} else {
			err = m.LOMSRepo.DeleteOutbox(ctx, msg.MsgID)
			if err != nil {
				if result != nil {
					result = errors.WithMessage(result, err.Error())
				} else {
					result = err
				}
			}
		}
	}
	return result
}
