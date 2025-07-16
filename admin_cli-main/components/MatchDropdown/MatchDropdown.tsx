import { JSX, useEffect, useState } from "react";
import { MatchDropdownProps } from "./MatchDropdown.props";
import styles from './MatchDropdown.module.css';
import { Dropdown } from "../Dropdown/Dropdown";
import { IMatchDropdown } from "@/interfaces/MatchDropdown.interface";
import { LS_bookmaker1, LS_bookmaker2, LS_sport } from "@/helpers/LocalStorageItemKey";

export const MatchDropdown = ({ bookmakers, sports, setSelect }: MatchDropdownProps): JSX.Element => {
    const [selectedBookmaker1, setSelectedBookmaker1] = useState<string>("");
    const [selectedBookmaker2, setSelectedBookmaker2] = useState<string>("");
    const [selectedSport, setSelectedSport] = useState<string>("");

    const [bookmakers1, setbookmakers1] = useState<string[]>([...bookmakers]);
    const [bookmakers2, setbookmakers2] = useState<string[]>([...bookmakers]);

    const handleBookmaker1Select = (selected: string) => {
        setSelectedBookmaker1(selected)

        setbookmakers2(bookmakers.filter((number) => number !== selected))
    };

    const handleBookmaker2Select = (selected: string) => {
        setSelectedBookmaker2(selected)

        setbookmakers1(bookmakers.filter((number) => number !== selected))
    };

    const handleSportSelect = (selected: string) => {
        setSelectedSport(selected)
    };

    useEffect(() => {
        if (selectedBookmaker1 != "" && selectedBookmaker2 != "" && selectedSport != "") {
            localStorage.setItem(LS_bookmaker1, selectedBookmaker1)
            localStorage.setItem(LS_bookmaker2, selectedBookmaker2)
            localStorage.setItem(LS_sport, selectedSport)
        }
    }, [selectedBookmaker1, selectedBookmaker2, selectedSport])

    useEffect(() => {
        if (bookmakers) {     
            const bookmaker1 = localStorage.getItem(LS_bookmaker1)
            const bookmaker2 = localStorage.getItem(LS_bookmaker2)
            const sport = localStorage.getItem(LS_sport)
            console.log(bookmaker1)
            if (bookmaker1) {
                setSelectedBookmaker1(bookmaker1)
            }
            if (bookmaker2) {
                setSelectedBookmaker2(bookmaker2)
            } 
            if (sport) {
                setSelectedSport(sport)
            } 
        }
    }, [bookmakers])

    useEffect(() => {
        if (selectedBookmaker1 != "" && selectedBookmaker2 != "" && selectedSport != "") {
            const select: IMatchDropdown = {
                bookmaker1: selectedBookmaker1,
                bookmaker2: selectedBookmaker2,
                sport: selectedSport,
            }
            setSelect(select)
        }
    }, [selectedBookmaker1, selectedBookmaker2, selectedSport, bookmakers1, bookmakers2]) // eslint-disable-line react-hooks/exhaustive-deps

    return (
        <div className={styles.wrapper}>
            <Dropdown elements={bookmakers1} selected={selectedBookmaker1} onSelect={handleBookmaker1Select} />
            <div className={styles.sports}>
                <Dropdown elements={sports} selected={selectedSport} onSelect={handleSportSelect} />
            </div>
            <Dropdown elements={bookmakers2} selected={selectedBookmaker2} onSelect={handleBookmaker2Select} />
        </div>
    );
};