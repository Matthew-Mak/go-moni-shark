package messageHandler

import (
	"github.com/Matthew-Mak/go-moni-shark/config"
	"github.com/Matthew-Mak/go-moni-shark/pkg/images"
	"github.com/Matthew-Mak/go-moni-shark/pkg/storage"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"time"
)

var (
	BotId         string
	SliceOfImages []images.Image
	SliceOfAkula  []images.Image
	ImagesPath    string = "../storage/images.txt"
	AkulaPath     string = "../storage/akula.txt"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	// if m.content contains botid (Mentions) and "ping" then send "pong!"
	if m.Content == "<@"+BotId+"> ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong me daddy~")
	}

	if m.Content == config.BotPrefix+"ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong!")
	}
}

func Commands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		log.Println("CommandPing: Not yet implemented")
		return
	}

	switch i.ApplicationCommandData().Name {
	case "ping":
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "pong me daddy~",
			},
		})
		if err != nil {
			log.Printf("Failed to respond to /ping: %v", err)
		}
		log.Println("Pong me daddy~")
	case "media":
		src := rand.NewSource(time.Now().UnixNano())
		randomInt := rand.New(src)
		randomIndex := randomInt.Intn(len(SliceOfImages))

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: SliceOfImages[randomIndex].Link,
			},
		})
		if err != nil {
			log.Printf("Failed to respond to /media: %v", err)
		}
		log.Println("Media is up!" + SliceOfImages[randomIndex].Link)
	case "akula":
		src := rand.NewSource(time.Now().UnixNano())
		randomInt := rand.New(src)
		randomIndex := randomInt.Intn(len(SliceOfAkula))

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: SliceOfAkula[randomIndex].Link,
			},
		})
		if err != nil {
			log.Printf("Failed to respond to /media: %v", err)
		}
		log.Println("Media is up!" + SliceOfAkula[randomIndex].Link)
	case "add_media":
		var link string
		for _, opt := range i.ApplicationCommandData().Options {
			if opt.Name == "link" {
				link = opt.StringValue()
				break
			}
		}

		// Валидация на пустой аргумент
		if link == "" {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Аргумент оказался пустым",
				},
			})
			if err != nil {
				log.Printf("Failed to respond to /add_media: %v", err)
			}
		} else if !storage.IsValidDiscordAttachmentURL(link) { // Валидация на префикс ссылки
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Неверный формат ссылки",
				},
			})
			if err != nil {
				log.Printf("Failed to respond to /add_media: %v", err)
			}
		} else { // Добавление ссылки
			image := images.Image{Link: link}
			err := storage.AddImage(image, ImagesPath)
			if err != nil {
				log.Printf("Error adding image: %v", err)
				return
			}
			SliceOfImages = append(SliceOfImages, image)

			// Ответ об успешном добавлении ссылки
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ссылка была успешно добавлена",
				},
			})
			if err != nil {
				log.Printf("Failed to respond to /add_media: %v", err)
			}
		}
	}
}
