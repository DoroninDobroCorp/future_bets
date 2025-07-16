import { JSX, useEffect, useState } from "react";
import styles from "./EmployeesPageComponent.module.css";
import { useUnit } from "effector-react";
import { $employees, EmployeesStore, setEmployees } from "@/stores/EmployeesStore";
import { EmployeeStatus } from "@/components/EmployeeStatus/EmployeeStatus";

export const EmployeesPageComponent = (): JSX.Element => {
    const employees = useUnit($employees);
    const [lastUpdated, setLastUpdated] = useState<string>("");
    const [isLoading, setIsLoading] = useState<boolean>(false);

    useEffect(() => {
        let isMounted = true;

        const updateData = async () => {
            if (isLoading) return;

            setIsLoading(true);
            try {
                const data = await EmployeesStore.getEmployees();
                if (isMounted) {
                    setEmployees(data);
                    setLastUpdated(new Date().toLocaleTimeString());
                }
            } catch (error) {
                console.error("Error fetching employees:", error);
            } finally {
                if (isMounted) setIsLoading(false);
            }
        };

        updateData();

        const interval = setInterval(updateData, 15000);

        return () => {
            isMounted = false;
            clearInterval(interval);
        };
    }, []);

    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h1 className={styles.title}>Сотрудники</h1>
                {lastUpdated && (
                    <div className={styles.refreshInfo}>Обновлено: {lastUpdated}</div>
                )}
            </div>

            <div className={styles.employeesList}>
                {employees ? (
                    employees.map(employee => (
                        <EmployeeStatus key={`${employee.name}-${employee.group}`} employee={employee} />
                    ))
                ) : (
                    <div className={styles.loading}>Загрузка данных..</div>
                )}
            </div>
        </div>
    );
};