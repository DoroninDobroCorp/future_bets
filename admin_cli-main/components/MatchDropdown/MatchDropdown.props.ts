import { IMatchDropdown, IPairDropdown } from "@/interfaces/MatchDropdown.interface";
import { DetailedHTMLProps, HTMLAttributes } from "react";

export interface MatchDropdownProps extends DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement> {
    bookmakers: string[];
    sports: string[];
    setSelect: (selectedItems: IMatchDropdown) => void;
}

export interface PairDropdownProps extends DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement> {
    bookmakers: string[];
    selected?: IPairDropdown;
    setSelect: (selectedItems: IPairDropdown) => void;
}