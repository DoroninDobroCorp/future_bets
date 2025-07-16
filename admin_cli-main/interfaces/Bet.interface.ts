export interface ICalculatedBet {
    originalAmount: number;
    adjustedAmount: number;
    percentage: number;
}

export interface ICalcSumBetWithUsers {
    usersCount: number
    calcBet: ICalculatedBet
}