package cache

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Кеш для отправки ответа аукциона
type CacheResponseDecorator struct {
	updateInterval time.Duration
	ar             responser
	isActive       bool

	sync.RWMutex
	resp interface{}
	err  error
}

func NewCacheResponseDecorator(updateInterval time.Duration, ar responser) *CacheResponseDecorator {
	return &CacheResponseDecorator{
		updateInterval: updateInterval,
		ar:             ar,
	}
}

// Делаем recover при панике на cache,
// когда отрабатывает recover - сache падает
// и берём данные из сборщика response'а
func (c *CacheResponseDecorator) handlePanic() {
	if r := recover(); r != nil {
		// Делаем данные кэша недоступными
		c.setIsActive(false)
		log.Warn().Msgf("recovered from auction cache panic: %s", r)
	}
}

// Делаем данные кэша недоступными или доступными
func (c *CacheResponseDecorator) setIsActive(status bool) {
	c.Lock()
	defer c.Unlock()
	c.isActive = status
}

func (c *CacheResponseDecorator) makeUpdate(ctx context.Context) {
	// Если нет ошибки, кэш - активен иначе неактивен
	if err := c.update(ctx); err != nil {
		c.setIsActive(false)
		log.Err(err).Msg("CacheResponseDecorator.makeUpdate")
		return
	}
	c.setIsActive(true)
}

// Включить обновление response'а с интервалом
func (c *CacheResponseDecorator) Start(ctx context.Context) {
	// Обновим данные
	go func() {
		c.makeUpdate(ctx)
		log.Info().Msg("cached_updater started")
		// Обновляем данные с итервалом
		for range time.NewTicker(c.updateInterval).C {
			c.makeUpdate(ctx)
		}
	}()
}

// Получение response'а
func (c *CacheResponseDecorator) Get(ctx context.Context) (interface{}, error) {
	c.RLock()
	// Если cache не активирован, соберём респонс
	if !c.isActive || c.resp == nil {
		c.RUnlock()
		log.Info().Msg("CacheResponseDecorator.Get get_uncahed_response")
		return c.ar.Get(ctx)
	}
	defer c.RUnlock()
	log.Info().Msg("CacheResponseDecorator.Get get_cache_response")
	return c.resp, c.err
}

// Обновление response'а
func (c *CacheResponseDecorator) update(ctx context.Context) error {
	// Делаем recover при панике
	defer c.handlePanic()
	// Собираем response
	resp, err := c.ar.Get(ctx)
	c.Lock()
	c.resp, c.err = resp, err
	c.Unlock()
	return err
}
