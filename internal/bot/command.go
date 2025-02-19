package bot

import "regexp"

type command struct {
	ID     int
	isBase bool
	rgx    *regexp.Regexp

	child int
}

var commandsToIDs = map[string]int{
	"add category":             1,
	"add category description": 2,
	"add record":               3,
	"show categories":          4,
	"show records":             5,
	"get time boundaries":      6,
}

var commandReplies = map[int]string{
	1: "Please, input category name",
	3: "Please, input category name and amount e.g. 'category 100.5'\nOptionally you can add description e.g. 'category 100.5 description'",
	4: "Please, input the number of categories you want to see:\n\n" +
		" - 'n' for n number of categories\n" +
		" - 'all' for all categories\n" +
		" - 'category name' for one specific category\n\n" +
		"Optionally you can add 'full' to see descriptions as well",
	5: "Please, input the category name",
}

var commandsByIDs = map[int]command{
	1: {
		ID:     1,
		isBase: true,
		rgx:    regexp.MustCompile(`^([a-zA-Z0-9]{1,10})$`), //TODO: update
		child:  2,
	},
	2: {
		ID:     2,
		isBase: false,
		rgx:    regexp.MustCompile(`^([a-zA-Z0-9 ]+)$`), //TODO: update
		child:  0,
	},
	3: {
		ID:     3,
		isBase: true,
		rgx:    regexp.MustCompile(`^\s*(?P<category>[a-zA-Z0-9]{1,10})\s*(?P<amount>\d+(?:\.\d+)?)\s*(?<description>[a-zA-Z0-9 ]+)?$`),
		child:  0,
	},
	4: {
		ID:     4,
		isBase: true,
		rgx:    regexp.MustCompile(`^(?:(?P<number>\d+)|(?P<category>[a-zA-Z0-9]{1,10}))\s*(?P<isfull>full)?$`),
		child:  0,
	},
	5: {
		ID:     5,
		isBase: true,
		rgx:    regexp.MustCompile(`^(?P<category>[a-zA-Z0-9]{1,10})$`),
		child:  6,
	},
	6: {
		ID:     6,
		isBase: false,
		rgx: regexp.MustCompile(
			`^(?P<number>(?:\d+)|(?:all))\s*` +
				`(?:(?:(?:last)?\s*(?P<ymd>(?:year)|(?:month)|(?:day)))|` +
				`(?:(?P<from>\d{2}.\d{2}.\d{4})\s*(?P<to>\d{2}.\d{2}.\d{4})?))` +
				`\s*(?P<full>full)?$`,
		),
		child: 0,
	},
}

func isCommand(s string) (command, bool) {
	id, ok := commandsToIDs[s]
	if !ok {
		return command{}, false
	}
	return commandsByIDs[id], true
}
