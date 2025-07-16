import { IPair } from "@/interfaces/Pair"
import { createEvent, createStore } from "effector"

const $pair = createStore<Map<string, IPair>>(new Map)
const $pairCount = createStore<number>(0)

export const setPair = createEvent<IPair[]>()
export const setPairCount = createEvent<number>()

$pair.on(setPair, (pairs: Map<string, IPair>, newPairs: IPair[] | null) => {
    const local = new Map<string, IPair>()
    if (newPairs) {
        const filtredPairs = newPairs.filter(pair => {
            const outcomes = pair.outcome.filter(outcome => outcome.roi > -30)
            return outcomes.length > 0
        })
        filtredPairs.map(pair => local.set(pair.first.bookmaker + pair.first.matchId + pair.second.bookmaker + pair.second.matchId + pair.sportName, pair))
    }
    return local
})

$pairCount.on(setPairCount, (count: number, newCount: number) => {
    return newCount
})

export { $pair, $pairCount }

export class PairStore { }

