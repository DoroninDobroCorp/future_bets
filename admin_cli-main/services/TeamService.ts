import $auto_matcher_api from "@/http/AutoMatcherApi";
import { IUnMatchTeamResponse } from "@/interfaces/response/team.interface";
import { AxiosResponse } from "axios";

export default class TeamService {
    static async getOnlineUnMatchTeam(sportName: string, firstBookmakerName: string, secondBookmakerName: string): Promise<AxiosResponse<IUnMatchTeamResponse>> {
        return $auto_matcher_api.get<IUnMatchTeamResponse>(`/hand-merge/teams/online-unmatch`, { withCredentials: false, params: { sportName, firstBookmakerName, secondBookmakerName } });
    }

    static async getOnlineUnMatchTeamPrematch(sportName: string, firstBookmakerName: string, secondBookmakerName: string): Promise<AxiosResponse<IUnMatchTeamResponse>> {
        return $auto_matcher_api.get<IUnMatchTeamResponse>(`/hand-merge/teams/online-unmatch-prematch`, { withCredentials: false, params: { sportName, firstBookmakerName, secondBookmakerName } });
    }

    static async getUnMatchTeam(sportName: string, firstBookmakerName: string, secondBookmakerName: string): Promise<AxiosResponse<IUnMatchTeamResponse>> {
        return $auto_matcher_api.get<IUnMatchTeamResponse>(`/hand-merge/teams/get-unmatch`, { withCredentials: false, params: { sportName, firstBookmakerName, secondBookmakerName } });
    }

    static async addTeamPair(strFirstTeamID: string, strSecondTeamID: string): Promise<AxiosResponse> {
        const firstTeamID = parseInt(strFirstTeamID, 10)
        const secondTeamID = parseInt(strSecondTeamID, 10)
        return $auto_matcher_api.post(`/hand-merge/teams/create-pair`, null, { withCredentials: false, params: { firstTeamID, secondTeamID } })
    }
}