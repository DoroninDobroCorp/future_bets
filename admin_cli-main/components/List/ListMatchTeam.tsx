import { JSX } from "react";
import { ListMatchTeamProps } from "./List.props";
import styles from './List.module.css';
import cn from 'classnames';

export const ListMatchTeam = ({ elements, onSelect, size, right, center }: ListMatchTeamProps): JSX.Element => {
	const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
		if (onSelect) {
			onSelect(event.target.value)
		}
	};

	elements.sort((x, y) => x.teamName.localeCompare(y.teamName));

	return (
		<div className={styles.wrapper}>
			<select className={cn(styles.select, right && styles.right, center && styles.center)} size={size} onChange={handleChange} >
				{elements && elements.map((item) => (
					<option key={item.teamID} value={item.teamID}>
						{item.teamName}
					</option>
				))}
			</select>
		</div>

	);
};