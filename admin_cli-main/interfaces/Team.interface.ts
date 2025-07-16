export interface IUnMatchTeam {
    leagueIDFirst: number
    leagueIDSecond: number
    bookmakerNameFirst: string
    bookmakerNameSecond: string
    leagueNameFirst: string
    leagueNameSecond: string
    teamsFirst: ITeam[]
    teamsSecond: ITeam[]
    sportName: string
}

export interface ITeam {
    teamID: number
    teamName: string
}