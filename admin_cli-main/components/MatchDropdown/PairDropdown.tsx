import { JSX, useEffect, useState } from "react";
import { PairDropdownProps } from "./MatchDropdown.props";
import styles from './MatchDropdown.module.css';
import { IPairDropdown } from "@/interfaces/MatchDropdown.interface";
import { Dropdown } from "../Dropdown/Dropdown";

export const PairDropdown = ({ bookmakers, selected, setSelect }: PairDropdownProps): JSX.Element => {
    const [selectedBookmaker1, setSelectedBookmaker1] = useState<string>("");
    const [selectedBookmaker2, setSelectedBookmaker2] = useState<string>("");

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

    useEffect(() => {
        if (selected && selected.bookmaker1) {
            setSelectedBookmaker1(selected.bookmaker1)
        }
        if (selected && selected.bookmaker2) {
            setSelectedBookmaker2(selected.bookmaker2)
        }
    }, [selected?.bookmaker1, selected?.bookmaker2]) // eslint-disable-line react-hooks/exhaustive-deps

    useEffect(() => {
        const select: IPairDropdown = {}
        if (selectedBookmaker1 != "") {
            select.bookmaker1 = selectedBookmaker1

        }
        if (selectedBookmaker2 != "") {
            select.bookmaker2 = selectedBookmaker2

        }
        setSelect(select)
    }, [selectedBookmaker1, selectedBookmaker2, bookmakers1, bookmakers2]) // eslint-disable-line react-hooks/exhaustive-deps

    return (
        <div className={styles.wrapper}>
            <Dropdown elements={bookmakers1} selected={selected?.bookmaker1} onSelect={handleBookmaker1Select} />
            <Dropdown elements={bookmakers2} selected={selected?.bookmaker2} onSelect={handleBookmaker2Select} />
        </div>
    );
};