import websockets

from datetime import datetime as dt, timedelta as td, timezone as tz
import asyncio
import json
import time

import config
from categories import generate_categories, categories, Outcome
from roi import calculate_roi
from betlogs import calc_bet, log_bet


class Statistic:
	def __init__(self, url, is_live) -> None:
		self.is_live = is_live
		self.analyzer_url = url
		self.actual_time = 4.5 if is_live else 35.0

		self.filters = config.FILTERS
		self.added_outcomes = []

		self.lad2_outcomes = {}
		self.lad2_blocked = []

	async def connect_to_analyzer(self):
		while True:
			async with websockets.connect(self.analyzer_url) as websocket:
				print("Подключено к Analyzer")

				try:
					await websocket.send(json.dumps(self.filters))
					while True:
						message = await websocket.recv()
						await self.process(message)

				except websockets.ConnectionClosed:
					print("Соединение с WebSocket сервером закрыто")

			await asyncio.sleep(300)

	async def run(self):
		await self.connect_to_analyzer()

	async def process_ladbrokes2(self, pair):
		sport = pair["sportName"]
		bookmaker = pair["second"]["bookmaker"]
		event = pair["first"]

		for outcome in [o for o in pair["outcome"] if isinstance(o, dict)]:
			if "G" in outcome["outcome"]:
				continue

			outcome_key = f"{event['matchId']}|{outcome['outcome']}"

			coef = outcome["score2"]["value"] - 0.05
			outcome["score2"]["value"] = coef

			roi = calculate_roi(
				coef, outcome["score1"]["value"],
				outcome["margin"], outcome["marketType"],
				bookmaker, sport
			)
			outcome["roi"] = roi

			lad2_outcome = self.lad2_outcomes.get(outcome_key)

			if lad2_outcome:
				if roi < 3:
					self.lad2_outcomes.pop(outcome_key)
					print(f"[LADBROKES2] Исход '{outcome_key}' (ROI: {roi}, Коэф: {coef}) уже не подходит")
				else:
					lad2_outcome.update(outcome)
					print(f"[LADBROKES2] Исход '{outcome_key}' (ROI: {roi}, Коэф: {coef}) обновлён")

					if lad2_outcome.updatedAt - lad2_outcome.addedAt > config.MIN_TIME:
						calculatedBet = calc_bet(lad2_outcome.pair)

						if calculatedBet:
							current_time = time.strftime("%M:%S")
							lad2_outcome.pair["second"]["bookmaker"] = "Ladbrokes2"
							is_logged = log_bet(lad2_outcome.pair, calculatedBet, 100, coef, current_time)

							if is_logged:
								self.lad2_outcomes.pop(outcome_key)
								self.lad2_blocked.append(event["matchId"])
								print(f"[LADBROKES2] Исход '{outcome_key}' (ROI: {roi}, Коэф: {coef}) был залогирован")
			else:
				if roi > 3:
					self.lad2_outcomes[outcome_key] = Outcome(outcome_key, pair, outcome)
					print(f"[LADBROKES2] Исход '{outcome_key}' (ROI: {roi}, Коэф: {coef}) добавлен")


	async def process(self, message):
		data = json.loads(message)
		if not data:
			return

		for pair in data:
			# Проверяем актуальность данных
			first_time, second_time = (
				dt.fromisoformat(pair["first"]["createdAt"]),
				dt.fromisoformat(pair["second"]["createdAt"])
			)
			if (dt.now(tz.utc) - first_time) > td(seconds=self.actual_time) or (dt.now(tz.utc) - second_time) > td(seconds=self.actual_time):
				continue

			sport = pair["sportName"]
			bookmaker = pair["second"]["bookmaker"]
			event = pair["first"]

			if bookmaker == "Ladbrokes" and event["matchId"] not in self.lad2_blocked and self.is_live:
				asyncio.create_task(self.process_ladbrokes2(pair))

			for outcome in pair["outcome"]:
				# Пропускаем геймы (ибо на них не можем найти результат)
				if "G" in outcome["outcome"]:
					continue

				outcome_key = f"{event['matchId']}|{outcome['outcome']}"
				roi, coef = outcome["roi"], outcome["score2"]["value"]

				# Если исход уже добавлен в какую-то категорию, ищем её и обновляем исход в категории
				if outcome_key in self.added_outcomes:
					for _, category in categories[sport][bookmaker].items():
						if category.has(outcome_key):
							is_updated = category.update(outcome_key, outcome)

							if is_updated:
								print(f"Исход '{outcome_key}' (ROI: {roi}, Коэф: {coef}) обновлён в категории '{category.name}'")
								updated_outcome = category.get(outcome_key)

								# Проверка насколько долго исход продержался в категории
								if updated_outcome.updatedAt - updated_outcome.addedAt > config.MIN_TIME:
									calculatedBet = calc_bet(updated_outcome.pair)

									if calculatedBet:
										coef = updated_outcome.pair["outcome"]["score2"]["value"]
										current_time = time.strftime("%M:%S")
										is_logged = log_bet(updated_outcome.pair, calculatedBet, category.ID, coef, current_time)

										if is_logged:
											category.delete(outcome_key)
											category.block(event)
											self.added_outcomes.remove(outcome_key)
											print(f"Исход '{outcome_key}' (ROI: {roi}, Коэф: {coef}) был залогирован и удалён из категории '{category.name}'")

							else:
								category.unblock(event)
								print(f"Исход '{outcome_key}' (ROI: {roi}, Коэф: {coef}) уже не подходит под категорию '{category.name}' (удалён)")

				# Если исход не был добавлен в какую-либо категорию, проверяем его на соответствие и добавляем в подходящую
				else:
					for _, category in categories[sport][bookmaker].items():
						if not category.is_blocked(event):
							is_added = category.add(outcome_key, pair, outcome)

							if is_added:
								self.added_outcomes.append(outcome_key)
								category.block(event)
								print(f"Исход '{outcome_key}' (ROI: {roi}, Коэф: {coef}) добавлен в категорию '{category.name}'")
							else:
								# print(f"Исход '{key}' не подходит под категорию {category.name}")
								pass

if __name__ == "__main__":
	generate_categories()

	statistic_live = Statistic("ws://188.253.24.91:7300/output", True)
	statistic_prematch = Statistic("ws://188.253.24.91:7301/output", False)

	async def main():
		live_task = asyncio.create_task(statistic_live.run())
		prematch_task = asyncio.create_task(statistic_prematch.run())

		await asyncio.gather(prematch_task)

	asyncio.run(main())