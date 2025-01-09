package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Event struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var producer Producer

type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

func NewProducer(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666) //nolint:all // fix autotests
	if err != nil {
		return fmt.Errorf("failed to OpenFile: %w", err)
	}

	producer = Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}
	return nil
}

func GetProducer() *Producer {
	return &producer
}

func (p *Producer) WriteEvent(event *Event) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal: %w", err)
	}

	if _, err := p.writer.Write(data); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	if err := p.writer.Flush(); err != nil {
		return fmt.Errorf("wrap: %w", err)
	}

	return nil
}

var consumer Consumer

type Consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewConsumer(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666) //nolint:all // fix autotests
	if err != nil {
		return fmt.Errorf("failed to OpenFile: %w", err)
	}

	consumer = Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}
	return nil
}

func GetConsumer() *Consumer {
	return &consumer
}

func (c *Consumer) GetEvents() ([]Event, error) {
	events := make([]Event, 0)
	c.scanner.Split(bufio.ScanLines)
	for c.scanner.Scan() {
		data := c.scanner.Bytes()
		event := Event{}
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, fmt.Errorf("failed to json.Unmarshal: %w", err)
		}
		events = append(events, event)
	}
	err := c.file.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to file.Close: %w", err)
	}
	return events, nil
}

func (c *Consumer) Close() error {
	err := c.file.Close()
	if err != nil {
		return fmt.Errorf("wrap: %w", err)
	}
	return nil
}
