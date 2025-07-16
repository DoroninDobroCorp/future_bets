import { ILeague } from "@/interfaces/League.interface";
import { ITeam, IUnMatchTeam } from "@/interfaces/Team.interface";

export interface ListLeagueProps {
    elements: ILeague[];
    right?: boolean;
    size: number;
    onSelect: (selectedItem: string) => void;
}

export interface ListMatchLeagueProps {
    elements: IUnMatchTeam[];
    right?: boolean;
    center?: boolean;
    size: number;
    onSelect: (selectedItem: string) => void;
}

export interface ListMatchTeamProps {
    elements: ITeam[];
    right?: boolean;
    center?: boolean;
    size: number;
    onSelect: (selectedItem: string) => void;
}