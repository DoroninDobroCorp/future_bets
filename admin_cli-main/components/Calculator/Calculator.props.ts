
import { IPairOneOutcome } from "@/interfaces/Pair";
import { DetailedHTMLProps, HTMLAttributes } from "react";

export interface CalculatorProps extends DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement> {
    pair: IPairOneOutcome
    userId: number
    setClose: (close: boolean) => void;
    setSentPair: (sent: IPairOneOutcome) => void;
}