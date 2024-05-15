package main

import "fmt"

func HelpString() string {
	s := "dbpull-cli list command:\nstart    (to start the app)\nversion  (check version app)\nhelp     (show list command)"

	resume := fmt.Sprintf("%s\nyou need a help?\n%s", skyText.Render(asciiLogo), s)

	return resume
}
