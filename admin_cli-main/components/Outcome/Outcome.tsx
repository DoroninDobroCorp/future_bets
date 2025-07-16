import { JSX } from "react";
import { OutcomeProps } from "./Outcome.props";
import styles from './Outcome.module.css';
import cn from 'classnames';

export const Outcome = ({ outcome, bookmaker1, bookmaker2, className, setSelectOutcome }: OutcomeProps): JSX.Element => {
    const handleClick = () => {
        setSelectOutcome(outcome)
    };

    return (
        <>
            <div className={cn(styles.wrapper, className)} onClick={handleClick}>
                <p className={cn(styles.p, outcome.roi >= 12 && styles.superOK, outcome.roi < 12 && outcome.roi >= 3 && styles.ok, outcome.roi < 3 && outcome.roi >= 0 && styles.warn)}><strong>{outcome.outcome}</strong> | {bookmaker1}: {outcome.score1.value} | {bookmaker2}: {outcome.score2.value} | ROI: {outcome.roi.toFixed(4)}% | MARGIN: {outcome.margin.toFixed(4)}</p>
            </div>
        </>
    )
};