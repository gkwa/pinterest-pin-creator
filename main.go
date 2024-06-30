package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"pin-creator/accessToken"
	"pin-creator/config"
	"pin-creator/pinterest"
	"pin-creator/schedule"

	log "github.com/sirupsen/logrus"
)

var cfg *config.Config

func main() {
	readConfig()

	log.Infof("Checking for pins to create in '%s'", cfg.ScheduleFilePath)

	scheduleReader := schedule.NewScheduleReader(cfg.ScheduleFilePath)
	nextPinData, err := scheduleReader.Next()
	if err != nil {
		log.Fatal(err.Error())
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

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	start := time.Now()
	err = createPin(ctx, nextPinData)
	duration := time.Since(start)

	if err != nil {
		log.Fatalf("Error creating pin: %v", err)
	}

	log.Infof("Pin creation took %s", duration.Truncate(time.Second))

	err = scheduleReader.SetCreated(nextPinData.Index)
	if err != nil {
		log.Fatalf("Error setting pin created to true. Error: %s", err.Error())
	}
}

func readConfig() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("config.yaml file not provided")
	}

	configFilePath := args[1]

	cr := config.NewReader(configFilePath)
	c, err := cr.Read()
	if err != nil {
		log.Fatal(err.Error())
	}
	cfg = c
}

func getToken() string {
	fmt.Println(cfg)
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
			log.Fatalf("error creating new access token. Error: %s", err.Error())
		}

		log.Info("Writing access token to file")
		err = tokenFileHandler.Write(token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing to token file handler: %v", err)
		}

		return token
	}
}

func getClient() pinterest.ClientInterface {
	token := getToken()
	return pinterest.NewClient(token)
}

func createPin(ctx context.Context, scheduledPinData *schedule.NextPinData) error {
	client := getClient()

	var boardId string
	var err error

	err = retry(ctx, func() error {
		boards, err := client.ListBoards()
		if err != nil {
			return err
		}

		boardId, err = boardIdByName(boards, scheduledPinData.BoardName)
		if err == nil {
			return nil
		}

		log.Info("Board not found. Creating new board.")
		err = client.CreateBoard(pinterest.BoardData{
			Name:        scheduledPinData.BoardName,
			Description: "Created by pin-creator",
			Privacy:     "PUBLIC",
		})
		if err != nil {
			return err
		}

		return errors.New("board not found after creation, retrying")
	})
	if err != nil {
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

	err = client.CreatePin(pinData)
	if err != nil {
		return fmt.Errorf("failed to create pin: %w", err)
	}

	log.Infof("Created Pin '%s' in board '%s'\n", pinData.Title, scheduledPinData.BoardName)
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
