import cn from 'classnames';
import styles from './Status.module.css';
import { JSX } from 'react';
import { StatusProps } from './Status.props';
import { StatusError, StatusErrorPrematch, StatusOK, StatusOKPrematch, StatusWarn, StatusWarnPrematch } from '@/helpers/TypeStatusByTime';

export const Status = ({ status, title }: StatusProps): JSX.Element => {
    return (
        <>
            <div className={styles.wrapper}>
                <div className={cn(styles.status,
                    status == StatusOK && styles.ok,
                    status == StatusWarn && styles.warn,
                    status == StatusError && styles.error,
                    status == StatusOKPrematch && styles.ok,
                    status == StatusWarnPrematch && styles.warn,
                    status == StatusErrorPrematch && styles.error,
                )}></div>
                <h4 className={cn(styles.title)}>{title}</h4>
            </div>
        </>
    )
};