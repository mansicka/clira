package event

type FileSaveEvent func(message string)

var onFileSaved FileSaveEvent

func SetFileSaveEventListener(listener FileSaveEvent) {
	onFileSaved = listener
}

func TriggerFileSaveEvent(message string) {
	if onFileSaved != nil {
		onFileSaved(message)
	}
}
