export interface IGeneralFilter {
    bookmakers: IBookmakerFilter[]
}

export interface IBookmakerFilter {
    name: string
    live: ILive
    prematch: ILive
}

export interface ILive {
    filter: boolean
    sports: string[]
}

