import Link from 'next/link';
import styles from '../styles/errors-page.module.css';
import { JSX } from 'react';

export function Error404(): JSX.Element {
    return (
      <div className={styles.wrapper}>
        <Link href="/">
        </Link>
        <h3>Страница не найдена. Ошибка 404</h3>
        <Link href="/">Вернутся на главную</Link>
      </div>
    )
}
  
export default Error404;