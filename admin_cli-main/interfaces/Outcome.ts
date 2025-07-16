export interface IOutcome {
    outcome: string
    roi: number
    margin: number
    score1: IOdd
    score2: IOdd
    marketType: number
}

export interface IOdd {
    value: number
}