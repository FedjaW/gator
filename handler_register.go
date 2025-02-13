package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fedjaw/gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
    if len(cmd.args) == 0 {
        return errors.New("name was not provided")
    }

    username := cmd.args[0]
    createUerParams := database.CreateUserParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Name: username,
    }

    user, err := s.db.CreateUser(context.Background(), createUerParams)
    if err != nil {
        log.Fatalf("%v", err)
    }

    s.cfg.SetUser(username)
    fmt.Printf("user: %+v was successfully registered", user)

    return nil
}
