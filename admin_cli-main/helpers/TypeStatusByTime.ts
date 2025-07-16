export type Status = {
    name: string
    from: number
    to: number
}

export const StatusOK: Status = { name: "OK", from: 0, to: 10 }
export const StatusWarn: Status = { name: "WARN", from: 10, to: 11 }
export const StatusError: Status = { name: "ERROR", from: 11, to: 1000 }

export const StatusOKPrematch: Status = { name: "OK", from: 0, to: 30 }
export const StatusWarnPrematch: Status = { name: "WARN", from: 30, to: 35 }
export const StatusErrorPrematch: Status = { name: "ERROR", from: 35, to: 1000 }

export function GetStatus(sub: number): Status {
    if (sub >= StatusOK.from && sub < StatusOK.to) {
        return StatusOK
    } else if (sub >= StatusWarn.from && sub < StatusWarn.to) {
        return StatusWarn
    }
    return StatusError
}

export function GetStatusPrematch(sub: number): Status {
    if (sub >= StatusOKPrematch.from && sub < StatusOKPrematch.to) {
        return StatusOKPrematch
    } else if (sub >= StatusWarnPrematch.from && sub < StatusWarnPrematch.to) {
        return StatusWarnPrematch
    }
    return StatusErrorPrematch
}