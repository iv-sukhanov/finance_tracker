package bot

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/sirupsen/logrus"
)

// contains command IDs and their properties:
//
//   - ID: command ID
//
//   - isBase: whether the command is a first command in the sequence
//
//   - rgx: regular expression that defines the expected input for the command
//
//   - action: function that will be called on the input when the command is selected
//
//   - child: ID of the next command in the sequence, if 0, it means the command is the last one
type command struct {
	ID     int
	isBase bool
	rgx    *regexp.Regexp
	action func(input []string, batch *any, srvc service.ServiceInterface, log *logrus.Logger, sender Sender, cl *client, cmd *command)

	child int
}

const (
	formatOut = "Monday, 02 Jan, 15:04"
	formatIn  = "02.01.2006"

	CommandAddCategory    = "\U0000270Fadd category"
	CommandAddRecord      = "\U0000270Fadd record"
	CommandShowCategories = "\U0001F9FEshow categories"
	CommandShowRecords    = "\U0001F9FEshow records"

	CallbackDataYesRecordsExel    = "yes_records_exel"
	CallbackDataNoRecordsExel     = "no_records_exel"
	CallbackDataYesCategoriesExel = "yes_categories_exel"
	CallbackDataNoCategoriesExel  = "no_categories_exel"

	filename = "report.xlsx"
)

var (
	// translates string commands to command IDs
	commandsToIDs = map[string]int{
		CommandAddCategory:    1,
		CommandAddRecord:      3,
		CommandShowCategories: 4,
		CommandShowRecords:    5,
	}

	// contains replies for each base command
	commandReplies = map[int]string{
		1: MessageAddCategory,
		3: MessageAddRecord,
		4: MessageShowCategories,
		5: MessageShowRecords,
	}

	// contains all registered commands
	// if a new command is added, it should be registered here
	commandsByIDs = map[int]command{
		1: {
			ID:     1,
			isBase: true,
			rgx:    regexp.MustCompile(`^(?<category_name>[a-zA-Z0-9 ]{1,20})$`),
			action: addCategoryAction,
			child:  2,
		},
		2: {
			ID:     2,
			isBase: false,
			rgx:    regexp.MustCompile(`^(?<category_descr>[a-zA-Z0-9 .,!]+)$`),
			action: addCategoryDescriptionAction,
			child:  0,
		},
		3: {
			ID:     3,
			isBase: true,
			rgx:    regexp.MustCompile(`^\s*(?P<category>[a-zA-Z0-9]{1,10})\s*(?P<amount>\d+(?:\.\d{1,2})?)(?:\s+(?<description>[a-zA-Z0-9 ]+))?$`),
			action: addRecordAction,
			child:  0,
		},
		4: {
			ID:     4,
			isBase: true,
			rgx:    regexp.MustCompile(`^(?:(?P<number>\d+)|(?P<category_or_all>[a-zA-Z0-9]{1,10}))(?:\s+(?P<isfull>full))?$`),
			action: showCategoriesAction,
			child:  8,
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
			child:  7,
		},
		7: {
			ID:     7,
			isBase: false,
			rgx: regexp.MustCompile(
				`^(?P<y_or_n>(?:` + CallbackDataYesRecordsExel + `)|(?:` + CallbackDataNoRecordsExel + `))$`,
			),
			action: returnRecordsExelAction,
			child:  0,
		},
		8: {
			ID:     8,
			isBase: false,
			rgx: regexp.MustCompile(
				`^(?P<y_or_n>(?:` + CallbackDataYesCategoriesExel + `)|(?:` + CallbackDataNoCategoriesExel + `))$`,
			),
			action: returnCategoriesExelAction,
			child:  0,
		},
	}

	// inline keyboard asking the user if they want to receive an EXEL file
	// with the records
	wantExelRecordsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", CallbackDataYesRecordsExel),
			tgbotapi.NewInlineKeyboardButtonData("No", CallbackDataNoRecordsExel),
		),
	)

	// inline keyboard asking the user if they want to receive an EXEL file
	// with the categories
	wantExelCategoriesKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", CallbackDataYesCategoriesExel),
			tgbotapi.NewInlineKeyboardButtonData("No", CallbackDataNoCategoriesExel),
		),
	)
)

