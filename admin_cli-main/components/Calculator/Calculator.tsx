import { JSX, useEffect, useState } from "react";
import { CalculatorProps } from "./Calculator.props";
import styles from './Calculator.module.css';
import { Status } from "../Status/Status";
import { GetStatus, GetStatusPrematch } from "@/helpers/TypeStatusByTime";
import { LogBetStore } from "@/stores/LogBetStore";
import { toast } from "react-toastify";
import $bet from "../../stores/LogBetStore";
import { useUnit } from "effector-react";
import cn from 'classnames';

export const Calculator = ({ pair, userId, setClose, setSentPair }: CalculatorProps): JSX.Element => {
    const bet = useUnit($bet)
    const [bookTimeout1, setBookTimeout1] = useState<number>(10);
    const [bookTimeout2, setBookTimeout2] = useState<number>(10);

    const handleClose = () => {
        setClose(false)
    }
    const handleSent = () => {
        setSentPair(pair)
    }

    const fetchBetLog = async () => {
        if (pair && bet) {
            const status = await LogBetStore.LogBet(pair, bet, userId)
            if (status == 200) {
                toast.success("Ставка залогирована")
            } else {
                toast.error("Не получилось залогировать ставки")
            }
            setTimeout(() => {
                handleSent()
            }, 5000)
        }
    }

    useEffect(() => {
        const interval = setInterval(async () => {
            await LogBetStore.GetCalcSumBet(userId.toString(), pair)
        }, 2000)
        return () => {
            clearInterval(interval)
        }
    }, [pair.outcome.outcome, pair.first.homeName]) // eslint-disable-line react-hooks/exhaustive-deps

    useEffect(() => {
        if (pair) {
            const interval = setInterval(() => {
                setBookTimeout1(new Date().getTime() / 1000 - new Date(pair.first.createdAt).getTime() / 1000)
                setBookTimeout2(new Date().getTime() / 1000 - new Date(pair.second.createdAt).getTime() / 1000)

            }, 500)
            return () => {
                clearInterval(interval);
            };
        }
    }, [pair.createdAt]) // eslint-disable-line react-hooks/exhaustive-deps

    return (
        <>
            <div className={styles.wrapper}>
                <h3 className={styles.h3}>Калькулятор ставки</h3>

                <div className={styles.status}>
                    <Status title={pair.first.bookmaker + " " + bookTimeout1.toFixed(2) + "s"} status={pair.isLive ? GetStatus(bookTimeout1) : GetStatusPrematch(bookTimeout1)} />
                    <Status title={pair.second.bookmaker + " " + bookTimeout2.toFixed(2) + "s"} status={pair.isLive ? GetStatus(bookTimeout2) : GetStatusPrematch(bookTimeout2)} />
                </div>

                <div className={styles.info}>
                    <p className={styles.p}><strong>{pair.first.bookmaker}:</strong> {pair.first.leagueName} - <strong><span className={styles.firstTeam}>{pair.first.homeName}</span> vs <span className={styles.secondTeam}>{pair.first.awayName}</span> | {!pair.isLive && "PREMATCH"}</strong></p>
                    <p className={styles.p}><strong>{pair.second.bookmaker}:</strong> {pair.second.leagueName} - <strong><span className={styles.firstTeam}>{pair.second.homeName}</span> vs <span className={styles.secondTeam}>{pair.second.awayName}</span></strong></p>
                    <p className={styles.p}><strong>Счёт:</strong> {pair.second.homeScore} - {pair.second.awayScore}</p>
                    <p className={styles.p}><strong>Исход:</strong> {pair.outcome.outcome}</p>
                    <p className={styles.p}><strong>Цена {pair.first.bookmaker}:</strong> {pair.outcome.score1.value}</p>
                    <p className={styles.p}><strong>Цена {pair.second.bookmaker}:</strong> {pair.outcome.score2.value}</p>
                    <p className={styles.p}><strong>Тип рынка :</strong> {pair.outcome.marketType < 0 ? 'Падающий' : 'Базовый'}</p>
                    <p className={cn(styles.p, pair.outcome.roi >= 12 && styles.roiSuperOk, pair.outcome.roi < 12 && pair.outcome.roi >= 3 && styles.roiOk, pair.outcome.roi < 3 && pair.outcome.roi >= 0 && styles.roiWarn, pair.outcome.roi < 0 && styles.roiError)}><strong>ROI: {pair.outcome.roi.toFixed(3)}</strong> <span className={styles.fakeROI}>ROI(5/):  {(pair.outcome.roi - 5/pair.outcome.score1.value).toFixed(3)}</span></p>
                    <p className={styles.p}><strong>MARGIN:</strong> {pair.outcome.margin}</p>
                </div>

                {bet && <div className={styles.sum}>
                    <h4 className={styles.h4}>Рекомендуемая сумма: {bet.calcBet.adjustedAmount}</h4>
                    <h4 className={styles.h4}>Доступная сумма: {bet.calcBet.originalAmount}</h4>
                    <h4 className={styles.h4}>Оставшийся процент: {bet.calcBet.percentage.toFixed(2)}%</h4>
                    <h4 className={styles.h4}>Cмотрят матч: {bet.usersCount}</h4>
                </div>}

                <div className={styles.closeBtn} onClick={handleClose}></div>

                <div className={styles.buttonContainer}>
                    <button className={styles.button} onClick={fetchBetLog} style={{
                        display: pair.isLive ? (
                            GetStatus(bookTimeout1).name == "ERROR" || GetStatus(bookTimeout2).name == "ERROR" ? 'none' : 'flex'
                        ) : (
                            GetStatusPrematch(bookTimeout1).name == "ERROR" || GetStatusPrematch(bookTimeout2).name == "ERROR" ? 'none' : 'flex'
                        )
                    }}>

                        Отправить ставку
                    </button>
                </div>

            </div>
        </>
    )
};