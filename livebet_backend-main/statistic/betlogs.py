import requests
import json

HEADERS = {
	"Content-Type": "application/json",
}

def calc_bet(pair):
	url = "http://188.253.24.91:7010/calc-bet"

	payload = {
		"userId": "0",
		"pair": pair
	}

	response = requests.post(url, data=json.dumps(payload), headers=HEADERS)
	if response.status_code != 200:
		print(f"[CalcBet.ERROR] Не удалось рассчитать ставку: {response.text}")
		return False

	return response.json()

def log_bet(pair, calculatedBet, sum, coef, time):
	url = "http://188.253.24.91:7010/log-test-bet-accept"

	payload = {
		"pair": pair,
		"bet": calculatedBet,
		"sum": sum,
		"coef": coef,
		"time": time,
		"userId": 0
	}

	response = requests.post(url, data=json.dumps(payload), headers=HEADERS)
	if response.status_code != 200:
		print(f"[LogTestBetAccept.ERROR] Не удалось логировать ставку: {response.text}")
		return False

	return True