// action function for the add category command, id 1
//
// it takes takes the input category name, stores it in the batch, and prompts the user to send the category description
func addCategoryAction(input []string, batch *any, _ service.ServiceInterface, log *logrus.Logger, sender Sender, cl *client, cmd *command) {

	// specified regex allways returns 2 tokens, so this check may be redundant,
	// but in case of future changes, it is better to keep it, to catch invalid regex changes
	// or to catch some errors I am unaware of
	if len(input) != 2 {
		log.Error("wrong input for add category command")
		cmd.becomeLast()
		msg := tgbotapi.NewMessage(cl.chanID, MessageInvalidNumberOfTockensAction+"\n"+internalErrorAditionalInfo)
		msg.ReplyMarkup = baseKeyboard
		sender.Send(msg)
		return
	}

	(*batch).(*ftracker.SpendingCategory).Category = input[1]
	log.Debug("action on add category command")
	sender.Send(
		tgbotapi.NewMessage(cl.chanID, MessageAddCategoryDescription),
	)
}

// action function for the add description to a category command, id 2
//
// it takes takes the input description and puts the category to the service.repository
func addCategoryDescriptionAction(input []string, batch *any, srvc service.ServiceInterface, log *logrus.Logger, sender Sender, cl *client, cmd *command) {

	// specified regex allways returns 2 tokens, so this check may be redundant,
	// but in case of future changes, it is better to keep it, to catch invalid regex changes
	// or to catch some errors I am unaware of
	if len(input) != 2 {
		log.Error("wrong input for add category command")
		msg := tgbotapi.NewMessage(cl.chanID, MessageInvalidNumberOfTockensAction+"\n"+internalErrorAditionalInfo)
		msg.ReplyMarkup = baseKeyboard
		sender.Send(msg)
		return
	}

	log.Debug("action on add category description command")

	(*batch).(*ftracker.SpendingCategory).Description = input[1]

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

	categoryToAdd := *(*batch).(*ftracker.SpendingCategory)
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

// action function for the add record, id 3
//
// it takes the input category name, amount and description, and adds the record to the service.repository
func addRecordAction(input []string, batch *any, srvc service.ServiceInterface, log *logrus.Logger, sender Sender, cl *client, cmd *command) {

	// specified regex allways returns 4 tokens, so this check may be redundant,
	// but in case of future changes, it is better to keep it, to catch invalid regex changes,
	// or to catch some errors I am unaware of
	if len(input) != 4 {
		log.Error("wrong tocken number for add record command")
		return
	}
	recordCategory := input[1:2]
	recordAmountLeft, recordAmountRight := utils.ExtractAmountParts(input[2])

	recordDescription := input[3]
	if len(recordDescription) == 0 {
		recordDescription = "spending"
	}

	msg := tgbotapi.NewMessage(cl.chanID, "")
	msg.ReplyMarkup = baseKeyboard
	defer func() {
		sender.Send(msg)
	}()

	if recordAmountLeft == "0" && recordAmountRight == "00" { //zero amount
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
		msg.Text = MessageNoCategoryFound
		return
	}

	(*batch).(*ftracker.SpendingRecord).CategoryGUID = categories[0].GUID
	amount, err := strconv.ParseUint(recordAmountLeft+recordAmountRight, 10, 32)
	if err != nil {
		log.WithError(err).Error("error on parsing amount")
		msg.Text = MessageAmountError + "\n" + internalErrorAditionalInfo
		return
	}
	(*batch).(*ftracker.SpendingRecord).Amount = uint32(amount)
	(*batch).(*ftracker.SpendingRecord).Description = recordDescription

	recordToAdd := *(*batch).(*ftracker.SpendingRecord)
	_, err = srvc.AddRecords([]ftracker.SpendingRecord{recordToAdd})
	if err != nil {
		log.WithError(err).Error("error on add record")
		msg.Text = MessageDatabaseError + "\n" + internalErrorAditionalInfo
	} else {
		msg.Text = MessageRecordSuccess
	}
}

// action function for the show categories command, id 4
//
// it takes the number of categories to display from the user and shows the categories from the service.repository
// then it asks the user if they want to receive an EXEL file with the categories
func showCategoriesAction(input []string, batch *any, srvc service.ServiceInterface, log *logrus.Logger, sender Sender, cl *client, cmd *command) {

	// specified regex allways returns 4 tokens, so this check may be redundant,
	// but in case of future changes, it is better to keep it, to catch invalid regex changes,
	// or to catch some errors I am unaware of
	if len(input) != 4 {
		log.Error("wrong tocken number for add record command")
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
		srvc.SpendingCategoriesWithOrder(service.OrderCategoriesByUpdatedAt, false),
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

	*batch = categories
	msg.Text = "Your categories:\n"
	if addDescription {
		for i, category := range categories {
			leftAmount, rightAmount := utils.ExtractAmountParts(category.Amount)
			msg.Text += fmt.Sprintf(MessageShowCategoriesFormatFull, i+1, category.Category, leftAmount, rightAmount, category.Description)
		}
	} else {
		for i, category := range categories {
			leftAmount, rightAmount := utils.ExtractAmountParts(category.Amount)
			msg.Text += fmt.Sprintf(MessageShowCategoriesFormat, i+1, category.Category, leftAmount, rightAmount)
		}
	}

	msg.Text += MessageWantEXEL
	msg.ReplyMarkup = wantExelCategoriesKeyboard
}

// action function for the show records command, id 5
//
// it takes the category name from the user and prompts then to input the time boundaries
func showRecordsAction(input []string, batch *any, srvc service.ServiceInterface, log *logrus.Logger, sender Sender, cl *client, cmd *command) {

	// specified regex allways returns 2 tokens, so this check may be redundant,
	// but in case of future changes, it is better to keep it, to catch invalid regex changes,
	// or to catch some errors I am unaware of
	if len(input) != 2 {
		log.Error("wrong tocken number for show records command")
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
		msg.ReplyMarkup = baseKeyboard
		cmd.becomeLast()
		return
	}

	if len(categories) == 0 {
		msg.Text = MessageNoCategoryFound
		msg.ReplyMarkup = baseKeyboard
		cmd.becomeLast()
		return
	}
	(*batch).(*repository.RecordOptions).CategoryGUIDs = []uuid.UUID{categories[0].GUID}
	msg.Text = MessageAddTimeDetails
}

// action function for the get time bounaries command, id 6
//
// it takes the the nubmer of records to display, the time boundaries and it the desctiption needed,
// then it diplays the records from the service.repository and asks the user if they want to receive an EXEL file
func getTimeBoundariesAction(input []string, batch *any, srvc service.ServiceInterface, log *logrus.Logger, sender Sender, cl *client, cmd *command) {

	// specified regex allways returns 6 tokens, so this check may be redundant,
	// but in case of future changes, it is better to keep it, to catch invalid regex changes,
	// or to catch some errors I am unaware of
	if len(input) != 6 {
		log.Error("wrong tocken number for set time boundaries command")
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

		if input[4] == "" {
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
			log.Error("invalid token for ymd time boundaries")
			msg.Text = MessageInvalidFixedTime + "\n" + internalErrorAditionalInfo
			return
		}
	}

	log.Debug("time boundaries: ", timeFrom, timeTo)
	recordOption := *(*batch).(*repository.RecordOptions)
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

	*batch = records
	var subtotal uint32 = 0
	if addDescription {
		for _, record := range records {
			leftAmount, rightAmount := utils.ExtractAmountParts(record.Amount)
			msg.Text += fmt.Sprintf(MessageShowRecordsFormatFull, record.CreatedAt.Format(formatOut), leftAmount, rightAmount, record.Description) //mb updated?
			subtotal += uint32(record.Amount)
		}
	} else {
		for _, record := range records {
			leftAmount, rightAmount := utils.ExtractAmountParts(record.Amount)
			msg.Text += fmt.Sprintf(MessageShowRecordsFormat, record.CreatedAt.Format(formatOut), leftAmount, rightAmount)
			subtotal += uint32(record.Amount)
		}
	}

	leftSubtotal, rightSubtotal := utils.ExtractAmountParts(subtotal)
	msg.Text = fmt.Sprintf(MessageShowRecordsFormatHeader, leftSubtotal, rightSubtotal) +
		msg.Text +
		"\n" +
		MessageWantEXEL

	msg.ReplyMarkup = wantExelRecordsKeyboard
}

// action function for the exel records command, id 7
//
// it retrieves the records from the batch and creates an EXEL file with them
// then it sends the file to the user
func returnRecordsExelAction(input []string, batch *any, service service.ServiceInterface, log *logrus.Logger, sender Sender, cl *client, cmd *command) {

	msg := tgbotapi.NewMessage(cl.chanID, "")
	msg.ReplyMarkup = baseKeyboard
	defer func() {
		sender.Send(msg)
	}()

	// specified regex allways returns 2 tokens, so this check may be redundant,
	// but in case of future changes, it is better to keep it, to catch invalid regex changes,
	// or to catch some errors I am unaware of
	if len(input) != 2 {
		log.Error("wrong callback input")
		msg.Text = MessageInvalidNumberOfTockensAction + "\n" + internalErrorAditionalInfo
		return
	}
	log.Debug("action on return records exel command, got: ", input[1])

	if input[1] == CallbackDataNoRecordsExel {
		msg.Text = MessageRecordsExelNo
		return
	}

	records, ok := (*batch).([]ftracker.SpendingRecord)
	if !ok {
		log.Errorf("wrong batch type for exel: %T", *batch)
		msg.Text = MessageInternalError + "\n" + internalErrorAditionalInfo
		return
	}

	file, err := service.CreateExelFromRecords(records)
	if err != nil {
		log.WithError(err).Error("error on create exel")
		msg.Text = MessageExelError + "\n" + internalErrorAditionalInfo
		return
	}
	var buffer bytes.Buffer
	err = file.Write(&buffer)
	if err != nil {
		log.WithError(err).Error("error on upload exel")
		msg.Text = MessageExelError + "\n" + internalErrorAditionalInfo
		return
	}

	document := tgbotapi.NewDocument(cl.chanID, tgbotapi.FileBytes{
		Name:  filename,
		Bytes: buffer.Bytes(),
	})
	msg.Text = MessageRecordsExelYes
	sender.SendDoc(document)
}

// action function for the exel categories command, id 8
//
// it retrieves the categories from the batch and creates an EXEL file with them
// then it sends the file to the user
func returnCategoriesExelAction(input []string, batch *any, service service.ServiceInterface, log *logrus.Logger, sender Sender, cl *client, cmd *command) {

	msg := tgbotapi.NewMessage(cl.chanID, "")
	msg.ReplyMarkup = baseKeyboard
	defer func() {
		sender.Send(msg)
	}()

	// specified regex allways returns 2 tokens, so this check may be redundant,
	// but in case of future changes, it is better to keep it, to catch invalid regex changes,
	// or to catch some errors I am unaware of
	if len(input) != 2 {
		log.Error("wrong callback input")
		msg.Text = MessageInvalidNumberOfTockensAction + "\n" + internalErrorAditionalInfo
		return
	}
	log.Debug("action on return categories exel command, got: ", input[1])

	if input[1] == CallbackDataNoCategoriesExel {
		msg.Text = MessageRecordsExelNo
		return
	}

	categories, ok := (*batch).([]ftracker.SpendingCategory)
	if !ok {
		log.Errorf("wrong batch type for exel: %T", *batch)
		msg.Text = MessageInternalError + "\n" + internalErrorAditionalInfo
		return
	}

	file, err := service.CreateExelFromCategories(categories)
	if err != nil {
		log.WithError(err).Error("error on create exel")
		msg.Text = MessageExelError + "\n" + internalErrorAditionalInfo
		return
	}
	var buffer bytes.Buffer
	err = file.Write(&buffer)
	if err != nil {
		log.WithError(err).Error("error on upload exel")
		msg.Text = MessageExelError + "\n" + internalErrorAditionalInfo
		return
	}

	document := tgbotapi.NewDocument(cl.chanID, tgbotapi.FileBytes{
		Name:  filename,
		Bytes: buffer.Bytes(),
	})
	msg.Text = MessageRecordsExelYes
	sender.SendDoc(document)
}

// validateInput function checks if the input matches the command's regex
// and returns the matched group if it does
func (c *command) validateInput(input string) []string {
	matches := c.rgx.FindAllStringSubmatch(input, 1)
	if len(matches) == 0 {
		return nil
	}
	return matches[0]
}

// in case we want to exit the command sequence
// we set the child to 0, so the command will become the last one
func (c *command) becomeLast() {
	c.child = 0
}

// checks if the command is the last one in the sequence
func (c *command) isLast() bool {
	return c.child == 0
}

// returns the next command in the sequence
func (c *command) next() command {
	return commandsByIDs[c.child]
}

// function to check if the input form user is a command
// it is applied only if there is no currently running process
func isCommand(s string) (command, bool) {
	id, ok := commandsToIDs[s]
	if !ok {
		return command{}, false
	}
	return commandsByIDs[id], true
}

// inits the batch for the command
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
