package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	EMOJI_COIN         = "\U0001FA99"
	EMOJI_SMILE        = "\U0001F642"
	EMOJI_SUNGLASSES   = "\U0001F60E"
	EMOJI_DONT_KNOW    = "\U0001F937"
	EMOJI_SAD          = "\U0001F63F"
	EMOJI_BICEPS       = "\U0001F4AA"
	EMOJI_BUTTON_START = "\U000025B6"
	EMOJI_BUTTON_END   = "\U000025C0"

	BUTTON_TEXT_PRINT_INTRO       = EMOJI_BUTTON_START + "View Introduction" + EMOJI_BUTTON_END
	BUTTON_TEXT_SKIP_INTRO        = EMOJI_BUTTON_START + "Skip Introduction" + EMOJI_BUTTON_END
	BUTTON_TEXT_BALANCE           = EMOJI_BUTTON_START + "Current Balance" + EMOJI_BUTTON_END
	BUTTON_TEXT_USEFUL_ACTIVITIES = EMOJI_BUTTON_START + "Useful Activities" + EMOJI_BUTTON_END
	BUTTON_TEXT_REWARDS           = EMOJI_BUTTON_START + "Rewards" + EMOJI_BUTTON_END
	BUTTON_TEXT_PRINT_MENU        = EMOJI_BUTTON_START + "MAIN MENU" + EMOJI_BUTTON_END

	BUTTON_CODE_PRINT_INTRO       = "print_intro"
	BUTTON_CODE_SKIP_INTRO        = "skip_intro"
	BUTTON_CODE_BALANCE           = "show_balance"
	BUTTON_CODE_USEFUL_ACTIVITIES = "show_useful_activities"
	BUTTON_CODE_REWARDS           = "show_rewards"
	BUTTON_CODE_PRINT_MENU        = "print_menu"

	TOKEN_NAME_IN_OS             = "6599775164:AAFenKJA1aPvIMeiiN4mz3Cnt1FckauXjjY"
	UPDATE_CONFIG_TIMEOUT        = 60
	MAX_USER_COINS        uint16 = 500
)

var gBot *tgbotapi.BotAPI
var gToken string
var gChatId int64

var gUsersInChat Users

var gUsefulActivities = Activities{
	// Саморазвитие
	{"yoga", "Йога (15 минут)", 1},
	{"meditation", "Медитация (15 минут)", 1},
	{"language", "Изучение языка (15 минут)", 1},
	{"swimming", "Плавание (15 минут)", 1},
	{"walk", "Ходьба (15 минут)", 1},
	{"chores", "работа", 1},

	// Работа
	{"work_learning", "Изучаем рабочие материалы(15 минут)", 1},
	{"portfolio_work", "Работа над проектом в портфолио. (15 минут)", 1},
	{"resume_edit", "Улучшение резюме (15 минут)", 1},

	// Творчество
	{"creative", "Творческая деятельность (15 минут)", 1},
	{"reading", "Чтение худ. литературы (15 минут)", 1},
}

var gRewards = Activities{
	// Развлечение
	{"watch_series", "Просмотр сериала (1 эпизод)", 10},
	{"watch_movie", "Смотреть фильм (1 фильм)", 30},
	{"social_nets", "Листать соц. сети (30 минут)", 10},

	// Еда
	{"eat_sweets", "300 калорий конфет", 60},
}

type User struct {
	id    int64
	name  string
	coins uint16
}
type Users []*User

type Activity struct {
	code, name string
	coins      uint16
}
type Activities []*Activity

func init() {
	_ = os.Setenv(TOKEN_NAME_IN_OS, "6599775164:AAFenKJA1aPvIMeiiN4mz3Cnt1FckauXjjY")
	gToken = os.Getenv(TOKEN_NAME_IN_OS)

	if gToken = os.Getenv(TOKEN_NAME_IN_OS); gToken == "" {
		panic(fmt.Errorf(`не удалось загрузить переменную среды "%s"`, TOKEN_NAME_IN_OS))
	}

	var err error
	if gBot, err = tgbotapi.NewBotAPI(gToken); err != nil {
		log.Panic(err)
	}
	gBot.Debug = true
}

func isStartMessage(update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text == "/start"
}

func isCallbackQuery(update *tgbotapi.Update) bool {
	return update.CallbackQuery != nil && update.CallbackQuery.Data != ""
}

func delay(seconds uint8) {
	time.Sleep(time.Second * time.Duration(seconds))
}

func sendStringMessage(msg string) {
	gBot.Send(tgbotapi.NewMessage(gChatId, msg))
}

func sendMessageWithDelay(delayInSec uint8, message string) {
	sendStringMessage(message)
	delay(delayInSec)
}

