package bot

import (
	"github.com/Matthew-Mak/go-moni-shark/config"
	"github.com/Matthew-Mak/go-moni-shark/handlers/messageHandler"
	"github.com/Matthew-Mak/go-moni-shark/pkg/storage"
	"github.com/bwmarrin/discordgo"
	"log"
)

var (
	goBot   *discordgo.Session
	guildID string = "1352324363952717856"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Проверить отклик бота",
	},
	{
		Name:        "media",
		Description: "Высылает медийку по вышу душу",
	},
	{
		Name:        "akula",
		Description: "Высылает Акулу по вышу душу",
	},
	{
		Name:        "add_media",
		Description: "Добавить медиа файл по ссылке",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "link",
				Description: "URL медиа файл (начиная с https://media.discordapp.net/... или https://cdn.discordapp.com/... )",
				Required:    true,
			},
		},
	},
}

func cleanupCommands(s *discordgo.Session, appID string) {
	commands, err := s.ApplicationCommands(appID, guildID)
	if err != nil {
		log.Println("Error getting commands:", err)
		return
	}

	for _, cmd := range commands {
		err = s.ApplicationCommandDelete(appID, guildID, cmd.ID)
		if err != nil {
			log.Println("Error deleting command:", cmd.ID, "Err: ", err)
			return
		} else {
			log.Println("Command deleted:", cmd.ID)
		}
	}
}

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Println("Error creating Discord session: ", err)
		return
	}

	user, err := goBot.User("@me")

	if err != nil {
		log.Println("Error getting user,", err)
		return
	}

	// чистим команды на случай если остались вырезанные
	cleanupCommands(goBot, config.AppId)

	// создаём команды
	for _, cmd := range commands {
		_, err := goBot.ApplicationCommandCreate(
			//config.GuildID, // или "" для глобальной команды (будет доступна везде, но дольше распространяется)
			config.AppId,
			guildID, // пустая строка = глобально; или укажите ID гильдии
			cmd,
		)
		if err != nil {
			log.Fatalf("Cannot create slash command %s: %v", cmd.Name, err)
		}
	}

	// вскрытие файла с картинками и перевод его в слайс
	messageHandler.SliceOfImages, err = storage.LoadImages(messageHandler.ImagesPath)
	if err != nil {
		log.Fatalf("Error loading images: %v", err)
		return
	}
	messageHandler.SliceOfAkula, err = storage.LoadImages(messageHandler.AkulaPath)
	if err != nil {
		log.Fatalf("Error loading Akulas!: %v", err)
		return
	}
	messageHandler.BotId = user.ID

	// добавляем хэндлеры, они отвечают за ответ на запросы
	goBot.AddHandler(messageHandler.Ping)
	goBot.AddHandler(messageHandler.Commands)

	err = goBot.Open()
	if err != nil {
		log.Println("Error opening Discord session,", err)
		return
	}

	err = goBot.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: "BYOND.com",
				Type: discordgo.ActivityTypeWatching,
			},
		},
		Status: "dnd",
	})
	if err != nil {
		log.Println("Error updating Discord status,", err)
		return
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")
}
