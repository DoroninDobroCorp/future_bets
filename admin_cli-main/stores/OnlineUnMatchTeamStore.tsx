import { ITeam, IUnMatchTeam } from "@/interfaces/Team.interface"
import TeamService from "@/services/TeamService"
import { createEvent, createStore } from "effector"

const $onlineUnMatchTeam = createStore<IUnMatchTeam[] | null>(null)
export const setUnMatchTeam = createEvent<IUnMatchTeam[] | null>()
export const deleteTeam = createEvent<string>()

export default $onlineUnMatchTeam
    .on(setUnMatchTeam, (teams: IUnMatchTeam[] | null, newTeams: IUnMatchTeam[] | null) => {
        return newTeams
    })
    .on(deleteTeam, (teams: IUnMatchTeam[] |  null, deleteTeam: string) => {
        let teamsFirst: ITeam[]
        let teamsSecond: ITeam[]
        let newTeams: IUnMatchTeam[] | undefined
        newTeams = teams?.map(key => {
            teamsFirst = key.teamsFirst.filter(teams => {
                return teams.teamID.toString() != deleteTeam
            })
            teamsSecond = key.teamsSecond.filter(teams => {
                return teams.teamID.toString() != deleteTeam
            })
            key.teamsFirst = teamsFirst
            key.teamsSecond = teamsSecond
            return key
        })

        if (newTeams) {
            newTeams = newTeams?.filter(key => {
                return key.teamsFirst.length > 0 && key.teamsSecond.length > 0
            })
        }

        return newTeams
    })

export class OnlineUnMatchTeamStore {
    static async getUnMatchTeams(sport: string, bookmaker1: string, bookmaker2: string): Promise<[IUnMatchTeam[] | null, number]> {
        return await TeamService.getOnlineUnMatchTeam(sport, bookmaker1, bookmaker2)
            .then((res) => {
                setUnMatchTeam(res.data.data)
                return [res.data.data, res.status] as [IUnMatchTeam[], number]
            })
            .catch((err) => {
                return [null, err.response.status]
            })
    }

    static async addLeaguePair(team1: string, team2: string): Promise<number> {
        return await TeamService.addTeamPair(team1, team2)
            .then((res) => {
                if (res.status == 200) {
                    deleteTeam(team1)
                    deleteTeam(team2)
                }
                return res.status as number
            })
            .catch((err) => {
                return err.response.status
            })
    }

}