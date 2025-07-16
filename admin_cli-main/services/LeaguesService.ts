import $auto_matcher_api from "@/http/AutoMatcherApi";
import { ILeagueResponse } from "@/interfaces/response/league.interface";
import { AxiosResponse } from "axios";

export default class LeaguesService {
    static async getOnlineUnMatchLeagues(sportName: string, firstBookmakerName: string, secondBookmakerName: string): Promise<AxiosResponse<ILeagueResponse>> {
        return $auto_matcher_api.get<ILeagueResponse>(`/hand-merge/leagues/online-unmatch`, { withCredentials: false, params: { sportName, firstBookmakerName, secondBookmakerName } });
    }

    static async getOnlineUnMatchLeaguesPrematch(sportName: string, firstBookmakerName: string, secondBookmakerName: string): Promise<AxiosResponse<ILeagueResponse>> {
        return $auto_matcher_api.get<ILeagueResponse>(`/hand-merge/leagues/online-unmatch-prematch`, { withCredentials: false, params: { sportName, firstBookmakerName, secondBookmakerName } });
    }

    static async getAllLeagues(sportName: string, firstBookmakerName: string, secondBookmakerName: string): Promise<AxiosResponse<ILeagueResponse>> {
        return $auto_matcher_api.get<ILeagueResponse>(`/hand-merge/leagues/`, { withCredentials: false, params: { sportName, firstBookmakerName, secondBookmakerName } });
    }

    static async getUnMatchLeague(sportName: string, firstBookmakerName: string, secondBookmakerName: string): Promise<AxiosResponse<ILeagueResponse>> {
        return $auto_matcher_api.get<ILeagueResponse>(`/hand-merge/leagues/get-unmatch`, { withCredentials: false, params: { sportName, firstBookmakerName, secondBookmakerName } });
    }

    static async addLeaguePair(strFirstLeagueID: string, strSecondLeagueID: string): Promise<AxiosResponse> {
        const firstLeagueID = parseInt(strFirstLeagueID, 10)
        const secondLeagueID = parseInt(strSecondLeagueID, 10)
        return $auto_matcher_api.post(`/hand-merge/leagues/create-pair`, null, { withCredentials: false, params: { firstLeagueID, secondLeagueID } })
    }
}