func printIntro(update *tgbotapi.Update) {
	sendMessageWithDelay(2, "Привет! "+EMOJI_SUNGLASSES)
	sendMessageWithDelay(7, "Существует множество полезных действий, регулярно совершая которые, мы улучшаем качество своей жизни. Однако зачастую веселее, проще или вкуснее сделать что-то вредное. Не так ли?")
	sendMessageWithDelay(7, "Мы предпочтем позалипать в YouTube Shorts вместо пары, купить M&M’s вместо овощей или поваляться в постели вместо того, чтобы заняться йогой.")
	sendMessageWithDelay(1, EMOJI_SAD)
	sendMessageWithDelay(10, "Каждый играл хотя бы в одну игру, где нужно прокачивать персонажа, делая его сильнее, умнее или красивее. Это приятно, потому что каждое действие приносит результат. Однако в реальной жизни систематические действия только со временем начинают становиться заметными. Давайте это изменим, ладно?")
	sendMessageWithDelay(1, EMOJI_SMILE)
	sendMessageWithDelay(14, "Перед вами две таблицы: «Полезные активности» и «Награды». В первой таблице перечислены простые короткие действия, за выполнение каждого из которых вы заработаете указанное количество монет. Во второй таблице вы увидите список действий, которые вы сможете выполнять только после оплаты за них монетами, заработанными на предыдущем шаге.")
	sendMessageWithDelay(1, EMOJI_COIN)
	sendMessageWithDelay(10, "Например, вы проводите полчаса, занимаясь йогой, за что получаете 2 монеты. После этого вам предстоит 2 часа изучения программирования, за которые вы получите 8 монет. Теперь вы можете посмотреть 1 серию аниме и безубыточно. Это так просто!")
	sendMessageWithDelay(6, "Отмечайте выполненные полезные действия, чтобы не потерять свои монеты. И не забудьте «купить» награду, прежде чем сделать это.")
}

func getKeyboardRow(buttonText, buttonCode string) []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonCode))
}

func askToPrintIntro() {
	msg := tgbotapi.NewMessage(gChatId, "Во вступительных сообщениях вы можете узнать назначение этого бота и правила игры. Что вы думаете?")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_PRINT_INTRO, BUTTON_CODE_PRINT_INTRO),
		getKeyboardRow(BUTTON_TEXT_SKIP_INTRO, BUTTON_CODE_SKIP_INTRO),
	)
	gBot.Send(msg)
}

func showMenu() {
	msg := tgbotapi.NewMessage(gChatId, "Выберите один из вариантов:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_BALANCE, BUTTON_CODE_BALANCE),
		getKeyboardRow(BUTTON_TEXT_USEFUL_ACTIVITIES, BUTTON_CODE_USEFUL_ACTIVITIES),
		getKeyboardRow(BUTTON_TEXT_REWARDS, BUTTON_CODE_REWARDS),
	)
	gBot.Send(msg)
}

func showBalance(user *User) {
	msg := fmt.Sprintf("%s, ваш кошелек пуст %s \nСделай что-нибудь полезное, чтобы заработать монеты", user.name, EMOJI_DONT_KNOW)
	if coins := user.coins; coins > 0 {
		msg = fmt.Sprintf("%s, у вас есть %d %s", user.name, coins, EMOJI_COIN)
	}
	sendStringMessage(msg)
	showMenu()
}

func callbackQueryFromIsMissing(update *tgbotapi.Update) bool {
	return update.CallbackQuery == nil || update.CallbackQuery.From == nil
}

func getUserFromUpdate(update *tgbotapi.Update) (user *User, found bool) {
	if callbackQueryFromIsMissing(update) {
		return
	}

	userId := update.CallbackQuery.From.ID
	for _, userInChat := range gUsersInChat {
		if userId == userInChat.id {
			return userInChat, true
		}
	}
	return
}

func storeUserFromUpdate(update *tgbotapi.Update) (user *User, found bool) {
	if callbackQueryFromIsMissing(update) {
		return
	}

	from := update.CallbackQuery.From
	user = &User{id: from.ID, name: strings.TrimSpace(from.FirstName + " " + from.LastName), coins: 0}
	gUsersInChat = append(gUsersInChat, user)
	return user, true
}

func showActivities(activities Activities, message string, isUseful bool) {
	activitiesButtonsRows := make([]([]tgbotapi.InlineKeyboardButton), 0, len(activities)+1)
	for _, activity := range activities {
		activityDescription := ""
		if isUseful {
			activityDescription = fmt.Sprintf("+ %d %s: %s", activity.coins, EMOJI_COIN, activity.name)
		} else {
			activityDescription = fmt.Sprintf("- %d %s: %s", activity.coins, EMOJI_COIN, activity.name)
		}
		activitiesButtonsRows = append(activitiesButtonsRows, getKeyboardRow(activityDescription, activity.code))
	}
	activitiesButtonsRows = append(activitiesButtonsRows, getKeyboardRow(BUTTON_TEXT_PRINT_MENU, BUTTON_CODE_PRINT_MENU))

	msg := tgbotapi.NewMessage(gChatId, message)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(activitiesButtonsRows...)
	gBot.Send(msg)
}

