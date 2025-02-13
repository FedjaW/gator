package main

import (
	"context"
	"errors"
	"fmt"
	"log"
)

func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) == 0 {
        return errors.New("no argument given for login")
    }

    username := cmd.args[0]

    _, err := s.db.GetUser(context.Background(), username)
    if err != nil {
        log.Fatalf("user %s does not exist", username)
    }

    err = s.cfg.SetUser(username)
    if err != nil {
        return err
    }
    
    fmt.Printf("user has been set to: %s\n", username)

    return nil
}
