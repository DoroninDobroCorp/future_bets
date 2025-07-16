export interface DropdownProps {
    elements: string[];
    selected?: string;
    onSelect: (selectedItem: string) => void;
}