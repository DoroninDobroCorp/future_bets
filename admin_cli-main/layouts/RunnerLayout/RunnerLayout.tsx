import { FunctionComponent, JSX } from 'react';
import styles from './RunnerLayout.module.css';
import { RunnerLayoutProps } from "./RunnerLayout.props";

const RunnerLayout = ({ children }: RunnerLayoutProps): JSX.Element => {

    return (
        <>
            <div className={styles.app}>
                <div className={styles.body}>
                    {children}
                </div>
            </div>
        </>
    );
};

export const withRunnerLayout = <T extends Record<string, unknown>>(Component: FunctionComponent<T>) => {
    return function withLayoutComponent(props: T): JSX.Element {
        return (
            <>
                <RunnerLayout>
                    <Component {...props} />
                </RunnerLayout>
            </>
        );
    };
};