import { JSX, useState } from "react";
import { MatchProps } from "./Match.props";
import styles from './Match.module.css';
import { Outcome } from "../Outcome/Outcome";
import cn from 'classnames';
import { IOutcome } from "@/interfaces/Outcome";
import { IPair } from "@/interfaces/Pair";
import Image from "next/image";
import { setFilterBet } from "../../stores/FilterBetStore";

export const Match = ({ keyMatch, pair, setSelectMatch }: MatchProps): JSX.Element => {
    const [arrow, setArrow] = useState<boolean>(false);

    const handleArrow = () => {
        setArrow(!arrow);
    };

    const handleOutcome = (out: IOutcome) => {
        const selectPair: IPair = pair
        const selectOutcome: IOutcome[] = [out]
        selectPair.outcome = selectOutcome
        setSelectMatch(selectPair, keyMatch)
    };

    const handleClose = () => {
        setFilterBet(pair.first.matchId)
    }

    return (
        <>
            <div className={cn(styles.wrapper, !pair.isLive && styles.prematch)}>
                <h3 className={styles.h3}>
                    {pair.sportName == 'Soccer' && <Image src="/football.png" width={20} height={20} alt="soccer" />}
                    {pair.sportName == 'Tennis' && <Image src="/tennis.png" width={20} height={20} alt="tennis" />}
                <span className={styles.span}>{pair.first.bookmaker} - {pair.second.bookmaker}</span></h3>
                <h3 className={styles.h3}><span className={styles.span}>{pair.first.homeName} vs {pair.first.awayName} | {pair.first.homeScore} - {pair.first.awayScore} | {!pair.isLive && "PREMATCH"}</span></h3>
                <h4 className={styles.h4}>{pair.second.homeName} : {pair.second.awayName} - ({pair.sportName})</h4>
                <p className={styles.p}>{pair.second.leagueName}</p>
                <div className={cn(styles.arrowIcon, arrow && styles.open)} onClick={handleArrow}>
                    <span className={styles.leftBar}></span>
                    <span className={styles.rightBar}></span>
                </div>
                <div className={cn(styles.outcomes)}>
                    {pair && pair.outcome.map((out, i) => {
                        return <Outcome key={out.outcome} className={cn(!arrow && i != 0 && styles.hide)}
                            outcome={out} bookmaker1={pair.first.bookmaker} bookmaker2={pair.second.bookmaker} setSelectOutcome={handleOutcome} />
                    })}
                </div>

                <div className={styles.closeBtn} onClick={handleClose}></div>
            </div>
        </>
    )
};