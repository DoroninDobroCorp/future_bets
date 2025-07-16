import { JSX } from 'react';
import styles from './Text.module.css';
import { TextProps } from './Text.props';

export const Text = ({ text }: TextProps): JSX.Element => {
    return (
        <div className={styles.text}>
            {text}
        </div>
    );
};