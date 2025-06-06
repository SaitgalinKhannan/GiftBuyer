package utils

import "strings"

func Contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// ChannelsToString преобразует []string в строку с запятыми
func ChannelsToString(channels []string) string {
	return strings.Join(channels, ",")
}

// StringToChannels преобразует строку в []string
func StringToChannels(s *string) []string {
	if s == nil || *s == "" {
		return nil
	}
	return strings.Split(*s, ",")
}
