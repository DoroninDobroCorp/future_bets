import { JSX, useEffect, useState } from "react";
import styles from './ClientMatchFilter.module.css';
import { ClientMatchFilterProps } from "./ClientMatchFilter.props";
import { IBookmakerFilter, IGeneralFilter } from "@/interfaces/GeneralFilter";
import cn from 'classnames';
import { LS_filter } from "@/helpers/LocalStorageItemKey";

export const ClientMatchFilter = ({ bookmakers, sports, setSelect }: ClientMatchFilterProps): JSX.Element => {
    const [filter, setFilter] = useState<Map<string, IBookmakerFilter>>(new Map());

    const [visible, setVisible] = useState<boolean>(true);

    const handleBookmaker = (bookmaker: string, checked: boolean) => {
        if (checked) {
            setFilter(oldFilter => {
                const newFilter = new Map(oldFilter)
                const oldBookmakerFilter: IBookmakerFilter | undefined = filter.get(bookmaker)
                let bookmakerFilter: IBookmakerFilter = { live: { filter: false, sports: [] }, prematch: { filter: false, sports: [] }, name: bookmaker }
                if (oldBookmakerFilter) { bookmakerFilter = oldBookmakerFilter }
                newFilter.set(bookmaker, bookmakerFilter)
                return newFilter
            })
        } else {
            setFilter(oldFilter => {
                const newFilter = new Map(oldFilter)
                newFilter.delete(bookmaker)
                return newFilter
            })
        }
    }

    const handleLive = (bookmaker: string, live: boolean, checked: boolean) => {
        if (checked) {
            setFilter(oldFilter => {
                const newFilter = new Map(oldFilter)
                const bookmakerFilter: IBookmakerFilter | undefined = filter.get(bookmaker)
                if (bookmakerFilter) {
                    if (live) { bookmakerFilter.live.filter = true } else { bookmakerFilter.prematch.filter = true }
                    newFilter.set(bookmaker, bookmakerFilter)
                }
                return newFilter
            })
        } else {
            setFilter(oldFilter => {
                const newFilter = new Map(oldFilter)
                const bookmakerFilter: IBookmakerFilter | undefined = filter.get(bookmaker)
                if (bookmakerFilter) {
                    if (live) {
                        bookmakerFilter.live.filter = false
                        bookmakerFilter.live.sports = []
                    } else {
                        bookmakerFilter.prematch.filter = false
                        bookmakerFilter.prematch.sports = []
                    }
                    newFilter.set(bookmaker, bookmakerFilter)
                }
                return newFilter
            })
        }
    }

    const handleSports = (bookmaker: string, live: boolean, sport: string, checked: boolean) => {
        if (checked) {
            setFilter(oldFilter => {
                const newFilter = new Map(oldFilter)
                const bookmakerFilter: IBookmakerFilter | undefined = filter.get(bookmaker)
                if (bookmakerFilter) {
                    if (live) { bookmakerFilter.live.sports.push(sport) } else { bookmakerFilter.prematch.sports.push(sport) }
                    newFilter.set(bookmaker, bookmakerFilter)
                }
                return newFilter
            })
        } else {
            setFilter(oldFilter => {
                const newFilter = new Map(oldFilter)
                const bookmakerFilter: IBookmakerFilter | undefined = filter.get(bookmaker)
                if (bookmakerFilter) {
                    if (live) { bookmakerFilter.live.sports = bookmakerFilter.live.sports.filter(s => s !== sport) } else { bookmakerFilter.prematch.sports = bookmakerFilter.prematch.sports.filter(s => s !== sport) }
                    newFilter.set(bookmaker, bookmakerFilter)
                }
                return newFilter
            })
        }
    }

    useEffect(() => {
        const rawSavedFilter = localStorage.getItem(LS_filter)
        if (rawSavedFilter) {
            const newMap: Map<string, IBookmakerFilter> = new Map(JSON.parse(rawSavedFilter))
            setFilter(newMap)
        }
    }, [])

    useEffect(() => {
        localStorage.setItem(LS_filter, JSON.stringify(Array.from(filter)))
    }, [filter])

    useEffect(() => {
        const selectFilter: IGeneralFilter = { bookmakers: [] }
        filter.forEach((value) => {
            selectFilter.bookmakers.push(value)
        })
        setSelect(selectFilter)
    }, [filter])  // eslint-disable-line react-hooks/exhaustive-deps

    return (
        <div className={styles.wrapperMain}>
            <div className={cn(styles.wrapper, !visible && styles.wrapperHide)}>
                {bookmakers && bookmakers.sort().map(bookmaker => {
                    return (<div key={bookmaker}>
                        <div className={styles.checkBoxGroup}>
                            <div className={styles.checkboxWrapper}>
                                <label><input className={styles.input} type="checkbox" checked={(filter && filter.get(bookmaker)) ? true : false} onChange={(event) => handleBookmaker(bookmaker, event.target.checked)} />
                                    <span className={styles.checkbox}></span><b>{bookmaker}</b></label>
                            </div>
                        </div>
                        {filter && filter.get(bookmaker) && <div className={styles.checkBoxGroup}>
                            <div className={styles.checkboxWrapper}>
                                <label><input className={styles.input} type="checkbox" checked={(filter.get(bookmaker)?.live.filter == true) ? true : false} onChange={(event) => handleLive(bookmaker, true, event.target.checked)} />
                                    <span className={styles.checkbox}></span><b>Live</b></label>

                            </div>
                            {filter && filter.get(bookmaker)?.live.filter && sports.map(sport => {
                                return <div key={sport} className={styles.checkboxWrapper}>
                                    <label key={sport}><input className={styles.input} type="checkbox" checked={filter.get(bookmaker)?.live.sports.includes(sport) ? true : false} onChange={(event) => handleSports(bookmaker, true, sport, event.target.checked)} />
                                        <span className={styles.checkbox}></span><p>{sport}</p></label>
                                </div>
                            })}
                        </div>}
                        {filter && filter.get(bookmaker) && <div className={styles.checkBoxGroup}>
                            <div className={styles.checkboxWrapper}>
                                <label><input className={styles.input} type="checkbox" checked={(filter.get(bookmaker)?.prematch.filter == true) ? true : false} onChange={(event) => handleLive(bookmaker, false, event.target.checked)} />
                                    <span className={styles.checkbox}></span><b>Prematch</b></label>
                            </div>
                            {filter && filter.get(bookmaker)?.prematch.filter && sports.map(sport => {
                                return <div key={sport} className={styles.checkboxWrapper}>
                                    <label><input className={styles.input} type="checkbox" checked={filter.get(bookmaker)?.prematch.sports.includes(sport) ? true : false} onChange={(event) => handleSports(bookmaker, false, sport, event.target.checked)} />
                                        <span className={styles.checkbox}></span><p>{sport}</p></label>
                                </div>
                            })}
                        </div>}
                    </div>)
                })}

            </div>
            <div className={cn(styles.arrowIcon, visible && styles.arrowIconTransform)} onClick={() => setVisible(!visible)}>
                <span className={styles.leftBar}></span>
                <span className={styles.rightBar}></span>
            </div>
        </div>

    );
};