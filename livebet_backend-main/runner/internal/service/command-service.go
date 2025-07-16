package service

import (
	"context"
	"livebets/runner/cmd/config"
	"livebets/runner/internal/entity"
	"livebets/runner/internal/storage"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/gofor-little/env"
	"github.com/rs/zerolog"
)

type CommandService struct {
	cfg              config.CommandConfig
	bookmakerStorage *storage.BookmakerStorage
	receiveChan      chan entity.Command
	logger           *zerolog.Logger
}

func NewCommandService(
	cfg config.CommandConfig,
	bookmakerStorage *storage.BookmakerStorage,
	logger *zerolog.Logger,
) *CommandService {
	receiveChan := make(chan entity.Command)
	return &CommandService{
		cfg:              cfg,
		bookmakerStorage: bookmakerStorage,
		receiveChan:      receiveChan,
		logger:           logger,
	}
}

func (c *CommandService) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	c.InitRun()

	for {
		select {
		case command := <-c.receiveChan:

			bookmakers := c.bookmakerStorage.ReadAll()

			bookmaker, ok := bookmakers[command.Name]
			if !ok {
				c.logger.Info().Msgf("[CommandService.Run] unknown command - %v", command)
				continue
			} else {
				c.logger.Info().Msgf("[CommandService.Run] register command - %v", command)
			}

			if err := c.writeEnvFile(bookmaker, command); err != nil {
				c.logger.Error().Err(err).Msg("[CommandService.Run] write to env error")
				continue
			}

			c.bookmakerStorage.SetReplicas(command)

			if err := c.runCommand(); err != nil {
				c.logger.Error().Err(err).Msg("[CommandService.Run] run docker compose error")
				continue
			}

		case <-ctx.Done():
			close(c.receiveChan)
			return
		}
	}
}

func (c *CommandService) SetCommand(ctx context.Context, command entity.Command) {
	c.receiveChan <- command
	return
}

func (c *CommandService) InitRun() error {
	var err error
	bookmakers := c.bookmakerStorage.ReadAll()
	for _, val := range bookmakers {

		if val.Replicas != 0 {

			err = c.writeEnvFile(val, entity.Command{Name: val.ReplicasName, Run: false})

			c.bookmakerStorage.SetReplicas(entity.Command{Name: val.Name, Run: false})

		}
	}
	if err != nil {
		c.logger.Error().Err(err).Msg("[CommandService.InitRun] write to env error")
		return err
	}

	if err = c.runCommand(); err != nil {
		c.logger.Error().Err(err).Msgf("[CommandService.InitRun] run docker compose error: %s", err)
		return err
	}

	return nil
}

func (c *CommandService) writeEnvFile(bookmaker entity.Bookmaker, command entity.Command) error {
	replicas := 0
	if command.Run {
		replicas = 1
	}

	return env.Write(bookmaker.ReplicasName, strconv.Itoa(replicas), c.cfg.EnvPath, true)
}

func (c *CommandService) runCommand() error {
	cmd := exec.Command("docker", "compose", "up", "-d", "--build")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
