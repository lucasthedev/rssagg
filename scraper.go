package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/lucasthedev/rssagg/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Resgatando em %v goroutines a cada %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedToFetch(context.Background(), int32(concurrency))

		if err != nil {
			log.Println("erro ao resgatar feeds")
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}

}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)

	if err != nil {
		log.Println(" Erro ao marcar fetched feed", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println(" Erro ao converter feed url", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Encontrou o post ", item.Title, " no feed ", feed.Name)
	}
	log.Printf("Feed %s coletado, %v posts encontrados", feed.Name, len(rssFeed.Channel.Item))
}
