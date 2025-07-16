import { DetailedHTMLProps, HTMLAttributes } from "react";
import {IEmployeeStatus} from "@/interfaces/EmployeeStatus.interface";

export interface EmployeeStatusProps extends DetailedHTMLProps<HTMLAttributes<HTMLDivElement>, HTMLDivElement> {
    employee: IEmployeeStatus
}