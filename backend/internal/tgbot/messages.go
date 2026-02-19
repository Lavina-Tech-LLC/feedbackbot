package tgbot

// Bot message templates with i18n support
var messages = map[string]map[string]string{
	"welcome": {
		"en": "👋 Welcome to FeedbackBot!\n\n📝 Send me a message and I'll deliver it **anonymously** to your team admin.\n\n**Commands:**\n/start — Show this welcome message\n/help — Show available commands\n/adminOnly <message> — Send feedback visible only to admin\n\n💡 Just type your feedback and send!",
		"ru": "👋 Добро пожаловать в FeedbackBot!\n\n📝 Отправьте мне сообщение, и я доставлю его **анонимно** вашему администратору.\n\n**Команды:**\n/start — Показать приветствие\n/help — Показать команды\n/adminOnly <сообщение> — Отправить отзыв только для админа\n\n💡 Просто напишите ваш отзыв и отправьте!",
		"uz": "👋 FeedbackBot'ga xush kelibsiz!\n\n📝 Menga xabar yuboring va men uni **anonim** ravishda administratoringizga yetkazaman.\n\n**Buyruqlar:**\n/start — Salom xabarini ko'rsatish\n/help — Buyruqlarni ko'rsatish\n/adminOnly <xabar> — Faqat admin uchun fikr yuborish\n\n💡 Fikringizni yozing va yuboring!",
	},
	"adminOnlyEmpty": {
		"en": "Please write your feedback after /adminOnly.\n\nExample: /adminOnly I think we should improve our standup meetings.",
		"ru": "Пожалуйста, напишите отзыв после /adminOnly.\n\nПример: /adminOnly Я думаю, нам стоит улучшить наши стендапы.",
		"uz": "/adminOnly dan keyin fikringizni yozing.\n\nMisol: /adminOnly Menimcha, standup yig'ilishlarimizni yaxshilashimiz kerak.",
	},
	"emptyMessage": {
		"en": "Please send a text message with your feedback.",
		"ru": "Пожалуйста, отправьте текстовое сообщение с вашим отзывом.",
		"uz": "Iltimos, fikringiz bilan matnli xabar yuboring.",
	},
	"noGroups": {
		"en": "❌ No active groups found. The bot needs to be added to a group first.",
		"ru": "❌ Активные группы не найдены. Сначала нужно добавить бота в группу.",
		"uz": "❌ Faol guruhlar topilmadi. Avval botni guruhga qo'shish kerak.",
	},
	"pickGroup": {
		"en": "📋 Which group is this feedback for?",
		"ru": "📋 Для какой группы этот отзыв?",
		"uz": "📋 Bu fikr qaysi guruh uchun?",
	},
	"feedbackSent": {
		"en": "✅ Your feedback has been submitted anonymously. Thank you!",
		"ru": "✅ Ваш отзыв отправлен анонимно. Спасибо!",
		"uz": "✅ Fikringiz anonim ravishda yuborildi. Rahmat!",
	},
	"feedbackSentAdminOnly": {
		"en": "✅ Your feedback has been sent privately to the admin. It will NOT be posted in the group.",
		"ru": "✅ Ваш отзыв отправлен приватно администратору. Он НЕ будет опубликован в группе.",
		"uz": "✅ Fikringiz maxfiy ravishda administratorga yuborildi. U guruhda JOYLANMAYDI.",
	},
	"sessionExpired": {
		"en": "⏳ Session expired. Please send your feedback again.",
		"ru": "⏳ Сессия истекла. Пожалуйста, отправьте отзыв заново.",
		"uz": "⏳ Sessiya tugadi. Iltimos, fikringizni qaytadan yuboring.",
	},
	"groupNotFound": {
		"en": "❌ Group not found.",
		"ru": "❌ Группа не найдена.",
		"uz": "❌ Guruh topilmadi.",
	},
	"unknownCommand": {
		"en": "🤔 Unknown command. Did you mean to send feedback? Just type your message!\n\nUse /help to see available commands.",
		"ru": "🤔 Неизвестная команда. Хотели отправить отзыв? Просто напишите сообщение!\n\nИспользуйте /help для списка команд.",
		"uz": "🤔 Noma'lum buyruq. Fikr yubormoqchi edingizmi? Xabaringizni yozing!\n\n/help — buyruqlar ro'yxati.",
	},
	"rateLimited": {
		"en": "⏰ You've sent too many messages. Please wait a bit before sending more feedback (max 10/hour).",
		"ru": "⏰ Вы отправили слишком много сообщений. Подождите немного (макс. 10/час).",
		"uz": "⏰ Juda ko'p xabar yubordingiz. Biroz kuting (max 10/soat).",
	},
}

// getMsg returns a localized message, falling back to English
func getMsg(key string, lang string) string {
	if msgs, ok := messages[key]; ok {
		if msg, ok := msgs[lang]; ok {
			return msg
		}
		if msg, ok := msgs["en"]; ok {
			return msg
		}
	}
	return ""
}

// detectLang returns the user's language code, defaulting to "en"
func detectLang(langCode string) string {
	switch langCode {
	case "ru":
		return "ru"
	case "uz":
		return "uz"
	default:
		return "en"
	}
}
