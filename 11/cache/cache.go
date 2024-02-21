package cache

import (
	"L2/11/domain"
	"L2/11/models"
	"log"
	"sync"
	"time"
)

type Cache struct {
	m    sync.Mutex
	pool map[string][]*models.Event
}

func NewCache() *Cache {
	return &Cache{
		m:    sync.Mutex{},
		pool: map[string][]*models.Event{},
	}
}

func (c *Cache) Create(e *models.Event) {
	c.m.Lock()
	defer c.m.Unlock()
	c.pool[e.Date] = append(c.pool[e.Date], e)
}

func (c *Cache) ReadDay(date string) ([]*models.Event, bool) {
	c.m.Lock()
	defer c.m.Unlock()
	val, ok := c.pool[date]
	if ok {
		return val, true
	}
	log.Println("Err at getting", date, ": no events detected")
	return []*models.Event{}, false
}

func (c *Cache) ReadWeek(date string) ([]models.Event, bool) {
	c.m.Lock()
	defer c.m.Unlock()
	result := []models.Event{}

	today, errParseToday := time.Parse("2006-01-02", date)

	if errParseToday != nil {
		log.Println("Error at parsing date: ", errParseToday)
		return nil, false
	}

	for _, ev := range c.pool {

		for _, event := range ev {

			evDate, errParseEv := time.Parse("2006-01-02", event.Date)
			if errParseEv != nil {
				log.Println("Error at parsing date: ", errParseEv)
				return nil, false
			}

			if evDate.After(today) && evDate.Before(today.Add(time.Hour*168)) {
				result = append(result, *event)
			}
		}

	}

	return result, true
}

func (c *Cache) ReadMonth(date string) ([]models.Event, bool) {
	c.m.Lock()
	defer c.m.Unlock()
	result := []models.Event{}

	today, errParseToday := time.Parse("2006-01-02", date)

	if errParseToday != nil {
		log.Println("Error at parsing date: ", errParseToday)
		return nil, false
	}

	for _, ev := range c.pool {

		for _, event := range ev {

			evDate, errParseEv := time.Parse("2006-01-02", event.Date)
			if errParseEv != nil {
				log.Println("Error at parsing date: ", errParseEv)
				return nil, false
			}

			if evDate.After(today) && evDate.Before(today.Add(time.Hour*24*30)) {
				result = append(result, *event)
			}
		}

	}

	return result, true
}

func (c *Cache) Delete(date, time string) {
	c.m.Lock()
	defer c.m.Unlock()
	ev := []*models.Event{}
	for i := 0; i < len(c.pool[date]); i += 1 {
		if c.pool[date][i].Time != time {
			ev = append(ev, c.pool[date][i])
		}
	}
	c.pool[date] = ev
}

func (c *Cache) Update(e models.Event, newDate, newTime string) {
	for _, event := range c.pool[e.Date] {
		if event.Time == e.Time && event.Date == e.Date {
			UpdatedEvent := models.Event{}
			UpdatedEvent.UserId = e.UserId
			if newTime != "" {
				UpdatedEvent.Time = newTime
			}
			if newDate != "" {
				UpdatedEvent.Date = newDate
				c.Delete(e.Date, e.Time)
				newEvent := domain.NewEvent(UpdatedEvent.Date, UpdatedEvent.Time, UpdatedEvent.UserId)
				c.Create(newEvent)
			}
			return
		}
	}
}
