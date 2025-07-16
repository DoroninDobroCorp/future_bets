import { IPair } from "@/interfaces/Pair";
import { DetailedHTMLProps, HTMLAttributes } from "react";

export interface MatchProps extends DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement> {
    pair: IPair
    keyMatch: string
    setSelectMatch: (selectedItem: IPair, key: string) => void;
}