package keyboard

import (
	"GiftBuyer/internal/model"
	"fmt"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func GoToChannelsKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("channels"),
		),
	)
}

func SettingsKeyboard(settings *model.UserSettings) *telego.InlineKeyboardMarkup {
	var turnAutoBuyButtonText string
	if settings.AutoBuyEnabled {
		turnAutoBuyButtonText = "🔴 Выключить"
	} else {
		turnAutoBuyButtonText = "🟢 Включить"
	}

	var onlyPremiumGiftButtonText string
	if settings.OnlyPremiumGift {
		onlyPremiumGiftButtonText = "🔴 Выключить premium подарки"
	} else {
		onlyPremiumGiftButtonText = "🟢 Включить premium подарки"
	}

	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(turnAutoBuyButtonText).WithCallbackData("toggle_auto_buy"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(onlyPremiumGiftButtonText).WithCallbackData("only_premium_gift"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔢 Лимит цены ОТ").WithCallbackData("price_from"),
			tu.InlineKeyboardButton("🔢 Лимит цены ДО").WithCallbackData("price_to"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔢 Лимит саплая").WithCallbackData("supply_limit"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Циклы автопокупки").WithCallbackData("auto_buy_cycles"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Каналы").WithCallbackData("channels"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("back_to_start"),
		),
	)
}

func SetPriceFromKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 15").WithCallbackData("set_price_from=15"),
			tu.InlineKeyboardButton("⭐️ 25").WithCallbackData("set_price_from=25"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 50").WithCallbackData("set_price_from=50"),
			tu.InlineKeyboardButton("⭐️ 100").WithCallbackData("set_price_from=100"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 200").WithCallbackData("set_price_from=200"),
			tu.InlineKeyboardButton("⭐️ 500").WithCallbackData("set_price_from=500"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 1000").WithCallbackData("set_price_from=1000"),
			tu.InlineKeyboardButton("⭐️ 1500").WithCallbackData("set_price_from=1500"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 2000").WithCallbackData("set_price_from=2000"),
			tu.InlineKeyboardButton("⭐️ 2500").WithCallbackData("set_price_from=2500"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 3000").WithCallbackData("set_price_from=3000"),
			tu.InlineKeyboardButton("⭐️ 5000").WithCallbackData("set_price_from=5000"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 10000").WithCallbackData("set_price_from=10000"),
			tu.InlineKeyboardButton("⭐️ 20000").WithCallbackData("set_price_from=20000"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ Убрать лимит").WithCallbackData("set_price_from=nil"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("settings"),
		),
	)
}

func SetPriceToKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 15").WithCallbackData("set_price_to=15"),
			tu.InlineKeyboardButton("⭐️ 25").WithCallbackData("set_price_to=25"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 50").WithCallbackData("set_price_to=50"),
			tu.InlineKeyboardButton("⭐️ 100").WithCallbackData("set_price_to=100"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 200").WithCallbackData("set_price_to=200"),
			tu.InlineKeyboardButton("⭐️ 500").WithCallbackData("set_price_to=500"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 1000").WithCallbackData("set_price_to=1000"),
			tu.InlineKeyboardButton("⭐️ 1500").WithCallbackData("set_price_to=1500"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 2000").WithCallbackData("set_price_to=2000"),
			tu.InlineKeyboardButton("⭐️ 2500").WithCallbackData("set_price_to=2500"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 3000").WithCallbackData("set_price_to=3000"),
			tu.InlineKeyboardButton("⭐️ 5000").WithCallbackData("set_price_to=5000"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 7500").WithCallbackData("set_price_to=7500"),
			tu.InlineKeyboardButton("⭐️ 10000").WithCallbackData("set_price_to=10000"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ 15000").WithCallbackData("set_price_to=15000"),
			tu.InlineKeyboardButton("⭐️ 20000").WithCallbackData("set_price_to=20000"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ Убрать лимит").WithCallbackData("set_price_to=nil"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("settings"),
		),
	)
}

func SetSupplyLimitKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("500").WithCallbackData("set_supply_limit=500"),
			tu.InlineKeyboardButton("1000").WithCallbackData("set_supply_limit=1000"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("1500").WithCallbackData("set_supply_limit=1500"),
			tu.InlineKeyboardButton("1999").WithCallbackData("set_supply_limit=1999"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("2000").WithCallbackData("set_supply_limit=2000"),
			tu.InlineKeyboardButton("3000").WithCallbackData("set_supply_limit=3000"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("5000").WithCallbackData("set_supply_limit=5000"),
			tu.InlineKeyboardButton("7500").WithCallbackData("set_supply_limit=7500"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("10000").WithCallbackData("set_supply_limit=10000"),
			tu.InlineKeyboardButton("15000").WithCallbackData("set_supply_limit=15000"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("25000").WithCallbackData("set_supply_limit=25000"),
			tu.InlineKeyboardButton("50000").WithCallbackData("set_supply_limit=50000"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("100000").WithCallbackData("set_supply_limit=100000"),
			tu.InlineKeyboardButton("250000").WithCallbackData("set_supply_limit=250000"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⭐️ Убрать лимит").WithCallbackData("set_supply_limit=nil"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("settings"),
		),
	)
}

func SetAutoBuyCyclesKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("1").WithCallbackData("set_auto_buy_cycles=1"),
			tu.InlineKeyboardButton("2").WithCallbackData("set_auto_buy_cycles=2"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("3").WithCallbackData("set_auto_buy_cycles=3"),
			tu.InlineKeyboardButton("5").WithCallbackData("set_auto_buy_cycles=5"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("10").WithCallbackData("set_auto_buy_cycles=10"),
			tu.InlineKeyboardButton("20").WithCallbackData("set_auto_buy_cycles=20"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("30").WithCallbackData("set_auto_buy_cycles=30"),
			tu.InlineKeyboardButton("50").WithCallbackData("set_auto_buy_cycles=50"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("75").WithCallbackData("set_auto_buy_cycles=75"),
			tu.InlineKeyboardButton("100").WithCallbackData("set_auto_buy_cycles=100"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔁 Бесконечно").WithCallbackData("set_auto_buy_cycles=infinite"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("settings"),
		),
	)
}

func ChannelsKeyboard(channels []string) *telego.InlineKeyboardMarkup {
	var rows [][]telego.InlineKeyboardButton

	// Если каналы заданы — добавляем кнопки
	if channels != nil && len(channels) > 0 {
		for i := 0; i < len(channels); i += 2 {
			row := make([]telego.InlineKeyboardButton, 0, 2)

			rows = append(rows, tu.InlineKeyboardRow(
				tu.InlineKeyboardButton((channels)[i]).WithCallbackData("channel="+(channels)[i]),
			))

			if i+1 < len(channels) {
				rows = append(rows, tu.InlineKeyboardRow(
					tu.InlineKeyboardButton((channels)[i+1]).WithCallbackData("channel="+(channels)[i+1]),
				))
			}

			rows = append(rows, row)
		}
	}

	// Добавляем кнопку "Добавить канал"
	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("➕ Добавить канал").WithCallbackData("add_channel"),
	))

	// Добавляем кнопку "Назад"
	rows = append(rows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("settings"),
	))

	return &telego.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func ChannelSettingsKeyboard(username string) *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Удалить канал").WithCallbackData(fmt.Sprintf("delete_channel=%s", username)),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🔙 Назад").WithCallbackData("channels"),
		),
	)
}
