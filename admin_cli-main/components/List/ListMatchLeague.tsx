import { JSX } from "react";
import { ListMatchLeagueProps } from "./List.props";
import styles from './List.module.css';
import cn from 'classnames';

export const ListMatchLeague = ({ elements, onSelect, size, right, center }: ListMatchLeagueProps): JSX.Element => {
	const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
		if (onSelect) {
			onSelect(event.target.value)
		}
	};

	elements.sort((x, y) => x.leagueNameFirst.localeCompare(y.leagueNameFirst));

	return (
		<div className={styles.wrapper}>
			<select className={cn(styles.select, right && styles.right, center && styles.center)} size={size} onChange={handleChange} >
				{elements && elements.map((item) => (
					<option key={item.leagueIDFirst + "-" + item.leagueIDSecond} value={item.leagueIDFirst + "-" + item.leagueIDSecond}>
						{item.leagueNameFirst + "  -  " + item.leagueNameSecond}
					</option>
				))}
			</select>
		</div>

	);
};