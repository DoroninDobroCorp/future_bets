import $manager_api from "@/http/ManagerApi";
import {AxiosResponse} from "axios";
import {IGeneralFilter} from "@/interfaces/GeneralFilter";

export default class EmployeesService {
	static async sendWorkTime(userId: number, time: number): Promise<AxiosResponse> {
		return $manager_api.post(`/work_time`, { userId, time }, { withCredentials: false });
	}

	static async sendFilters(userId: number, filters: IGeneralFilter): Promise<AxiosResponse> {
		return $manager_api.post(`/filters`, { userId, filters }, { withCredentials: false });
	}

	static async sendToken(token: string): Promise<AxiosResponse> {
		return $manager_api.post(`/token`, { token }, { withCredentials: false })
	}

	static async getEmployees(): Promise<AxiosResponse> {
		return $manager_api.get("/employees")
	}
}