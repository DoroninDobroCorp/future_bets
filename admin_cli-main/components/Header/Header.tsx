import { JSX } from 'react';
import styles from './Header.module.css';
import { HeaderProps } from './Header.props';

export const Header = ({header}: HeaderProps): JSX.Element => {
	return (
		<>
			<div className={styles.header}>{header}</div>
		</>
	);
};