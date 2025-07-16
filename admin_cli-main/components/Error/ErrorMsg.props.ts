import { DetailedHTMLProps, HTMLAttributes } from "react";

export interface ErrorMsgProps extends DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement> {
    message: string;
}