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
	baseCtx := context.Background()
	ctx := logger.NewContext(baseCtx)

	myLogger := logger.NewLogger(logger.LoggerConfig{UseJSON: false, LogLevel: 0})
	ctx = logger.WithLogger(ctx, myLogger)
	readConfig(ctx)

	log := logger.FromContext(ctx)
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

	// client := getClient(ctx)
	// err = client.DeleteBoards(ctx, "testboard\\d+")
	// if err != nil {
	// 	log.Error(err, "error deleting boards")
	// 	os.Exit(1)
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
	log := logger.FromContext(ctx)
	tokenFileHandler := accessToken.NewAccessTokenFileHandler(cfg.AccessTokenPath)

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

	boardCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	boardId, err := pinterest.CreateOrFindBoard(boardCtx, client, log, scheduledPinData.BoardName)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Error(err, "Timeout occurred while creating or finding board")
			return fmt.Errorf("timeout occurred while creating or finding board: %w", err)
		}
		return fmt.Errorf("failed to create or find board: %w", err)
	}

	pinData := pinterest.PinData{
		BoardId:     boardId,
		ImgPath:     scheduledPinData.ImagePath,
		Link:        scheduledPinData.Link,
		Title:       scheduledPinData.Title,
		Description: scheduledPinData.Description,
		AltText:     scheduledPinData.Description,
	}

	pinCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	err = client.CreatePin(pinCtx, pinData)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Error(err, "Timeout occurred while creating pin")
			return fmt.Errorf("timeout occurred while creating pin: %w", err)
		}
		return fmt.Errorf("failed to create pin: %w", err)
	}

	log.Info(fmt.Sprintf("Created Pin '%s' in board '%s'", pinData.Title, scheduledPinData.BoardName))
	return nil
}
