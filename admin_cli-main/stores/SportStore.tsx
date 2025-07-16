import SportService from "@/services/SportService";
import { createEvent, createStore } from "effector";

const $sport = createStore<string[] | null>(null)
export const setSports = createEvent<string[] | null>()

export default $sport
    .on(setSports, (sports: string[] | null, newSports: string[] | null) => {
        return newSports
    })

export class SportStore {
    static async getSports(): Promise<[string[] | null, number]> {
        return await SportService.getSports()
            .then((res) => {
                setSports(res.data.data)
                return [res.data.data, res.status] as [string[], number]
            })
            .catch((err) => {
                return [null, err.response.status]
            })
    }
}