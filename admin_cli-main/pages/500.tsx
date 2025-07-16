import Link from 'next/link';
import styles from '../styles/errors-page.module.css';
import { JSX } from 'react';

export function Error500(): JSX.Element {
    return (
      <div className={styles.wrapper}>
        <Link href="/"></Link>
        <h3>Internal Server Error (status code 500)</h3>
      </div>
    )
}
  
export default Error500;