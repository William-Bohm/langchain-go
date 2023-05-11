package tools

import "fmt"

func GetColoredText(text string, color string) string {
	colorStr := textColorMapping[color]
	return fmt.Sprintf("\033[%sm\033[1;3m%s\033[0m", colorStr, text)
}

var textColorMapping = map[string]string{
	"black":   "30",
	"red":     "31",
	"green":   "32",
	"yellow":  "33",
	"blue":    "34",
	"magenta": "35",
	"cyan":    "36",
	"white":   "37",
}

func PrintText(text string, color *string, end string) {
	if color == nil {
		fmt.Print(text)
	} else {
		textToPrint := GetColoredText(text, *color)
		fmt.Print(textToPrint)
	}
	fmt.Print(end)
}

func GetColorMapping(items []string, excludedColors ...string) map[string]string {
	colors := make([]string, len(textColorMapping))
	i := 0
	for k := range textColorMapping {
		colors[i] = k
		i++
	}

	if excludedColors != nil {
		tempColors := colors[:0]
		for _, c := range colors {
			excluded := false
			for _, ec := range excludedColors {
				if c == ec {
					excluded = true
					break
				}
			}
			if !excluded {
				tempColors = append(tempColors, c)
			}
		}
		colors = tempColors
	}

	colorMapping := make(map[string]string)
	for i, item := range items {
		colorMapping[item] = colors[i%len(colors)]
	}

	return colorMapping
}
