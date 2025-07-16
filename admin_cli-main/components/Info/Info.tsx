import { JSX } from "react"
import styles from "./Info.module.css"
import { InfoProps } from "./Info.props"

export const Info = ({text}: InfoProps): JSX.Element => {
	return (
		<>
			<div className={styles.info}>{text}</div>
		</>
	);
};
