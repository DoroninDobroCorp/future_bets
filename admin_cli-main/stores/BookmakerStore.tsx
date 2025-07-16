import BookmakerService from "@/services/BookmakerService"
import { createEvent, createStore } from "effector"

const $bookmaker = createStore<string[] | null>(null)
export const setBookmakers = createEvent<string[] | null>()

export default $bookmaker
    .on(setBookmakers, (bookmakers: string[] | null, newBookmakers: string[] | null) => {
        return newBookmakers
    })

export class BookmakerStore {
    static async getBookmakers(): Promise<[string[] | null, number]> {
        return await BookmakerService.getBookmakers()
            .then((res) => {
                setBookmakers(res.data.data)
                return [res.data.data, res.status] as [string[], number]
            })
            .catch((err) => {
                return [null, err.response.status]
            })
    }
}