package main

import "context"

type PaymentTaken struct {
	PaymentID string
	Amount    int
}

type PaymentsHandler struct {
	repo *PaymentsRepository
}

func NewPaymentsHandler(repo *PaymentsRepository) *PaymentsHandler {
	return &PaymentsHandler{
		repo: repo,
	}
}

func (p *PaymentsHandler) HandlePaymentTaken(ctx context.Context, event *PaymentTaken) error {
	return p.repo.SavePaymentTaken(ctx, event)
}

type PaymentsRepository struct {
	payments []PaymentTaken
	dedup    map[string]struct{}
}

func (p *PaymentsRepository) Payments() []PaymentTaken {
	return p.payments
}

func NewPaymentsRepository() *PaymentsRepository {
	dedup := make(map[string]struct{})
	return &PaymentsRepository{dedup: dedup}
}

func (p *PaymentsRepository) SavePaymentTaken(ctx context.Context, event *PaymentTaken) error {
	_, ok := p.dedup[event.PaymentID]
	if ok {
		return nil
	}
	p.dedup[event.PaymentID] = struct{}{}
	p.payments = append(p.payments, *event)
	return nil
}
