import {createEvent, createStore} from "effector";
import EmployeesService from "@/services/EmployeesService";
import {IGeneralFilter} from "@/interfaces/GeneralFilter";
import {IEmployeeStatus} from "@/interfaces/EmployeeStatus.interface";

export const toggleWorkMode = createEvent()
export const $workMode = createStore<boolean>(false)
    .on(toggleWorkMode, (state) => !state)

export const setStartedWork = createEvent<number>()
export const $startedWork = createStore<number>(0)
    .on(setStartedWork, (_, time) => time)

export const setEmployees = createEvent<IEmployeeStatus[]>()
export const $employees = createStore<IEmployeeStatus[]>([])
    .on(setEmployees, (_, employees) =>  employees)

export class EmployeesStore {
    static async sendWorkTime(userId: number, time: number) {
        return await EmployeesService.sendWorkTime(userId, time)
            .then((res) => {
                return res.status;
            })
            .catch((err) => {
                return err.response;
            })
    }

    static async sendFilters(userId: number, filters: IGeneralFilter) {
        return await EmployeesService.sendFilters(userId, filters)
            .then((res) => {
                return res.status;
            })
            .catch((err) => {
                return err.response;
            })
    }

    static async sendToken(token: string) {
        return await EmployeesService.sendToken(token)
            .then((res) => {
                return res.data;
            })
            .catch((err) => {
                return err.response;
            })
    }

    static async getEmployees() {
        return await EmployeesService.getEmployees()
            .then((res) => {
                return res.data
            })
            .catch((err) => {
                return err.response
            })
    }
}