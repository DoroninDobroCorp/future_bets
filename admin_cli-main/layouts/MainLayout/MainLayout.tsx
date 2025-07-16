import styles from './MainLayout.module.css';
import { FunctionComponent, JSX, useEffect } from "react";
import { LayoutProps } from "./MainLayout.props";
import { Menu } from '@/components/Menu/Menu';
import { SportStore } from "../../stores/SportStore";
import { toast } from 'react-toastify';
import { useRouter } from 'next/router';
import { BookmakerStore}  from '../../stores/BookmakerStore';
import { useUnit } from 'effector-react';
import $sport from "../../stores/SportStore";
import $bookmaker from "../../stores/BookmakerStore";

const Layout = ({ children }: LayoutProps): JSX.Element => {
	const bookmakers = useUnit($bookmaker)
	const sports = useUnit($sport)

	const router = useRouter()

	const fetchSports = async () => {
        const [, status] = await SportStore.getSports()
		if (status == 200) {
			toast.success("Виды спорта получены")
		} else {
			toast.error("Виды спорта не получены")
		}
    }

	const fetchBookmakers = async () => {
        const [, status] = await BookmakerStore.getBookmakers()
		if (status == 200) {
			toast.success("Букмекеры получены")
		} else {
			toast.error("Букмекеры не получены")
		}
    }

	useEffect(() => {
		if (!sports) {
			fetchSports()
		}
		if (!bookmakers) {
			fetchBookmakers()
		}
	}, [router]) // eslint-disable-line react-hooks/exhaustive-deps

	return (
		<>
			<div className={styles.app}>
				<Menu />
				<div className={styles.body}>
					{children}
				</div>
			</div>
		</>
	);
};

export const withLayout = <T extends Record<string, unknown>>(Component: FunctionComponent<T>) => {
	return function withLayoutComponent(props: T): JSX.Element {
		return (
			<>
				<Layout>
					<Component {...props} />
				</Layout>
			</>
		);
	};
};