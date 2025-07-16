import { BookmakerStore } from "@/stores/BookmakerStore";
import { useRouter } from "next/router";
import { FunctionComponent, JSX, useEffect } from "react";
import { toast } from "react-toastify";
import $bookmaker from "../../stores/BookmakerStore";
import $sport from "../../stores/SportStore";
import { toggleWorkMode, $workMode, setStartedWork, $startedWork } from "@/stores/EmployeesStore";
import { useUnit } from "effector-react";
import styles from './ClientLayout.module.css';
import { ClientLayoutProps } from "./ClientLayout.props";
import { SportStore } from "@/stores/SportStore";
import { EmployeesStore } from "@/stores/EmployeesStore";

const ClientLayout = ({ children, userId }: ClientLayoutProps): JSX.Element => {
    const router = useRouter()
    const bookmakers = useUnit($bookmaker)
    const sports = useUnit($sport)
    const workMode = useUnit($workMode)
    const startedWork = useUnit($startedWork)

    const fetchBookmakers = async () => {
        const [, status] = await BookmakerStore.getBookmakers()
        if (status == 200) {
            toast.success("Букмекеры получены")
        } else {
            toast.error("Букмекеры не получены")
        }
    }

    const fetchSports = async () => {
        const [, status] = await SportStore.getSports()
        if (status == 200) {
            toast.success("Виды спорта получены")
        } else {
            toast.error("Виды спорта не получены")
        }
    }

    useEffect(() => {
        if (!bookmakers) {
            fetchBookmakers()
        }
        if (!sports) {
            fetchSports()
        }
    }, [router]) // eslint-disable-line react-hooks/exhaustive-deps

    // work mode logic
    const changeWorkMode = async (workMode: boolean, startedWork: number) => {
        toggleWorkMode();
        const now = new Date().getTime() / 1000;

        if (!workMode) {
            setStartedWork(now)
            await EmployeesStore.sendWorkTime(userId, 0)
        } else {
            const workDuration = now - startedWork;
            await EmployeesStore.sendWorkTime(userId, Math.round(workDuration))
        }
    }

    return (
        <>
            <div className={styles.app}>
                <button className={styles.workModeButton} onClick={() => changeWorkMode(workMode, startedWork)}>
                    { workMode ? "Не работаю" : "Работаю" }
                </button>
                <div className={styles.body}>
                    {children}
                </div>
            </div>
        </>
    );
};

export const withClientLayout = <T extends Record<string, unknown> & { userId: number }>(Component: FunctionComponent<T>) => {
    return function withLayoutComponent(props: T): JSX.Element {
        return (
            <>
                <ClientLayout userId={props.userId}>
                    <Component {...props} />
                </ClientLayout>
            </>
        );
    };
};