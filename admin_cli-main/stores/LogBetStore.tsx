import { ICalcSumBetWithUsers } from "@/interfaces/Bet.interface";
import { IPairOneOutcome } from "@/interfaces/Pair";
import LogBetService from "@/services/LogBetService";
import { createEvent, createStore } from "effector";

const $bet = createStore<ICalcSumBetWithUsers | null>(null)
export const setBet = createEvent<ICalcSumBetWithUsers | null>()

export default $bet
    .on(setBet, (bet: ICalcSumBetWithUsers | null, newBet: ICalcSumBetWithUsers | null) => {
        return newBet
    })

export class LogBetStore {
    static async GetCalcSumBet(userId: string, pair: IPairOneOutcome): Promise<[ICalcSumBetWithUsers | null, number]> {
        return await LogBetService.GetCalcSumBet(userId, pair)
            .then((res) => {
                setBet(res.data)
                return [res.data, res.status] as [ICalcSumBetWithUsers, number]
            })
            .catch((err) => {
                return [null, err.response]
            })
    }

    static async LogBet(pair: IPairOneOutcome, bet: ICalcSumBetWithUsers, userId: number): Promise<number> {
        return await LogBetService.LogBet(pair, bet, userId)
            .then((res) => {
                return res.status
            })
            .catch((err) => {
                return err.response
            })
    }
}
