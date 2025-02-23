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
	action func(input []string, batch any, srvc service.ServiceInterface, log *logrus.Logger, sender Sender, cl *Client, cmd *command)

	child int
}

const (
	formatOut = "Monday, 02 Jan, 15:04"
	formatIn  = "02.01.2006"
)

const (
	CommandAddCategory            = "add category"
	CommandAddCategoryDescription = "add category description"
	CommandAddRecord              = "add record"
	CommandShowCategories         = "show categories"
	CommandShowRecords            = "show records"
	CommandGetTimeBoundaries      = "get time boundaries"
)

var (
	internalErrorAditionalInfo = fmt.Sprintf("Please, contact @%s to share this interesting case", os.Getenv("TELEGRAM_USERNAME"))

	commandsToIDs = map[string]int{
		CommandAddCategory:            1,
		CommandAddCategoryDescription: 2,
		CommandAddRecord:              3,
		CommandShowCategories:         4,
		CommandShowRecords:            5,
		CommandGetTimeBoundaries:      6,
	}

	commandReplies = map[int]string{
		1: MessageAddCategory,
		3: MessageAddRecord,
		4: MessageShowCategories,
		5: MessageShowRecords,
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

func addCategoryAction(input []string, batch any, _ service.ServiceInterface, log *logrus.Logger, sender Sender, cl *Client, cmd *command) {
	if len(input) != 2 {
		log.Debug("wrong input for add category command")
		cmd.becomeLast()
		msg := tgbotapi.NewMessage(cl.chanID, MessageInvalidNumberOfTockensAction+"\n"+internalErrorAditionalInfo)
		msg.ReplyMarkup = baseKeyboard
		sender.Send(msg)
		return
	}

	batch.(*ftracker.SpendingCategory).Category = input[1]
	log.Debug("action on add category command")
	sender.Send(
		tgbotapi.NewMessage(cl.chanID, MessageAddCategoryDescription),
	)
}

func addCategoryDescriptionAction(input []string, batch any, srvc service.ServiceInterface, log *logrus.Logger, sender Sender, cl *Client, cmd *command) {

	if len(input) != 2 {
		log.Debug("wrong input for add category command")
		msg := tgbotapi.NewMessage(cl.chanID, MessageInvalidNumberOfTockensAction+"\n"+internalErrorAditionalInfo)
		msg.ReplyMarkup = baseKeyboard
		sender.Send(msg)
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
		msg.Text = MessageDatabaseError + "\n" + internalErrorAditionalInfo
		return
	}

	categoryToAdd := *batch.(*ftracker.SpendingCategory)
	categoryToAdd.UserGUID = cl.userGUID
	_, err = srvc.AddCategories([]ftracker.SpendingCategory{categoryToAdd})
	if err != nil {

		if utils.IsUniqueConstrainViolation(err) {
			msg.Text = MessageCategoryDuplicate
			return
		}

		log.WithError(err).Error("error on add category")
		msg.Text = MessageDatabaseError + "\n" + internalErrorAditionalInfo
	} else {
		msg.Text = MessageCategorySuccess
	}
}

func addRecordAction(input []string, batch any, srvc service.ServiceInterface, log *logrus.Logger, sender Sender, cl *Client, cmd *command) {
	if len(input) != 4 {
		log.Debug("wrong tocken number for add record command")
		return
	}
	recordCategory := input[1:2]
	recordAmountLeft, recordAmountRight := splitAmount(input[2])

	recordDescription := input[3]
	if len(recordDescription) == 0 {
		recordDescription = "spending"
	}

	msg := tgbotapi.NewMessage(cl.chanID, "")
	msg.ReplyMarkup = baseKeyboard
	defer func() {
		sender.Send(msg)
	}()

	if recordAmountLeft == "0" && recordAmountRight == "00" {
		msg.Text = MessageZeroAmount
		return
	}

	log.Debug("category to lookup: ", recordCategory)

	categories, err := srvc.GetCategories(srvc.SpendingCategoriesWithCategories(recordCategory))
	if err != nil {
		log.WithError(err).Error("error on get category")
		msg.Text = MessageDatabaseError + "\n" + internalErrorAditionalInfo
		return
	}

	if len(categories) == 0 {
		msg.Text = "There is no such category"
		return
	}

	batch.(*ftracker.SpendingRecord).CategoryGUID = categories[0].GUID
	amount, err := strconv.ParseUint(recordAmountLeft+recordAmountRight, 10, 32)
	if err != nil {
		log.WithError(err).Error("error on parsing amount")
		msg.Text = MessageAmountError + "\n" + internalErrorAditionalInfo
		return
	}
	batch.(*ftracker.SpendingRecord).Amount = uint32(amount)
	batch.(*ftracker.SpendingRecord).Description = recordDescription

	recordToAdd := *batch.(*ftracker.SpendingRecord)
	_, err = srvc.AddRecords([]ftracker.SpendingRecord{recordToAdd})
	if err != nil {
		log.WithError(err).Error("error on add record")
		msg.Text = MessageDatabaseError + "\n" + internalErrorAditionalInfo
	} else {
		msg.Text = MessageRecordSuccess
	}
}

func showCategoriesAction(input []string, batch any, srvc service.ServiceInterface, log *logrus.Logger, sender Sender, cl *Client, cmd *command) {
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
			msg.Text = MessageLimitError + "\n" + internalErrorAditionalInfo
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
		srvc.SpendingCategoriesWithUserGUIDs([]uuid.UUID{cl.userGUID}),
		srvc.SpendingCategoriesWithLimit(categoriesLimit),
		srvc.SpendingCategoriesWithCategories(categoryNames),
		srvc.SpendingCategoriesWithOrder(service.OrderCategoriesByCreatedAt, false),
	)
	if err != nil {
		log.WithError(err).Error("error on get categories")
		msg.Text = MessageDatabaseError + "\n" + internalErrorAditionalInfo
		return
	}

	if len(categories) == 0 {
		if len(categoryNames) == 0 {
			msg.Text = MessageUnderflowCategories
		} else {
			msg.Text = MessageNoCategoryFound
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
			leftAmount, rightAmount := splitAmount(category.Amount)
			msg.Text += fmt.Sprintf("%d. %s - %s.%s\u20AC\n%s\n\n", i+1, category.Category, leftAmount, rightAmount, category.Description)
		}
	} else {
		for i, category := range categories {
			leftAmount, rightAmount := splitAmount(category.Amount)
			msg.Text += fmt.Sprintf("%d. %s - %s.%s\u20AC\n", i+1, category.Category, leftAmount, rightAmount)
		}
	}
}

func showRecordsAction(input []string, batch any, srvc service.ServiceInterface, log *logrus.Logger, sender Sender, cl *Client, cmd *command) {

	if len(input) != 2 {
		log.Debug("wrong tocken number for show records command")
		return
	}
	recordCategory := input[1:2]

	msg := tgbotapi.NewMessage(cl.chanID, "")
	defer func() {
		sender.Send(msg)
	}()

	categories, err := srvc.GetCategories(srvc.SpendingCategoriesWithCategories(recordCategory))
	if err != nil {
		log.WithError(err).Error("error on get category")
		msg.Text = MessageDatabaseError + "\n" + internalErrorAditionalInfo
		return
	}

	if len(categories) == 0 {
		msg.Text = MessageNoCategoryFound
		msg.ReplyMarkup = baseKeyboard
		cmd.becomeLast()
		return
	}
	batch.(*repository.RecordOptions).CategoryGUIDs = []uuid.UUID{categories[0].GUID}

	sender.Send(
		tgbotapi.NewMessage(cl.chanID, MessageAddTimeDetails), //update
	)
}

func getTimeBoundariesAction(input []string, batch any, srvc service.ServiceInterface, log *logrus.Logger, sender Sender, cl *Client, cmd *command) {
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
			msg.Text = MessageLimitError + "\n" + internalErrorAditionalInfo
			return
		}
	}

	var timeFrom, timeTo time.Time
	if input[2] == "" {
		timeFrom, err = time.Parse(formatIn, input[3])
		if err != nil {
			log.WithError(err).Error("error on parsing time from")
			msg.Text = MessageInvalidFromDate
			return
		}

		if input[4] != "" {
			timeTo = time.Now()
		} else {
			timeTo, err = time.Parse(formatIn, input[4])
			if err != nil {
				log.WithError(err).Error("error on parsing time to")
				msg.Text = MessageInvalidToDate
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
			msg.Text = MessageInvalidFixedTime
			return
		}
	}

	recordOption := *batch.(*repository.RecordOptions)
	records, err := srvc.GetRecords(
		srvc.SpendingRecordsWithCategoryGUIDs(recordOption.CategoryGUIDs),
		srvc.SpendingRecordsWithTimeFrame(timeFrom, timeTo),
		srvc.SpendingRecordsWithLimit(recordsLimit),
		srvc.SpendingRecordsWithOrder(service.OrderRecordsByCreatedAt, false),
	)
	if err != nil {
		log.WithError(err).Error("error on get records")
		msg.Text = MessageDatabaseError + "\n" + internalErrorAditionalInfo
		return
	}

	if len(records) == 0 {
		msg.Text = MessageUnderflowRecords
		return
	}

	var subtotal uint32 = 0
	if addDescription {
		for _, record := range records {
			leftAmount, rightAmount := splitAmount(record.Amount)
			msg.Text += fmt.Sprintf("[%s] %s.%s\u20AC - %s\n", record.CreatedAt.Format(formatOut), leftAmount, rightAmount, record.Description) //mb updated?
			subtotal += uint32(record.Amount)
		}
	} else {
		for _, record := range records {
			leftAmount, rightAmount := splitAmount(record.Amount)
			msg.Text += fmt.Sprintf("[%s] %s.%s\u20AC\n", record.CreatedAt.Format(formatOut), leftAmount, rightAmount)
			subtotal += uint32(record.Amount)
		}
	}

	leftSubtotal, rightSubtotal := splitAmount(subtotal)
	msg.Text = fmt.Sprintf("Subtotal: %s.%s\u20AC\n\n", leftSubtotal, rightSubtotal) + msg.Text
}

func (c *command) validateInput(input string) []string {
	matches := c.rgx.FindAllStringSubmatch(input, 1)
	if len(matches) == 0 {
		return nil
	}
	return matches[0]
}

func (c *command) becomeLast() {
	c.child = 0
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

func splitAmount(amount any) (left string, rignt string) {

	switch amount := amount.(type) {
	case string:
		splitedAmount := strings.Split(amount, ".")
		left = splitedAmount[0]
		if len(splitedAmount) == 1 {
			rignt = "00"
		} else if len(splitedAmount[1]) == 1 {
			rignt = splitedAmount[1] + "0"
		} else {
			rignt = splitedAmount[1]
		}
	case uint32:
		left = strconv.FormatUint(uint64(amount/100), 10)
		rignt = strconv.FormatUint(uint64(amount%100), 10)
		if len(rignt) == 1 {
			rignt = "0" + rignt
		}
	case uint64:
		left = strconv.FormatUint(uint64(amount/100), 10)
		rignt = strconv.FormatUint(uint64(amount%100), 10)
		if len(rignt) == 1 {
			rignt = "0" + rignt
		}
	}

	return left, rignt
}
