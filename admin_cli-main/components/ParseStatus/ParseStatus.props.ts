import { IParseStatus } from "@/interfaces/ParseStatus.interface";
import { DetailedHTMLProps, HTMLAttributes } from "react";

export interface ParseStatusProps extends DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement> {
    bookmaker: IParseStatus
}