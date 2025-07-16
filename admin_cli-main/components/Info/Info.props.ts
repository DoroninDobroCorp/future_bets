import { DetailedHTMLProps, HTMLAttributes } from "react";

export interface InfoProps extends DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement> {
	text: string;
}