package keyboard

import (
	"fmt"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func GoMainKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("back_to_start"),
		),
	)
}

func StartKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("📝 Вход в аккаунт").WithCallbackData("login"),
		),
		tu.InlineKeyboardRow(
			// tu.InlineKeyboardButton("👤 Профиль").WithCallbackData("profile"),
			tu.InlineKeyboardButton("⚙️ Настройки автопокупки").WithCallbackData("settings"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("👨🏼‍💻 Помощь").WithURL("https://t.me/gmkmnv"),
		),
	)
}

func BuyStarsKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ Купить звезды дёшево").WithURL("https://split.tg/?ref=UQAEBELEbfrTtfyaT1ny28DCdQzSu34-mLv7gY-1czOlZMWL"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("back_to_start"),
		),
	)
}

func TopUpBalanceKeyboard(amount int, link string) *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow( // Row 1
			tu.InlineKeyboardButton(fmt.Sprintf("⭐️ Оплатить %d STARS", amount)).WithURL(link),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ Купить звезды дёшево").WithURL("https://split.tg/?ref=UQAEBELEbfrTtfyaT1ny28DCdQzSu34-mLv7gY-1czOlZMWL"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("back_to_start"),
		),
	)
}

func ProfileKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ Пополнить баланс").WithCallbackData("top_up_balance"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("back_to_start"),
		),
	)
}
