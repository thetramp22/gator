package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/thetramp22/blog_aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("takes only 1 optional argument: <limit>")
	}
	limit := 2
	if len(cmd.args) == 1 {
		limitString := cmd.args[0]
		i, err := strconv.Atoi(limitString)
		if err != nil {
			return fmt.Errorf("could not convert '%v' to interger", limitString)
		}
		limit = i
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error getting posts: %v", err)
	}

	for _, post := range posts {
		printPost(post)
	}

	return nil
}

func printPost(post database.Post) {
	fmt.Printf("Title: %v\n", post.Title)
	fmt.Printf("URL: %v\n", post.Url)
	if post.Description.Valid == true {
		fmt.Printf("Description: %v\n", post.Description.String)
	}
	if post.PublishedAt.Valid == true {
		fmt.Printf("Publish Date: %v\n", post.PublishedAt.Time.Format("Mon Jan 2"))
	}
	fmt.Println("==========")
}
