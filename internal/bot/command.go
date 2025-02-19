package bot

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/sirupsen/logrus"
)

type command struct {
	ID     int
	isBase bool
	rgx    *regexp.Regexp
	action func(input []string, batch any, srvc *service.Service, log *logrus.Logger, sender *MessageSender, cl *Client, cmd *command)

	child int
}

var (
	internalErrorAditionalInfo = fmt.Sprintf("Please, contact @%s to share this interesting case", os.Getenv("TELEGRAM_USERNAME"))

	commandsToIDs = map[string]int{
		"add category":             1,
		"add category description": 2,
		"add record":               3,
		"show categories":          4,
		"show records":             5,
		"get time boundaries":      6,
	}

	commandReplies = map[int]string{
		1: "Please, input category name",
		3: "Please, input category name and amount e.g. 'category 100.5'\nOptionally you can add description e.g. 'category 100.5 description'",
		4: "Please, input the number of categories you want to see:\n\n" +
			" - 'n' for n number of categories\n" +
			" - 'all' for all categories\n" +
			" - 'category name' for one specific category\n\n" +
			"Optionally you can add 'full' to see descriptions as well",
		5: "Please, input the category name",
	}

	commandsByIDs = map[int]command{
		1: {
			ID:     1,
			isBase: true,
			rgx:    regexp.MustCompile(`^([a-zA-Z0-9]{1,10})$`), //TODO: update
			action: addCategoryAction,
			child:  2,
		},
		2: {
			ID:     2,
			isBase: false,
			rgx:    regexp.MustCompile(`^([a-zA-Z0-9 ]+)$`), //TODO: update
			action: addCategoryDescriptionAction,
			child:  0,
		},
		3: {
			ID:     3,
			isBase: true,
			rgx:    regexp.MustCompile(`^\s*(?P<category>[a-zA-Z0-9]{1,10})\s*(?P<amount>\d+(?:\.\d{1,2})?)(?<description>\s+[a-zA-Z0-9 ]+)?$`),
			action: addRecordAction,
			child:  0,
		},
		4: {
			ID:     4,
			isBase: true,
			rgx:    regexp.MustCompile(`^(?:(?P<number>\d+)|(?P<category>[a-zA-Z0-9]{1,10}))\s*(?P<isfull>full)?$`),
			action: showCategoriesAction,
			child:  0,
		},
		5: {
			ID:     5,
			isBase: true,
			rgx:    regexp.MustCompile(`^(?P<category>[a-zA-Z0-9]{1,10})$`),
			action: showRecordsAction,
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
			action: getTimeBoundariesAction,
			child:  0,
		},
	}
)

func addCategoryAction(input []string, batch any, _ *service.Service, log *logrus.Logger, sender *MessageSender, cl *Client, cmd *command) {
	if len(input) != 2 {
		log.Debug("wrong input for add category command")
		return
	}

	batch.(*ftracker.SpendingCategory).Category = input[1]
	log.Debug("action on add category command")
	sender.Send(
		tgbotapi.NewMessage(cl.chanID, "Please, type description to a new category"),
	)
}

func addCategoryDescriptionAction(input []string, batch any, srvc *service.Service, log *logrus.Logger, sender *MessageSender, cl *Client, cmd *command) {

	if len(input) != 2 {
		log.Debug("wrong input for add category command")
		return
	}

	log.Debug("action on add category description command")

	batch.(*ftracker.SpendingCategory).Description = input[1]

	msg := tgbotapi.NewMessage(cl.chanID, "")
	msg.ReplyMarkup = baseKeyboard
	defer func() {
		sender.Send(msg)
	}()

	err := cl.populateUserGUID(srvc, log)
	if err != nil {
		log.WithError(err).Error("error on fill user guid")
		msg.Text = "Sorry, something went wrong with the users database :(\n" + internalErrorAditionalInfo
		return
	}

	categoryToAdd := *batch.(*ftracker.SpendingCategory)
	categoryToAdd.UserGUID = cl.userGUID
	_, err = srvc.AddCategories([]ftracker.SpendingCategory{categoryToAdd})
	if err != nil {

		if utils.IsUniqueConstrainViolation(err) {
			msg.Text = "Category with that name already exists"
			return
		}

		log.WithError(err).Error("error on add category")
		msg.Text = "Sorry, something went wrong with the database adding the category:(\n" + internalErrorAditionalInfo
	} else {
		msg.Text = "Category added successfully"
	}
}

