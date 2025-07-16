import { JSX, useMemo } from "react";
import styles from './EmployeeStatus.module.css';
import { EmployeeStatusProps } from "./EmployeeStatus.props";

export const EmployeeStatus = ({ employee }: EmployeeStatusProps): JSX.Element => {
    const activeFilters = useMemo(() => {
        // Добавляем проверки на существование свойств
        if (!employee.filters || !employee.filters.bookmakers) {
            return [];
        }

        return employee.filters.bookmakers.flatMap(bookmaker => {
            const filters = [];

            if (bookmaker?.live?.filter) {
                filters.push({
                    type: `${bookmaker.name || 'Unknown'} (live)`,
                    sports: bookmaker.live.sports || []
                });
            }

            if (bookmaker?.prematch?.filter) {
                filters.push({
                    type: `${bookmaker.name || 'Unknown'} (prematch)`,
                    sports: bookmaker.prematch.sports || []
                });
            }

            return filters;
        });
    }, [employee.filters]);

    return (
        <div className={styles.wrapper}>
            <div className={styles.mainInfo}>
                <h3 className={styles.h3}>{employee.name}</h3>

                <div className={styles.group}>{employee.group}</div>

                <div className={styles.statusContainer}>
                    <div className={`${styles.status} ${employee.works ? styles.statusWorking : styles.statusNotWorking}`}>
                        {employee.works ? "Работает" : "Не работает"}
                    </div>
                </div>
            </div>

            {activeFilters.length > 0 && employee.works && (
                <div className={styles.filtersContainer}>
                    {activeFilters.map((filter, index) => (
                        <div key={`${filter.type}-${index}`} className={styles.filterItem}>
                            <span className={styles.filterType}>{filter.type}</span>
                            {filter.sports.length > 0 && (
                                <div className={styles.sportsList}>
                                    {filter.sports.map((sport) => (
                                        <span key={sport} className={styles.sportTag}>
                                            {sport}
                                        </span>
                                    ))}
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};