package limiter

import (
	"context"
	"time"
)

// Простейший вариант рейт лимитера, без burst, ограничивает пропускание строго заданными интервалами.
// Построен на базе time.Ticker.
// Отправляет сообщения в канал Limiter.C не чаще установленного интервала.
// В сообщении содержится время, когда оно отправлено в канал. По нему можно анализировать использование полосы пропускания.

type Limiter struct {
	C      chan time.Time
	ctx    context.Context
	ticker *time.Ticker
}

// New создает новый лимитер и сразу его запускает
// В ctx передается контекст используемый для остановки лимитера при graceful shutdown
// В d передается минимальный период, который будет ограничивать лимитер
// Для использвания лимитера, необходимо перед лимитируемым кодом читать канал Limiter.C

func New(ctx context.Context, d time.Duration) *Limiter {
	result := Limiter{
		C:      make(chan time.Time),
		ctx:    ctx,
		ticker: time.NewTicker(d),
	}
	go func(ctx context.Context) {
		defer result.ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-result.ticker.C:
				select {
				case <-ctx.Done():
					return
				case result.C <- time.Now():
				}
			}
		}
	}(ctx)
	return &result
}

// Stop останавливает лимитер

func (l *Limiter) Stop() {
	l.ticker.Stop()
}
