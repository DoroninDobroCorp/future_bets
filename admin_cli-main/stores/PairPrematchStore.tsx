import { IPair } from "@/interfaces/Pair"
import { createEvent, createStore } from "effector"

const $pairPrematch = createStore<Map<string, IPair>>(new Map)
const $pairPrematchCount = createStore<number>(0)

export const setPairPrematch = createEvent<IPair[]>()
export const setPairPrematchCount = createEvent<number>()

$pairPrematch.on(setPairPrematch, (pairs: Map<string, IPair>, newPairs: IPair[] | null) => {
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

$pairPrematchCount.on(setPairPrematchCount, (count: number, newCount: number) => {
    return newCount
})

export { $pairPrematch, $pairPrematchCount }

export class PairPrematchStore { }
