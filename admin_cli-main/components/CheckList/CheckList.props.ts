import { Candidate } from "@/interfaces/Candidate";

export interface CheckListProps {
    label1: string;
    label2: string;
    second_label: string;
    elements: Candidate[];
    onSelect?: (selectedItems: Candidate[]) => void;
}