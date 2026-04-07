package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/thetramp22/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("1 argument required: <time_between_reqs>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error parsing duration: %v", err)
	}

	fmt.Printf("collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed to scrape: %v", err)
	}

	s.db.MarkFeedFetched(context.Background(), nextFeed.ID)

	fetchedFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error getting feed: %v", err)
	}

	err = saveFeedToDB(s, *fetchedFeed, nextFeed.ID)
	if err != nil {
		return err
	}

	return nil
}

func saveFeedToDB(s *state, feed RSSFeed, feedId uuid.UUID) error {
	for _, item := range feed.Channel.Item {
		// Parse using RFC1123Z (includes numeric timezone offset)
		pubTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			// Fallback: try RFC1123 if the string uses named timezones like "MST"
			pubTime, err = time.Parse(time.RFC1123, item.PubDate)
		}

		nullTime := sql.NullTime{
			Time:  pubTime,
			Valid: err == nil,
		}

		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: nullTime,
			FeedID:      feedId,
		})
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code != "23505" {
				fmt.Println(err)
			}
		} else if err != nil {
			return err
		} else {
			fmt.Printf("post saved to db: %v\n", post.Title)
		}
	}
	return nil
}
