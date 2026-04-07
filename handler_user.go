package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thetramp22/gator/internal/database"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("1 argument <name> required")
	}
	name := cmd.args[0]
	ctx := context.Background()
	newUser, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("user already exists: %v", err)
	}
	err = s.cfg.SetUser(newUser.Name)
	if err != nil {
		return err
	}
	fmt.Printf("user '%v' was created\n", newUser.Name)
	fmt.Printf("%+v\n", newUser)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("1 argument <username> required")
	}
	username := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user '%v' does not exist: %v", username, err)
	}
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("current user set to: %v\n", s.cfg.CurrentUserName)
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("no arguments required")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users: %v", err)
	}
	for _, user := range users {
		line := fmt.Sprintf("* %v", user.Name)
		if user.Name == s.cfg.CurrentUserName {
			line += " (current)"
		}
		fmt.Println(line)
	}
	return nil
}
