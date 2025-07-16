import $runner_api from "@/http/RunnerApi";
import { IParseStatus } from "@/interfaces/ParseStatus.interface";
import { AxiosResponse } from "axios";

export default class ParseStatusService {
    static async getStatus(): Promise<AxiosResponse<IParseStatus[]>> {
        return $runner_api.get<IParseStatus[]>(`/status`, { withCredentials: false });
    }

    static async setCommand(bookmaker: string, run: boolean): Promise<AxiosResponse> {
        return $runner_api.post(`/set-command`, null, { withCredentials: false, params: { bookmaker, run } });
    }
}