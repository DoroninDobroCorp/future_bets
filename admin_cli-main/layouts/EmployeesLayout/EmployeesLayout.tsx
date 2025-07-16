import { FunctionComponent, JSX } from 'react';
import styles from './EmployeesLayout.module.css';
import { EmployeesLayoutProps } from "./EmployeesLayout.props";

const EmployeesLayout = ({ children }: EmployeesLayoutProps): JSX.Element => {

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

export const withEmployeesLayout = <T extends Record<string, unknown>>(Component: FunctionComponent<T>) => {
    return function withLayoutComponent(props: T): JSX.Element {
        return (
            <>
                <EmployeesLayout>
                    <Component {...props} />
                </EmployeesLayout>
            </>
        );
    };
};