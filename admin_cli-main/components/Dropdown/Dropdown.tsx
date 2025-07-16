import { JSX } from 'react';
import styles from './Dropdown.module.css';
import { DropdownProps } from './Dropdown.props';

export const Dropdown = ({ elements, selected, onSelect }: DropdownProps): JSX.Element => {
	const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
		if (onSelect) {
			onSelect(event.target.value);
		}
	};

	return (
			<select className={styles.select} onChange={handleChange} value={selected && selected}>
				<option value=""></option>
				{elements.map((element) => (
					<option key={element} value={element}>
						{element}
					</option>
				))}
			</select>
	);
};