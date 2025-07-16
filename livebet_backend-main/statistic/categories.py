from copy import deepcopy
import json
import time


class Outcome:
    def __init__(self, key, pair, outcome) -> None:
        self.key = key
        self.pair = pair
        self.pair["outcome"] = outcome

        self.addedAt = time.time()
        self.updatedAt = None

    def update(self, outcome):
        self.updatedAt = time.time()
        self.pair["outcome"] = outcome


class Category:
    def __init__(self, ID, name, roiLow, roiHigh, coefLow, coefHigh) -> None:
        self.ID = ID
        self.name = name
        
        self.blocked_matches = []
        self.outcomes = {}

        self.roiLow = roiLow
        self.roiHigh = roiHigh
        self.coefLow = coefLow
        self.coefHigh = coefHigh

    def check(self, outcome):
        if self.roiLow <= outcome["roi"] < self.roiHigh and self.coefLow <= outcome["score2"]["value"] < self.coefHigh:
            return True
        return False

    def has(self, key):
        return key in self.outcomes

    def get(self, key):
        return self.outcomes[key]

    def add(self, key, pair, outcome):
        if not self.check(outcome):
            return False

        self.outcomes[key] = Outcome(
            key=key,
            pair=pair,
            outcome=outcome
        )
        self.outcomes[key].addedAt = time.time()
        self.outcomes[key].updatedAt = time.time()

        return True

    def update(self, key, outcome):
        if not self.check(outcome):
            self.delete(key)
            return False

        self.outcomes[key].update(outcome)
        return True

    def block(self, event):
        if event["matchId"] not in self.blocked_matches:
            self.blocked_matches.append(event["matchId"])

    def unblock(self, event):
        if event["matchId"] in self.blocked_matches:
            self.blocked_matches.remove(event["matchId"])

    def is_blocked(self, event):
        return event["matchId"] in self.blocked_matches

    def delete(self, key):
        self.outcomes.pop(key)


categories = {
    "Soccer": {
        "Lobbet": {},
        "Ladbrokes": {},
        "StarCasino": {},
        "Unibet": {}
    },
    "Tennis": {
        "Lobbet": {},
        "Ladbrokes": {},
        "StarCasino": {},
        "Unibet": {}
    },
    "Basketball": {
        "Lobbet": {},
        "Ladbrokes": {},
        "StarCasino": {},
        "Unibet": {}
    }
}

sport_categories = [
    "Soccer",
    "Tennis",
    "Basketball"
]

bookmaker_categories = [
    "Lobbet",
    "Ladbrokes",
    "StarCasino",
    "Unibet"
]

roi_categories = [
    [-2, 0],
    [0, 3],
    [3, 6],
    [6, 10],
    [10, 15],
    [15, 1000] # 1000 - max
]

coef_categories = [
    [-1, 1.8], # -1 - min
    [1.8, 2.6],
    [2.6, 1000] # 1000 - max
]

def generate_categories():
    json_categories = deepcopy(categories)

    count = 0
    for sport in sport_categories:
        for bookmaker in bookmaker_categories:
            for roi_category in roi_categories:
                for coef_category in coef_categories:
                    category = Category(
                        ID=count+100,
                        name=f"{count+100}|{sport}|{bookmaker}|{roi_category[0]}_{roi_category[1]}|{coef_category[0]}_{coef_category[1]}",
                        roiLow=roi_category[0],
                        roiHigh=roi_category[1],
                        coefLow=coef_category[0],
                        coefHigh=coef_category[1]
                    )
                    categories[sport][bookmaker][category.name] = category
                    json_categories[sport][bookmaker][category.name] = category.ID
                    count += 1

    with open("categories.json", "w") as f:
        json.dump(json_categories, f, indent=4)
