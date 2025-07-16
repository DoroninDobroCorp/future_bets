import $calculator_api from "@/http/CalcApi";
import $manager_api from "@/http/ManagerApi";
import { ICalcSumBetWithUsers } from "@/interfaces/Bet.interface";
import { IPairOneOutcome } from "@/interfaces/Pair";
import { AxiosResponse } from "axios";

export default class LogBetService {
    static async LogBet(pair: IPairOneOutcome, bet: ICalcSumBetWithUsers, userId: number): Promise<AxiosResponse> {
        return $manager_api.post(`/log_bet`, { pair, bet, userId }, { withCredentials: false });
    }

    static async GetCalcSumBet(userId: string, pair: IPairOneOutcome): Promise<AxiosResponse<ICalcSumBetWithUsers>> {
        return $calculator_api.post<ICalcSumBetWithUsers>(`/calc-bet`, { userId, pair }, { withCredentials: false });
    }
}