func addRecordAction(input []string, batch any, srvc *service.Service, log *logrus.Logger, sender *MessageSender, cl *Client, cmd *command) {
	if len(input) != 4 {
		log.Debug("wrong tocken number for add record command")
		return
	}
	recordCategory := input[1:2]
	recordAmount := input[2]
	if splitedAmount := strings.Split(recordAmount, "."); len(splitedAmount) == 1 {
		recordAmount += "00"
	} else if len(splitedAmount[1]) == 1 {
		recordAmount = strings.Join(splitedAmount, "") + "0"
	} else {
		recordAmount = strings.Join(splitedAmount, "")
	}

	recordDescription := input[3]
	if len(recordDescription) == 0 {
		recordDescription = "spending"
	}

	msg := tgbotapi.NewMessage(cl.chanID, "")
	msg.ReplyMarkup = baseKeyboard
	defer func() {
		sender.Send(msg)
	}()

	log.Debug("category to lookup: ", recordCategory)

	categories, err := srvc.GetCategories(srvc.SpendingCategory.WithCategories(recordCategory))
	if err != nil {
		log.WithError(err).Error("error on get category")
		msg.Text = "Sorry, something went wrong with the database getting the category :(\n" + internalErrorAditionalInfo
		return
	}

	if len(categories) == 0 {
		msg.Text = "There is no such category"
		return
	}

	batch.(*ftracker.SpendingRecord).CategoryGUID = categories[0].GUID
	amount, err := strconv.ParseUint(recordAmount, 10, 16)
	if err != nil {
		log.WithError(err).Error("error on parsing amount")
		msg.Text = "Wow, there is something wrong with the amount you've entered\n" + internalErrorAditionalInfo
		return
	}
	batch.(*ftracker.SpendingRecord).Amount = uint16(amount)
	batch.(*ftracker.SpendingRecord).Description = recordDescription

	recordToAdd := *batch.(*ftracker.SpendingRecord)
	_, err = srvc.AddRecords([]ftracker.SpendingRecord{recordToAdd})
	if err != nil {
		log.WithError(err).Error("error on add record")
		msg.Text = "Sorry, something went wrong with the database adding the record :(\n" + internalErrorAditionalInfo
	} else {
		msg.Text = "Record added successfully"
	}
}

func showCategoriesAction(input []string, batch any, srvc *service.Service, log *logrus.Logger, sender *MessageSender, cl *Client, cmd *command) {
	if len(input) != 4 {
		log.Debug("wrong tocken number for add record command")
		return
	}

	msg := tgbotapi.NewMessage(cl.chanID, "")
	msg.ReplyMarkup = baseKeyboard
	defer func() {
		sender.Send(msg)
	}()

	var categoriesLimit int
	var categoryNames []string
	var err error

	switch instruction := input[2]; instruction {
	case "":
		categoriesLimit, err = strconv.Atoi(input[1])
		if err != nil {
			log.WithError(err).Error("error on parsing limit")
			msg.Text = "Ooopsie, there is something wrong with the limit you've entered\n" + internalErrorAditionalInfo
			return
		}
	case "all":
		categoriesLimit = 0
	default:
		categoriesLimit = 1
		categoryNames = []string{instruction}
	}
	addDescription := input[3] == "full"

	cl.populateUserGUID(srvc, log)
	categories, err := srvc.GetCategories(
		srvc.SpendingCategory.WithUserGUIDs([]uuid.UUID{cl.userGUID}),
		srvc.SpendingCategory.WithLimit(categoriesLimit),
		srvc.SpendingCategory.WithCategories(categoryNames),
		srvc.SpendingCategory.WithOrder(repository.LastModifiedOrder),
	)
	if err != nil {
		log.WithError(err).Error("error on get categories")
		msg.Text = "Sorry, something went wrong with the database getting the categories :(\n" + internalErrorAditionalInfo
		return
	}

	if len(categories) == 0 {
		if len(categoryNames) == 0 {
			msg.Text = "You don't have any categories yet"
		} else {
			msg.Text = "There is no such category"
		}
		return
	}

	msg.Text = "Your categories:\n"
	var format string
	if addDescription {
		format += "%d. %s - %f\n%s\n\n"
	} else {
		format += "%d. %s - %f\n"
	}
	if addDescription {
		for i, category := range categories {
			msg.Text += fmt.Sprintf("%d. %s - %d.%d\u20AC\n%s\n\n", i+1, category.Category, category.Amount/100, category.Amount%100, category.Description)
		}
	} else {
		for i, category := range categories {
			msg.Text += fmt.Sprintf("%d. %s - %d.%d\u20AC\n", i+1, category.Category, category.Amount/100, category.Amount%100)
		}
	}
}

