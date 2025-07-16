import { JSX, useEffect } from "react";
import styles from "./RunnerPageComponent.module.css";
import { ParseStatusStore } from "@/stores/ParseStatusStore";
import $parseStatus from "../../stores/ParseStatusStore";
import { useUnit } from "effector-react";
import { ParseStatus } from "@/components/ParseStatus/ParseStatus";

export const RunnerPageComponent = ({ }): JSX.Element => {
    const parseStatus = useUnit($parseStatus)

    useEffect(() => {
        const interval = setInterval(async () => {
            await ParseStatusStore.getStatus()
        }, 5000)
        return () => {
            clearInterval(interval)
        }
    }, [parseStatus])

    if (parseStatus) {
        parseStatus.sort((x, y) => x.name.localeCompare(y.name));
    }
    
    return (
        <>
            <div className={styles.wrapper}>
                {parseStatus ? parseStatus?.map(status => {
                    return <ParseStatus key={status.name} bookmaker={status} />
                }) 
                : <h1 className={styles.h1}>ЖДИТЕ!</h1>}
            </div>
        </>
    )
}