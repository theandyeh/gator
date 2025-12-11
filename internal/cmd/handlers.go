package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/theandyeh/gator/internal/app"
	"github.com/theandyeh/gator/internal/database"
	"github.com/theandyeh/gator/internal/rss"
)

//Middleware

func MiddlewareLoggedIn(handler func(s *app.State, cmd Command, user database.User) error) func(*app.State, Command) error {
	return func(s *app.State, c Command) error {
		u, err := s.Db.GetUser(context.Background(), s.Cfg.Current_db_user)
		if err != nil {
			return fmt.Errorf("middleware logged in error: %w", err)
		}

		return handler(s, c, u)
	}
}

//Handlers

func HandlerLogin(s *app.State, c Command) error {
	if len(c.Args) < 1 {
		return fmt.Errorf("login handler error: no username provided for login command")
	}

	if _, err := s.Db.GetUser(context.Background(), c.Args[0]); err != nil {
		return fmt.Errorf("login handler error: %w", err)
	}

	if err := s.Cfg.SetUser(c.Args[0]); err != nil {
		return fmt.Errorf("login handler error: %w", err)
	}

	fmt.Printf("Successfuly set DB User: %s", c.Args[0])
	return nil
}

func HandlerRegister(s *app.State, c Command) error {
	if len(c.Args) < 1 {
		return fmt.Errorf("register handler error: no username provided for register command")
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      c.Args[0],
	}

	if u, err := s.Db.GetUser(context.Background(), c.Args[0]); err == nil && u.Name == c.Args[0] {
		fmt.Printf("register handler error: user %s already exists\n", c.Args[0])
		os.Exit(1)
	} else if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return fmt.Errorf("register handler error: %w", err)
		}
	}

	if _, err := s.Db.CreateUser(context.Background(), user); err != nil {
		return fmt.Errorf("register handler error: %w", err)
	}

	s.Cfg.Current_db_user = c.Args[0]
	s.Cfg.SetUser(c.Args[0])

	fmt.Printf("Successfuly registered and set DB User: %s\n", c.Args[0])
	fmt.Printf("User details:\n %v", user)
	return nil
}

func HandlerReset(s *app.State, c Command) error {
	if err := s.Db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("reset handler error: %w", err)
	}

	s.Cfg.Current_db_user = ""
	s.Cfg.SetUser("")

	fmt.Println("Successfully reset database and cleared current user in config")
	return nil
}

func HandlerUsers(s *app.State, c Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("users handler error: %w", err)
	}

	fmt.Println("Registered Users:")
	for _, user := range users {
		if user.Name == s.Cfg.Current_db_user {
			fmt.Printf("- %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("- %s\n", user.Name)
	}

	return nil
}

func HandlerAgg(s *app.State, c Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("agg handler error: %w", err)
	}

	fmt.Printf("Feed Title: %s\n", feed.Channel.Title)
	fmt.Printf("Feed Link: %s\n", feed.Channel.Link)
	fmt.Printf("Feed Description: %s\n", feed.Channel.Description)

	fmt.Println("Feed Items:")
	for _, item := range feed.Channel.Item {
		fmt.Printf("- Title: %s\n", item.Title)
		fmt.Printf("  Link: %s\n", item.Link)
		fmt.Printf("  Description: %s\n", item.Description)
		fmt.Printf("  Published Date: %s\n", item.PubDate)
		fmt.Println("--------------------------")
	}
	return nil
}

func HandlerAddFeed(s *app.State, c Command, user database.User) error {
	if len(c.Args) < 2 {
		return fmt.Errorf("addfeed handler error: not enough arguments provided, expected feed name and URL")
	}

	if !isValidUrl(c.Args[1]) {
		return fmt.Errorf("addfeed handler error: invalid feed URL provided, check args (addfeed <name> <url>")
	}

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      c.Args[0],
		Url:       c.Args[1],
		UserID:    user.ID,
	}

	if _, err := s.Db.CreateFeed(context.Background(), feed); err != nil {
		return fmt.Errorf("addfeed handler error creating feed: %w", err)
	}

	followP := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}

	if _, err := s.Db.CreateFeedFollow(context.Background(), followP); err != nil {
		return fmt.Errorf("addfeed handler error creating feed follow: %w", err)
	}

	fmt.Printf("Successfully added and followed feed:\n%v", feed)

	return nil
}

func HandlerFeeds(s *app.State, c Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("feeds handler error retrieving data: %w", err)
	}

	fmt.Println("Registered Feeds:")
	for _, feed := range feeds {
		fmt.Printf("- Name: %s\n", feed.Name)
		fmt.Printf("  URL: %s\n", feed.Url)
		fmt.Printf("  User Name: %s\n", feed.Username.String)
		fmt.Println("--------------------------")
	}

	return nil
}

func HandlerFollow(s *app.State, c Command, user database.User) error {
	if len(c.Args) < 1 {
		return fmt.Errorf("follow handler error: no feed url provided to follow")
	}

	feedUrl := c.Args[0]
	if !isValidUrl(feedUrl) {
		return fmt.Errorf("follow handler error: invalid feed URL provided")
	}

	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("follow handler error retrieving feeds: %w", err)
	}

	var feedRegistered bool

	for _, f := range feeds {
		if f.Url == feedUrl {
			feedRegistered = true

			followP := database.CreateFeedFollowParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				FeedID:    f.ID,
				UserID:    user.ID,
			}

			followRow, err := s.Db.CreateFeedFollow(context.Background(), followP)
			if err != nil {
				return fmt.Errorf("follow handler error creating feed follow: %w", err)
			}

			fmt.Printf("Successfully followed feed:\n%v", followRow)
			return nil
		}
	}

	if !feedRegistered {
		feed, err := rss.FetchFeed(context.Background(), feedUrl)
		if err != nil {
			return fmt.Errorf("follow handler error fetching feed: %w", err)
		}

		newFeed := database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      feed.Channel.Title,
			Url:       feedUrl,
			UserID:    user.ID,
		}

		createdFeed, err := s.Db.CreateFeed(context.Background(), newFeed)
		if err != nil {
			return fmt.Errorf("follow handler error creating new feed: %w", err)
		}

		followP := database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FeedID:    createdFeed.ID,
			UserID:    user.ID,
		}

		followRow, err := s.Db.CreateFeedFollow(context.Background(), followP)
		if err != nil {
			return fmt.Errorf("follow handler error creating feed follow for new feed: %w", err)
		}

		fmt.Printf("Successfully added and followed new feed:\n%v", followRow)
	}
	return nil
}

func HandlerFollowing(s *app.State, c Command, user database.User) error {
	following, err := s.Db.GetFeedFollowsByUserID(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("following handler error retrieving data: %w", err)
	}

	fmt.Printf("Feeds followed by user %s:\n", user.Name)
	for _, follow := range following {
		fmt.Printf("- Feed Name: %s\n", follow.FeedName)
		fmt.Printf("  Feed ID: %s\n", follow.FeedID)
		fmt.Printf("  User Name: %s\n", follow.UserName)
		fmt.Println("--------------------------")
	}

	return nil
}

//HELPERS

func isValidUrl(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}
