extra_percents = [
	{'min': 2.29, 'max': 2.75, 'extra_percent': 1.03},
	{'min': 2.75, 'max': 3.2, 'extra_percent': 1.04},
	{'min': 3.2, 'max': 3.7, 'extra_percent': 1.05}
]

# Получение дополнительного процента
def get_extra_percent(pinnacle_odd: float) -> float:
	for ep in extra_percents:
		if ep['min'] <= pinnacle_odd < ep['max']:
			return ep['extra_percent']
	return 1.0

# Расчет ROI
def calculate_roi(
		sansa_odd: float,
		pinnacle_odd: float,
		margin: float,
		market_type: int,
		second_bookmaker_name: str,
		sport_name: str
) -> float:
	extra_percent = get_extra_percent(pinnacle_odd)

	if second_bookmaker_name == "LOBBET":
		if sport_name == "TENNIS":
			return (sansa_odd / (pinnacle_odd * margin * extra_percent) - 1 - 0.03) * 100 * 0.67
		if market_type == 0:
			return (sansa_odd / (pinnacle_odd * margin * extra_percent) - 1 - 0.03) * 100 * 0.67
		if market_type < 0:
			return (sansa_odd / (pinnacle_odd * margin * extra_percent) - 1 - 0.015) * 100 * 0.75
		return (sansa_odd / (pinnacle_odd * margin * extra_percent) - 1 - 0.03) * 100 * 0.67

	elif second_bookmaker_name == "LADBROKES":
		if sport_name == "TENNIS":
			return (sansa_odd / (pinnacle_odd * margin * extra_percent) - 1 - 0.02) * 100 * 0.75
		if market_type == 0:
			return (sansa_odd / (pinnacle_odd * margin * extra_percent) - 1 - 0.02) * 100 * 0.75
		if market_type < 0:
			return (sansa_odd / (pinnacle_odd * margin * extra_percent) - 1) * 100 * 0.85
		return (sansa_odd / (pinnacle_odd * margin * extra_percent) - 1 - 0.02) * 100 * 0.75

	# Default case
	return (sansa_odd / (pinnacle_odd * margin * extra_percent) - 1 - 0.03) * 100 * 0.67