func showRecordsAction(input []string, batch any, srvc *service.Service, log *logrus.Logger, sender *MessageSender, cl *Client, cmd *command) {

	if len(input) != 2 {
		log.Debug("wrong tocken number for show records command")
		return
	}
	recordCategory := input[1:2]

	msg := tgbotapi.NewMessage(cl.chanID, "")
	defer func() {
		sender.Send(msg)
	}()

	categories, err := srvc.GetCategories(srvc.SpendingCategory.WithCategories(recordCategory))
	if err != nil {
		log.WithError(err).Error("error on get category")
		msg.Text = "Sorry, something went wrong with the database getting the category :(\n" + internalErrorAditionalInfo
		return
	}

	if len(categories) == 0 {
		msg.Text = "There is no such category"
		msg.ReplyMarkup = baseKeyboard
		cmd.child = 0
		return
	}
	batch.(*repository.RecordOptions).CategoryGUIDs = []uuid.UUID{categories[0].GUID}

	sender.Send(
		tgbotapi.NewMessage(cl.chanID, "Please, type the time period for the records ('from' 'to') in dd.mm.yyyy format:\n'dd.mm.yyyy dd.mm.yyyy'"), //update
	)
}

func getTimeBoundariesAction(input []string, batch any, srvc *service.Service, log *logrus.Logger, sender *MessageSender, cl *Client, cmd *command) {
	if len(input) != 6 {
		log.Debug("wrong tocken number for set time boundaries command")
		return
	}

	msg := tgbotapi.NewMessage(cl.chanID, "")
	msg.ReplyMarkup = baseKeyboard
	defer func() {
		sender.Send(msg)
	}()

	addDescription := input[5] == "full"
	var recordsLimit int
	var err error
	if input[1] == "all" {
		recordsLimit = 0
	} else {
		recordsLimit, err = strconv.Atoi(input[1])
		if err != nil {
			log.WithError(err).Error("error on parsing limit")
			msg.Text = "Ooopsie, there is something wrong with the number of records you've entered\n" + internalErrorAditionalInfo
			return
		}
	}

	var timeFrom, timeTo time.Time
	if input[2] == "" {
		timeFrom, err = time.Parse("02.01.2006", input[3])
		if err != nil {
			log.WithError(err).Error("error on parsing time from")
			msg.Text = "Wow, there is something wrong with the 'from' date you've entered"
			return
		}

		if input[4] != "" {
			timeTo = time.Now()
		} else {
			timeTo, err = time.Parse("02.01.2006", input[4])
			if err != nil {
				log.WithError(err).Error("error on parsing time to")
				msg.Text = "Wow, there is something wrong with the 'to' date you've entered"
				return
			}
		}
	} else {
		switch input[2] {
		case "year":
			timeTo = time.Now()
			timeFrom = timeTo.AddDate(-1, 0, 0)
		case "month":
			timeTo = time.Now()
			timeFrom = timeTo.AddDate(0, -1, 0)
		case "day":
			timeTo = time.Now()
			timeFrom = timeTo.AddDate(0, 0, -1)
		default:
			log.Debug("invalid token for ymd time boundaries")
			msg.Text = "Wow, there is something wrong with the time period you've entered"
			return
		}
	}

	recordOption := *batch.(*repository.RecordOptions)
	records, err := srvc.GetRecords(
		srvc.SpendingRecord.WithCategoryGUIDs(recordOption.CategoryGUIDs),
		srvc.SpendingRecord.WithTimeFrame(timeFrom, timeTo),
		srvc.SpendingRecord.WithLimit(recordsLimit),
		//add ordered by
	)
	if err != nil {
		log.WithError(err).Error("error on get records")
		msg.Text = "Sorry, something went wrong with the database getting the records :(\n" + internalErrorAditionalInfo
		return
	}

	if len(records) == 0 {
		msg.Text = "There are no records for this category and time period"
		return
	}

	var subtotal uint32 = 0
	if addDescription {
		for _, record := range records {
			msg.Text += fmt.Sprintf("[%s] %d.%d\u20AC - %s\n", record.CreatedAt.Format("Monday, 02 Jan, 15:04"), record.Amount/100, record.Amount%100, record.Description) //mb updated?
			subtotal += uint32(record.Amount)
		}
	} else {
		for _, record := range records {
			msg.Text += fmt.Sprintf("[%s] %d.%d\u20AC\n", record.CreatedAt.Format("Monday, 02 Jan, 15:04"), record.Amount/100, record.Amount%100)
			subtotal += uint32(record.Amount)
		}
	}

	msg.Text = fmt.Sprintf("Subtotal: %d.%d\u20AC\n\n", subtotal/100, subtotal%100) + msg.Text
}

func (c *command) validateInput(input string) []string {
	matches := c.rgx.FindAllStringSubmatch(input, 1)
	if len(matches) == 0 {
		return nil
	}
	return matches[0]
}

func (c *command) isLast() bool {
	return c.child == 0
}

func (c *command) next() command {
	return commandsByIDs[c.child]
}

func isCommand(s string) (command, bool) {
	id, ok := commandsToIDs[s]
	if !ok {
		return command{}, false
	}
	return commandsByIDs[id], true
}

func initBatch(id int) any {
	switch id {
	case 1:
		return &ftracker.SpendingCategory{}
	case 3:
		return &ftracker.SpendingRecord{}
	case 5:
		return &repository.RecordOptions{}
	}

	return nil
}