func showUsefulActivities() {
	showActivities(gUsefulActivities, "Тыкай полезную активность или вернитесь в главное меню:", true)
}

func showRewards() {
	showActivities(gRewards, "Купи награду или возвращаяйся в главное меню:", false)
}

func findActivity(activities Activities, choiceCode string) (activity *Activity, found bool) {
	for _, activity := range activities {
		if choiceCode == activity.code {
			return activity, true
		}
	}
	return
}

func processUsefulActivity(activity *Activity, user *User) {
	errorMsg := ""
	if activity.coins == 0 {
		errorMsg = fmt.Sprintf(`У активность "%s" не указана стоимость`, activity.name)
	} else if user.coins+activity.coins > MAX_USER_COINS {
		errorMsg = fmt.Sprintf("У тебя не может быть больше %d %s", MAX_USER_COINS, EMOJI_COIN)
	}

	resultMessage := ""
	if errorMsg != "" {
		resultMessage = fmt.Sprintf("%s, прости , но %s %s Твой баланс остается неизменным.", user.name, errorMsg, EMOJI_SAD)
	} else {
		user.coins += activity.coins
		resultMessage = fmt.Sprintf(`%s, эта "%s"деятельность завершена! %d %s был добавлен в ваш аккаунт. Так держать! %s%s Теперь у вас есть %d %s`,
			user.name, activity.name, activity.coins, EMOJI_COIN, EMOJI_BICEPS, EMOJI_SUNGLASSES, user.coins, EMOJI_COIN)
	}
	sendStringMessage(resultMessage)
}

func processReward(activity *Activity, user *User) {
	errorMsg := ""
	if activity.coins == 0 {
		errorMsg = fmt.Sprintf(`награда "%s" не имеет указанной стоимости`, activity.name)
	} else if user.coins < activity.coins {
		errorMsg = fmt.Sprintf(`у тебя сейчас есть %d %s. Ты не можешь себе позволить "%s" за %d %s`, user.coins, EMOJI_COIN, activity.name, activity.coins, EMOJI_COIN)
	}

	resultMessage := ""
	if errorMsg != "" {
		resultMessage = fmt.Sprintf("%s, прости, но %s %s твой баланс остается неизменным, награда недоступна %s", user.name, errorMsg, EMOJI_SAD, EMOJI_DONT_KNOW)
	} else {
		user.coins -= activity.coins
		resultMessage = fmt.Sprintf(`%s, награда"%s" оплачена, начинай! %d %s было списано с вашего счета. Теперь у тебя есть %d %s`, user.name, activity.name, activity.coins, EMOJI_COIN, user.coins, EMOJI_COIN)
	}
	sendStringMessage(resultMessage)
}

func updateProcessing(update *tgbotapi.Update) {
	user, found := getUserFromUpdate(update)
	if !found {
		if user, found = storeUserFromUpdate(update); !found {
			sendStringMessage("Невозможно идентифицировать пользователя")
			return
		}
	}

	choiceCode := update.CallbackQuery.Data
	log.Printf("[%T] %s", time.Now(), choiceCode)

	switch choiceCode {
	case BUTTON_CODE_BALANCE:
		showBalance(user)
	case BUTTON_CODE_USEFUL_ACTIVITIES:
		showUsefulActivities()
	case BUTTON_CODE_REWARDS:
		showRewards()
	case BUTTON_CODE_PRINT_INTRO:
		printIntro(update)
		showMenu()
	case BUTTON_CODE_SKIP_INTRO:
		showMenu()
	case BUTTON_CODE_PRINT_MENU:
		showMenu()
	default:
		if usefulActivity, found := findActivity(gUsefulActivities, choiceCode); found {
			processUsefulActivity(usefulActivity, user)

			delay(2)
			showUsefulActivities()
			return
		}

		if reward, found := findActivity(gRewards, choiceCode); found {
			processReward(reward, user)

			delay(2)
			showRewards()
			return
		}

		log.Printf(`[%T] ERROR: Неизвестный код! "%s"`, time.Now(), choiceCode)
		msg := fmt.Sprintf("%s, Извините, я не знаю код '%s' %s Пожалуйста, сообщи об этой ошибке моему создателю.", user.name, choiceCode, EMOJI_SAD)
		sendStringMessage(msg)
	}
}

func main() {
	log.Printf("Авторизован на аккаунте %s", gBot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = UPDATE_CONFIG_TIMEOUT

	for update := range gBot.GetUpdatesChan(updateConfig) {
		if isCallbackQuery(&update) {
			updateProcessing(&update)
		} else if isStartMessage(&update) {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			gChatId = update.Message.Chat.ID
			askToPrintIntro()
		}
	}
}
