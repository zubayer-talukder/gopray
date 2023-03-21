package application

import (
	"context"

	"github.com/escalopa/gopray/pkg/core"
)

type PrayerRepository interface {
	StorePrayer(ctx context.Context, times core.PrayerTimes) error
	GetPrayer(ctx context.Context, day, month int) (core.PrayerTimes, error)
}

type SubscriberRepository interface {
	StoreSubscriber(ctx context.Context, id int) error
	RemoveSubscribe(ctx context.Context, id int) error
	GetSubscribers(ctx context.Context) ([]int, error)
}

type LanguageRepository interface {
	GetLang(ctx context.Context, id int) (string, error)
	SetLang(ctx context.Context, id int, lang string) error
}

type HistoryRepository interface {
	// Default prayers
	GetPrayerMessageID(ctx context.Context, userID int) (int, error)
	StorePrayerMessageID(ctx context.Context, userID int, messageID int) error
	// Gomaa
	GetGomaaMessageID(ctx context.Context, userID int) (int, error)
	StoreGomaaMessageID(ctx context.Context, userID int, messageID int) error
}

type Parser interface {
	ParseSchedule(ctx context.Context) error
}

type Notifier interface {
	// NotifyPrayers notifies subscribers about the upcoming prayer and when the prayer has started.
	// The first argument is a function that is called when the prayer is about to start.
	// The second argument is a function that is called when the prayer has started.
	NotifyPrayers(context.Context, func(id []int, name, time string), func(id []int, name string))
	// NotifyGomaa notifies subscribers about the gomaa prayer at the specified hour of friday.
	NotifyGomaa(context.Context, func(id []int, time string))
}
