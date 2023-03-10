package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/shmn7iii/gigi/internal/bosyu"
	"github.com/shmn7iii/gigi/internal/eula"
	"github.com/shmn7iii/gigi/internal/hey"
	"github.com/shmn7iii/gigi/internal/osiire"

	"github.com/bwmarrin/discordgo"
	"github.com/imdario/mergo"
)

var s *discordgo.Session

func composeCommands() (commands []*discordgo.ApplicationCommand) {
	commands = append(commands, bosyu.GetCommandsArray()...)
	commands = append(commands, osiire.GetCommandsArray()...)
	commands = append(commands, eula.GetCommandsArray()...)
	commands = append(commands, hey.GetCommandsArray()...)

	return
}

func composeCommandHandlers() (commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	mergo.Merge(&commandHandlers, bosyu.GetCommandHandlersMap())
	mergo.Merge(&commandHandlers, osiire.GetCommandHandlersMap())
	mergo.Merge(&commandHandlers, eula.GetCommandHandlersMap())
	mergo.Merge(&commandHandlers, hey.GetCommandHandlersMap())

	return
}

func composeComponentHandlers() (componentsHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	mergo.Merge(&componentsHandlers, bosyu.GetComponentHandlersMap())
	mergo.Merge(&componentsHandlers, eula.GetComponentHandlersMap())

	return
}

func init() {
	var err error
	s, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	var commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) = composeCommandHandlers()
	var componentsHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) = composeComponentHandlers()

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	s.AddHandler(hey.MessageCreate)

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("[gigi] ohayo.")
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	commands := composeCommands()
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop

	for _, cmd := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
		if err != nil {
			log.Fatalf("Cannot delete %q command: %v", cmd.Name, err)
		}
	}

	log.Println("[gigi] byebye.")
}
