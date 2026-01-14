package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"project_3sem/internal/middleware"
	"project_3sem/internal/responses"
	"time"
)

func (h *OrderHandler) GetOrderStatuses(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.IdKey).(string)

	// Получаем все сайты пользователя
	sites, err := h.RepoSites.GetUserSites(userId)
	if err != nil {
		log.Printf("Ошибка получения сайтов: %v", err)
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	var allOrders []map[string]interface{}

	// Для каждого сайта получаем заказы
	for _, site := range sites {
		orders, err := h.RepoOrders.GetOrdersBySiteID(site.ID)
		if err != nil {
			log.Printf("Ошибка получения заказов для сайта %s: %v", site.ID, err)
			continue
		}

		for _, order := range orders {
			items := make(map[string]interface{})
			if err := json.Unmarshal(order.Items, &items); err != nil {
				log.Printf("Ошибка парсинга items для заказа %s: %v", order.ID, err)
				continue
			}

			orderMap := map[string]interface{}{
				"id":           order.ID.String(),
				"site_id":      site.ID,
				"site_name":    site.Subdomain,
				"items":        items,
				"total_amount": order.TotalAmount,
				"status":       order.Status,
				"created_at":   order.CreatedAt.Format(time.RFC3339),
			}
			allOrders = append(allOrders, orderMap)
		}
	}

	resp := map[string]interface{}{
		"orders": allOrders,
	}

	responses.SendJSONResp(w, resp, http.StatusOK)
}

func (h *OrderHandler) GetSellerBalances(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.IdKey).(string)

	sites, err := h.RepoSites.GetUserSites(userId)
	if err != nil {
		log.Printf("Ошибка получения сайтов: %v", err)
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	balances := make(map[string]float64)

	for _, site := range sites {
		// Суммируем успешные платежи по заказам
		paidOrders, err := h.RepoOrders.GetPaidOrdersAmountBySiteID(site.ID)
		if err != nil {
			log.Printf("Ошибка получения суммы оплаченных заказов для сайта %s: %v", site.ID, err)
			continue
		}

		balances[site.Subdomain] = paidOrders
	}

	resp := map[string]interface{}{
		"balances": balances,
	}

	responses.SendJSONResp(w, resp, http.StatusOK)
}
