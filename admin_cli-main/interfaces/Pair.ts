import { IOutcome } from "./Outcome"

export interface IMatchToPair {
	bookmaker: string
	leagueName: string
	homeScore: number
	awayScore: number
	homeName: string
	awayName: string
	matchId: string
	createdAt: Date
}

export interface IPair {
	first: IMatchToPair
	second: IMatchToPair
	outcome: IOutcome[]
	isLive: boolean
	sportName: string
	createdAt: Date
}

export interface IPairOneOutcome {
	first: IMatchToPair
	second: IMatchToPair
	outcome: IOutcome
	isLive: boolean
	sportName: string
	createdAt: Date
}