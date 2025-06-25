package bot

import (
	"github.com/Matthew-Mak/go-moni-shark/config"
	"github.com/Matthew-Mak/go-moni-shark/handlers/messageHandler"
	"github.com/Matthew-Mak/go-moni-shark/pkg/storage"
	"github.com/bwmarrin/discordgo"
	"log"
)

var (
	goBot *discordgo.Session
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

	// создаём команды
	for _, cmd := range commands {
		_, err := goBot.ApplicationCommandCreate(
			//config.GuildID, // или "" для глобальной команды (будет доступна везде, но дольше распространяется)
			config.AppId,
			"1352324363952717856", // пустая строка = глобально; или укажите ID гильдии
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
	messageHandler.BotId = user.ID

	// добавляем хэндлеры, они отвечают за ответ на запросы
	goBot.AddHandler(messageHandler.Ping)
	goBot.AddHandler(messageHandler.Commands)

	err = goBot.Open()

	if err != nil {
		log.Println("Error opening Discord session,", err)
		return
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")
}
