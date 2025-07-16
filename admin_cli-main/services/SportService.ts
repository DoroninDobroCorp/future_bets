import $auto_matcher_api from "@/http/AutoMatcherApi";
import { ISportsResponse } from "@/interfaces/response/sports.interface";
import { AxiosResponse } from "axios";

export default class SportService {
    static async getSports(): Promise<AxiosResponse<ISportsResponse>> {
        return $auto_matcher_api.get<ISportsResponse>(`/hand-merge/sports`, { withCredentials: false });
    }
}