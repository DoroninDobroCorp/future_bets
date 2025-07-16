import $auto_matcher_api from "@/http/AutoMatcherApi";
import { IBookmakersResponse } from "@/interfaces/response/bookmaker.interface";
import { AxiosResponse } from "axios";

export default class BookmakerService {
    static async getBookmakers(): Promise<AxiosResponse<IBookmakersResponse>> {
        return $auto_matcher_api.get<IBookmakersResponse>(`/hand-merge/bookmakers`, { withCredentials: false });
    }
}