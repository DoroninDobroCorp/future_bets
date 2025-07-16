import { IOutcome } from "@/interfaces/Outcome";
import { DetailedHTMLProps, HTMLAttributes } from "react";

export interface OutcomeProps extends DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement> {
    outcome: IOutcome
    bookmaker1: string
    bookmaker2: string
    className?: string
    setSelectOutcome: (selectedItem: IOutcome) => void;
}