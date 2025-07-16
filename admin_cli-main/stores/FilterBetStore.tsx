import { createEvent, createStore } from "effector"

const $filterBet = createStore<string[]>([])
export const setFilterBet = createEvent<string>()

export default $filterBet
    .on(setFilterBet, (ids: string[], newID: string) => {
        ids?.push(newID)
        return ids
    })