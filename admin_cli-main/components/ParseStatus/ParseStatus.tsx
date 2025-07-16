import { JSX } from "react";
import styles from './ParseStatus.module.css';
import { ParseStatusProps } from "./ParseStatus.props";
import { Status } from "../Status/Status";
import { StatusError, StatusOK } from "@/helpers/TypeStatusByTime";
import cn from 'classnames';
import { ParseStatusStore } from "@/stores/ParseStatusStore";
import { toast } from "react-toastify";

export const ParseStatus = ({ bookmaker }: ParseStatusProps): JSX.Element => {

    const date = new Date(bookmaker.createdAt)

    const fetchON = async () => {
        const status = await ParseStatusStore.setCommand(bookmaker.name, true)
        if (status == 200) {
            toast.success("Парсер запускается... Ждите!")
        } else {
            toast.error("Ошибка запуска")
        }
    }

    const fetchOFF = async () => {
        const status = await ParseStatusStore.setCommand(bookmaker.name, false)
        if (status == 200) {
            toast.success("Парсер выключается... Ждите!")
        } else {
            toast.error("Ошибка запуска")
        }
    }

    return (
        <>
            <div className={styles.wrapper}>
                <h3 className={styles.h3}>{bookmaker.name}</h3>

                <div className={styles.status}>
                    <Status status={bookmaker.status == "ON" ? StatusOK : StatusError} title="" />
                </div>
                

                <h4 className={styles.h4}>{date.getMinutes()}:{date.getSeconds()}</h4>

                {bookmaker.status == "ON" ?
                    <div className={styles.buttonContainer} onClick={fetchOFF}>
                        <button className={cn(styles.button, styles.off)}>
                            Выключить
                        </button>
                    </div>
                    :
                    <div className={styles.buttonContainer}>
                        <button className={cn(styles.button)} onClick={fetchON}>
                            Включить
                        </button>
                    </div>
                }
            </div>
        </>
    )
};