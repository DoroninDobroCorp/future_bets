import { IGeneralFilter } from "@/interfaces/GeneralFilter";
import { DetailedHTMLProps, HTMLAttributes } from "react";

export interface ClientMatchFilterProps extends DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement> {
    bookmakers: string[];
    sports: string[];
    setSelect: (selectedItems: IGeneralFilter) => void;
}