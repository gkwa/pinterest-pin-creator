package pinterest

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
)

const (
	defaultWaitTime = 2 * time.Second
)

type ErrBoardNotFound struct {
	BoardName string
}

func (e ErrBoardNotFound) Error() string {
	return fmt.Sprintf("board %s not found", e.BoardName)
}

func findOrCreateBoard(ctx context.Context, client ClientInterface, log logr.Logger, boardName string) (string, error) {
	boardID, err := findBoard(ctx, client, log, boardName)
	if err == nil {
		return boardID, nil
	}

	if _, ok := err.(ErrBoardNotFound); !ok {
		return "", fmt.Errorf("error finding board: %w", err)
	}

	log.V(1).Info("Board not found. Creating new board.", "boardName", boardName)
	err = client.CreateBoard(ctx, BoardData{
		Name:        boardName,
		Description: "Created by pin-creator",
		Privacy:     "PUBLIC",
	})
	if err != nil {
		return "", fmt.Errorf("error creating board: %w", err)
	}

	log.V(2).Info("Board creation request sent successfully")
	time.Sleep(defaultWaitTime)

	return findBoard(ctx, client, log, boardName)
}
