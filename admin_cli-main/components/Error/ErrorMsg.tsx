import { JSX } from "react";
import styles from "./ErrorMsg.module.css";
import { ErrorMsgProps } from "./ErrorMsg.props";


export const ErrorMsg = ({message}: ErrorMsgProps): JSX.Element => {
    return (
        <div className={styles.errorMessage}>
			{message}
		</div>
    );
};