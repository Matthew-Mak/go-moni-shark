package storage

import (
	"bufio"
	"fmt"
	"github.com/Matthew-Mak/go-moni-shark/errors"
	"github.com/Matthew-Mak/go-moni-shark/pkg/images"
	"log"
	"os"
	"strings"
)

func IsValidDiscordAttachmentURL(url string) bool {
	const (
		mediaPrefix = "https://media.discordapp.net/attachments/"
		cdnPrefix   = "https://cdn.discordapp.com/attachments/"
	)
	return strings.HasPrefix(url, mediaPrefix) || strings.HasPrefix(url, cdnPrefix)
}

func AddImage(image images.Image, path string) error {
	file, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return fmt.Errorf("%v, %w", path, errors.ErrCreateFile)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	writer := bufio.NewWriter(file)
	line := fmt.Sprintf("%v\n", image.Link)
	_, err = writer.WriteString(line)
	if err != nil {
		return fmt.Errorf("%v, %w", path, errors.ErrWriteLine)
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("%v, %w", path, errors.ErrInvolvingWriter)
	}
	log.Printf("Added image: %v", image.Link)
	return nil
}

func LoadImages(path string) ([]images.Image, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []images.Image{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("%v: %w", path, errors.ErrOpenFile)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", path, errors.ErrOpenFile)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	fileImages := make([]images.Image, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		link := strings.TrimSpace(line)

		image := images.Image{
			Link: link,
		}
		fileImages = append(fileImages, image)
	}
	return fileImages, nil
}
