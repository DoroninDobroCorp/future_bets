import { Status } from "@/helpers/TypeStatusByTime";
import { DetailedHTMLProps, HTMLAttributes } from "react";

export interface StatusProps extends DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement> {
    status: Status;
    title: string;
}