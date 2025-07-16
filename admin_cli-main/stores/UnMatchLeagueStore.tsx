import { ILeague } from "@/interfaces/League.interface"
import LeaguesService from "@/services/LeaguesService"
import { createEvent, createStore } from "effector"

const $unMatchLeague = createStore<ILeague[] | null>(null)
export const setLeague = createEvent<ILeague[] | null>()
export const deleteLeague = createEvent<string>()

export default $unMatchLeague
    .on(setLeague, (leagues: ILeague[] | null, newLeagues: ILeague[] | null) => {
        const sortLeagues = newLeagues?.filter(leg => {
            const name = leg.leagueName.split(" ")
            if (name[name.length - 1] != "corners") {
                return leg
            }
        })
        return sortLeagues
    })
    .on(deleteLeague, (leagues: ILeague[] |  null, deleteLeague: string) => {
        return leagues?.filter(key => key.id.toString() != deleteLeague)
    })

export class UnMatchLeagueStore {
    static async getLeagues(sport: string, bookmaker1: string, bookmaker2: string): Promise<[ILeague[] | null, number]> {
        return await LeaguesService.getUnMatchLeague(sport, bookmaker1, bookmaker2)
            .then((res) => {
                setLeague(res.data.data)
                return [res.data.data, res.status] as [ILeague[], number]
            })
            .catch((err) => {
                return [null, err.response.status]
            })
    }

    static async addLeaguePair(league1: string, league2: string): Promise<number> {
        return await LeaguesService.addLeaguePair(league1, league2)
        .then((res) => {
            if (res.status == 200) {
                deleteLeague(league1)
                deleteLeague(league2)
            }
            return res.status as number
        })
        .catch((err) => {
            return err.response.status
        })
    }
}