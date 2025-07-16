import { IParseStatus } from "@/interfaces/ParseStatus.interface"
import ParseStatusService from "@/services/ParseStatusService"
import { createEvent, createStore } from "effector"

const $parseStatus = createStore<IParseStatus[] | null>(null)
export const setParseStatus = createEvent<IParseStatus[] | null>()

export default $parseStatus
    .on(setParseStatus, (status: IParseStatus[] | null, newStatus: IParseStatus[] | null) => {
        return newStatus
    })

export class ParseStatusStore {
    static async getStatus(): Promise<[IParseStatus[] | null, number]> {
        return await ParseStatusService.getStatus()
            .then((res) => {
                setParseStatus(res.data)
                return [res.data, res.status] as [IParseStatus[], number]
            })
            .catch((err) => {
                return [null, err.response]
            })
    }

    static async setCommand(bookmaker: string, run: boolean): Promise<number> {
        return await ParseStatusService.setCommand(bookmaker, run)
            .then((res) => {
                return res.status
            })
            .catch((err) => {
                return err.response
            })
    }
}