package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"pin-creator/accessToken"
	"pin-creator/config"
	"pin-creator/pinterest"
	"pin-creator/schedule"

	"pin-creator/internal/logger"
)

var cfg *config.Config

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Set log level to Info (0) or Debug (-1) here
	ctx = logger.NewContext(ctx)

	myLogger := logger.NewLogger(logger.LoggerConfig{UseJSON: false, LogLevel: 1})
	ctx = logger.WithLogger(ctx, myLogger)
	readConfig(ctx)

	log := logger.FromContext(ctx)
	log.Info("Doing something", "step", 1)
	log.V(1).Info("Debug info", "details", "some debug details")

	log.Info("Checking for pins to create in", cfg.ScheduleFilePath)

	scheduleReader := schedule.NewScheduleReader(cfg.ScheduleFilePath)
	nextPinData, err := scheduleReader.Next()
	if err != nil {
		log.Error(err, "error reading next schedule")
		os.Exit(1)
	}
	if nextPinData == nil {
		log.Info("No pin scheduled for creation")
		return
	}

	// client := getClient()
	// err = client.DeleteBoards("testboard\\d+")
	// if err != nil {
	// 	log.Fatalf("Error deleting boards: %v", err)
	// }

	start := time.Now()
	err = createPin(ctx, nextPinData)
	duration := time.Since(start)

	if err != nil {
		log.Error(err, "error creating pin")
		os.Exit(1)
	}

	log.Info(fmt.Sprintf("Pin creation took %s", duration.Truncate(time.Second)))

	err = scheduleReader.SetCreated(nextPinData.Index)
	if err != nil {
		log.Error(err, "error setting pin created to true")
		os.Exit(1)
	}
}

func readConfig(ctx context.Context) {
	log := logger.FromContext(ctx)
	args := os.Args
	if len(args) != 2 {
		log.Error(nil, "config.yaml file not provided")
		os.Exit(1)
	}

	configFilePath := args[1]

	cr := config.NewReader(configFilePath)
	c, err := cr.Read()
	if err != nil {
		log.Error(err, fmt.Sprintf("error reading %s", configFilePath))
		os.Exit(1)
	}
	cfg = c
}

func getToken(ctx context.Context) string {
	fmt.Println(cfg)
	tokenFileHandler := accessToken.NewAccessTokenFileHandler(cfg.AccessTokenPath)

	log := logger.FromContext(ctx)

	log.Info("Reading access token from file")
	token, err := tokenFileHandler.Read()
	if err == nil {
		return token
	} else {
		log.Info("No access token file found. Creating new token")

		tokenCreator := accessToken.NewAccessAccessTokenCreator(cfg.BrowserPath, cfg.RedirectPort)
		appId := os.Getenv("APP_ID")
		appSecret := os.Getenv("APP_SECRET")

		token, err := tokenCreator.NewToken(appId, appSecret)
		if err != nil {
			log.Error(err, "error creating new access token")
			os.Exit(1)
		}

		log.Info("Writing access token to file")
		err = tokenFileHandler.Write(token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing to token file handler: %v", err)
		}

		return token
	}
}

func getClient(ctx context.Context) pinterest.ClientInterface {
	token := getToken(ctx)
	return pinterest.NewClient(token)
}

func createPin(ctx context.Context, scheduledPinData *schedule.NextPinData) error {
	log := logger.FromContext(ctx)

	client := getClient(ctx)

	var boardId string
	var err error

	err = retry(ctx, func() error {
		boards, err := client.ListBoards(ctx)
		if err != nil {
			return err
		}

		boardId, err = boardIdByName(boards, scheduledPinData.BoardName)
		if err == nil {
			return nil
		}

		log.V(1).Info("Board not found. Creating new board.")
		err = client.CreateBoard(ctx, pinterest.BoardData{
			Name:        scheduledPinData.BoardName,
			Description: "Created by pin-creator",
			Privacy:     "PUBLIC",
		})
		if err != nil {
			return err
		}

		return fmt.Errorf("board not found after creation, retrying: %v", err)
	})
	if err != nil {
		return fmt.Errorf("failed to create or find board: %v", err)
	}

	pinData := pinterest.PinData{
		BoardId:     boardId,
		ImgPath:     scheduledPinData.ImagePath,
		Link:        scheduledPinData.Link,
		Title:       scheduledPinData.Title,
		Description: scheduledPinData.Description,
		AltText:     scheduledPinData.Description,
	}

	err = client.CreatePin(ctx, pinData)
	if err != nil {
		return fmt.Errorf("failed to create pin: %w", err)
	}

	log.Info(fmt.Sprintf("Created Pin '%s' in board '%s'", pinData.Title, scheduledPinData.BoardName))
	return nil
}

func boardIdByName(boards []pinterest.BoardInfo, boardName string) (string, error) {
	for _, board := range boards {
		if board.Name == boardName {
			return board.Id, nil
		}
	}

	return "", fmt.Errorf("board %s not found", boardName)
}

func retry(ctx context.Context, f func() error) error {
	backoff := time.Second
	for {
		err := f()
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("operation failed after retries: %w", err)
		case <-time.After(backoff):
			backoff *= 2
			if backoff > 60*time.Second {
				backoff = 60 * time.Second
			}
		}
	}
}
