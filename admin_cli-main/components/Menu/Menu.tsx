import { JSX, useState } from 'react';
import styles from './Menu.module.css';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

export const Menu = (): JSX.Element => {
    const pathname = usePathname();
    const [isCollapsed, setIsCollapsed] = useState(false);

    const menuItems = [
        { href: '/online-leagues', label: 'Онлайн лиги' },
        { href: '/online-matches', label: 'Онлайн матчи' },
        { href: '/online-leagues-prematch', label: 'Онлайн лиги prematch' },
        { href: '/online-matches-prematch', label: 'Онлайн матчи prematch' },
        { href: '/leagues', label: 'Лиги' },
        { href: '/matches', label: 'Матчи' },
        { href: '/all-leagues', label: 'Все лиги' },
        // { href: '/unpaired-matches', label: 'Матчи без пары' },
        { href: '/match-candidates', label: 'Кандидаты матчей' },
        { href: '/league-candidates', label: 'Кандидаты лиг' },
        { href: '/paired-leagues', label: 'Сопоставленные лиги' },
        { href: '/paired-matches', label: 'Сопоставленные матчи' },
    ];

    function getMenuItem(idx: number): JSX.Element {
        return (
            <li key={menuItems[idx].href} className={styles.menuItem}>
                <Link href={menuItems[idx].href} className={`${styles.menuLink} ${pathname === menuItems[idx].href ? styles.active : ''}`}>
                    {menuItems[idx].label}
                </Link>
            </li>
        );
    }

    return (
        <>
            <nav className={`${styles.menu} ${isCollapsed ? styles.menuCollapsed : ''}`}>
                <ul className={styles.menuList}>
                    {getMenuItem(0)}
                    {getMenuItem(1)}
                    {getMenuItem(2)}
                    {getMenuItem(3)}
                    {getMenuItem(4)}
                    {getMenuItem(5)}
                    {getMenuItem(6)}
                </ul>
                <hr className={styles.line} />
                <ul className={styles.menuList}>
                    {getMenuItem(4)}
                    {getMenuItem(5)}
                </ul>
            </nav>
            <button
                className={`${styles.toggleButton} ${isCollapsed ? styles.toggleButtonCollapsed : ''}`}
                onClick={() => setIsCollapsed(!isCollapsed)}
            >
                {isCollapsed ? '>' : '<'}
            </button>
        </>
    );
};