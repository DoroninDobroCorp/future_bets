import { JSX } from 'react';
import styles from './List.module.css';
import { ListLeagueProps } from './List.props';
import cn from 'classnames';

export const ListLeague = ({ elements, onSelect, size, right }: ListLeagueProps): JSX.Element => {
	const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
		if (onSelect) {
			onSelect(event.target.value)
		}
	};

	elements.sort((x, y) => x.leagueName.localeCompare(y.leagueName));

	return (
		<div className={styles.wrapper}>
			<select className={cn(styles.select, right && styles.right)} size={size} onChange={handleChange} >
				{elements && elements.map((item) => (
					<option key={item.id} value={item.id}>
						{item.leagueName}
					</option>
				))}
			</select>
		</div>

	);
};