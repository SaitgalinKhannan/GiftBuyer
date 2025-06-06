package utils

import (
	"GiftBuyer/internal/model"
	"fmt"
	"strconv"
)

// FormatAutoBuySettings формирует текст настроек автопокупки для отображения пользователю
func FormatAutoBuySettings(settings *model.UserSettings) string {
	// Форматируем статус автопокупки
	status := "🔴 Выключено"
	if settings.AutoBuyEnabled {
		status = "🟢 Включено"
	}

	// Форматируем лимит цены
	from := "0"
	to := "∞"
	if settings.PriceLimitFrom != nil {
		from = strconv.Itoa(*settings.PriceLimitFrom)
	}
	if settings.PriceLimitTo != nil {
		to = strconv.Itoa(*settings.PriceLimitTo)
	}

	// Форматируем лимит саплая
	supply := "Не задан"
	if settings.SupplyLimit != nil && *settings.SupplyLimit > 0 {
		supply = strconv.Itoa(*settings.SupplyLimit)
	}

	// Форматируем кол-во циклов покупки
	autoBuyCycles := "Бесконечно"
	if settings.AutoBuyCycles > 0 {
		autoBuyCycles = strconv.Itoa(settings.AutoBuyCycles)
	}

	return fmt.Sprintf(
		"<b>⚙️ Настройки автопокупки</b>\n"+
			"<b>Статус:</b> %s\n\n"+
			"<b>Лимит цены:</b>\n"+
			"От %s до %s ⭐️\n\n"+
			"<b>Лимит саплая:</b>\n"+
			"%s\n\n"+
			"<b>Количество циклов:</b>\n"+
			"%s",
		status,
		from,
		to,
		supply,
		autoBuyCycles,
	)
}
