import {IGeneralFilter} from "@/interfaces/GeneralFilter";

export interface IEmployeeStatus {
    name: string
    group: string
    works: boolean
    filters: IGeneralFilter
}