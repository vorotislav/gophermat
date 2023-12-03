package models

import "time"

// Balance хранит информацию о текущем баланса баллов пользователя и списанные баллы. всё хранится в копейках.
type Balance struct {
	Current  int `json:"current"`  //
	Withdraw int `json:"withdraw"` //
}

// BalanceWithdraw запрос на списание баллов со счёта.
type BalanceWithdraw struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

// BalanceWithdrawal показывает историю списаний баллов для каждого заказа.
type BalanceWithdrawal struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
