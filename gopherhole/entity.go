package gopherhole

import (
	"fmt"
	"strings"
)

type EntityType byte

const (
	TextFileEntity      EntityType = '0'
	DirectoryEntity     EntityType = '1'
	CsoNameServerEntity EntityType = '2'
	ErrorEntity         EntityType = '3'
	MacBinaryEntity     EntityType = '4'
	DosBinaryEntity     EntityType = '5'
	UnixUUEncodedEntity EntityType = '6'
	SearchServiceEntity EntityType = '7'
	TelnetEntity        EntityType = '8'
	BinaryEntity        EntityType = '9'
	CalendarEntity      EntityType = 'c'
	EventEntity         EntityType = 'e'
	GifEntity           EntityType = 'g'
	HtmlEntity          EntityType = 'h'
	InlineTextEntity    EntityType = 'i'
	SoundEntity         EntityType = 's'
	ImageEntity         EntityType = 'I'
	MIMEEntity          EntityType = 'M'
	TN3270Entity        EntityType = 'T'
)

var AllEntities = [...]EntityType{
	TextFileEntity,
	DirectoryEntity,
	CsoNameServerEntity,
	ErrorEntity,
	MacBinaryEntity,
	DosBinaryEntity,
	UnixUUEncodedEntity,
	SearchServiceEntity,
	TelnetEntity,
	BinaryEntity,
	CalendarEntity,
	EventEntity,
	GifEntity,
	HtmlEntity,
	InlineTextEntity,
	SoundEntity,
	ImageEntity,
	MIMEEntity,
	TN3270Entity,
}

var MimeTypeEntities = map[string]EntityType{
	"audio/":                      SoundEntity,
	"application/":                BinaryEntity,
	"application/binhex":          MacBinaryEntity,
	"application/mac-binhex40":    MacBinaryEntity,
	"application/mac-binhex":      MacBinaryEntity,
	"application/x-msdos-program": DosBinaryEntity,
	"image/":                      ImageEntity,
	"image/gif":                   GifEntity,
	"multipart/":                  MIMEEntity,
	"text/":                       TextFileEntity,
	"text/calendar":               CalendarEntity,
	"text/x-vcalendar":            CalendarEntity,
	"text/html":                   HtmlEntity,
	"video/":                      BinaryEntity,
}

func EntityForMimeType(mime string) (entity EntityType, ok bool) {
	splitMime := strings.Split(mime, "/")
	if len(splitMime) != 2 || strings.TrimSpace(splitMime[1]) == "" {
		return
	}

	entity, ok = MimeTypeEntities[mime]
	if !ok {
		group := fmt.Sprintf("%s/", splitMime[0])
		entity, ok = MimeTypeEntities[group]
	}

	return
}
