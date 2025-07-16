import { ILeague } from "@/interfaces/League.interface"
import LeaguesService from "@/services/LeaguesService"
import { createEvent, createStore } from "effector"

const $allLeagues = createStore<ILeague[] | null>(null)
export const setLeagues = createEvent<ILeague[] | null>()

export default $allLeagues
    .on(setLeagues, (leagues: ILeague[] | null, newLeagues: ILeague[] | null) => {
        const sortLeagues = newLeagues?.filter(leg => {
            const name = leg.leagueName.split(" ")
            if (name[name.length - 1] != "corners") {
                return leg
            }
        })
        return sortLeagues
    })

export class AllLeagueStore {
    static async getLeagues(sport: string, bookmaker1: string, bookmaker2: string): Promise<[ILeague[] | null, number]> {
        return await LeaguesService.getAllLeagues(sport, bookmaker1, bookmaker2)
            .then((res) => {
                setLeagues(res.data.data)
                return [res.data.data, res.status] as [ILeague[], number]
            })
            .catch((err) => {
                return [null, err.response.status]
            })
    }

    static async addLeaguePair(league1: string, league2: string): Promise<number> {
        return await LeaguesService.addLeaguePair(league1, league2)
            .then((res) => {
                return res.status as number
            })
            .catch((err) => {
                return err.response.status
            })
    }
}