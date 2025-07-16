import { JSX, useEffect, useState } from "react";
import styles from "./MatchesPageComponent.module.css";
import { Header } from "@/components/Header/Header";
import { ListMatchLeague } from "@/components/List/ListMatchLeague";
import { useUnit } from "effector-react";
import $sport from "../../stores/SportStore";
import $bookmaker from "../../stores/BookmakerStore";
import $unMatchTeam from "../../stores/UnMatchTeamStore"
import { IMatchDropdown } from "@/interfaces/MatchDropdown.interface";
import { MatchDropdown } from "@/components/MatchDropdown/MatchDropdown";
import { UnMatchTeamStore } from "../../stores/UnMatchTeamStore";
import { toast } from "react-toastify";
import { ListMatchTeam } from "@/components/List/ListMatchTeam";
import { ITeam } from "@/interfaces/Team.interface";

export const MatchesPageComponent = ({ }): JSX.Element => {
    const bookmakers = useUnit($bookmaker)
    const sports = useUnit($sport)
    const unMatchTeam = useUnit($unMatchTeam)

    const [selectDropdowns, setSelectDropdowns] = useState<IMatchDropdown>();
    const [selectedLeague, setSelectedLeague] = useState<string | null>(null)

    const [teams1, setTeams1] = useState<ITeam[] | null>(null);
    const [teams2, setTeams2] = useState<ITeam[] | null>(null);

    const handleLeagueSelect = (selected: string) => {
        setSelectedLeague(selected);
    };

    const [selectedTeam1, setSelectedTeam1] = useState<string | null>(null);
    const [selectedTeam2, setSelectedTeam2] = useState<string | null>(null);

    const handleTeamSelect1 = (selected: string) => {
        setSelectedTeam1(selected);
    };

    const handleTeamSelect2 = (selected: string) => {
        setSelectedTeam2(selected);
    };

    useEffect(() => {
        const fetchUnMatchTeams = async () => {
            if (selectDropdowns) {
                const [, status] = await UnMatchTeamStore.getUnMatchTeams(selectDropdowns.sport, selectDropdowns.bookmaker1, selectDropdowns.bookmaker2)
                if (status == 200) {
                    toast.success("Команды получены")
                } else if (status == 404) {
                    toast.warning("Нет несопоставленных команд")
                } else {
                    toast.error("Команды не получены")
                }
            }
        }

        fetchUnMatchTeams()
    }, [selectDropdowns])

    useEffect(() => {
        if (selectedLeague && unMatchTeam) {
            const ids = selectedLeague.split("-")
            const matchLeagues = unMatchTeam.filter(key => {
                return key.leagueIDFirst.toString() == ids[0] && key.leagueIDSecond.toString() == ids[1]
            })
            if (matchLeagues[0]) {
                setTeams1(matchLeagues[0].teamsFirst)
                setTeams2(matchLeagues[0].teamsSecond)
            }
        }

        if (selectedLeague == null) {
            setTeams1(null)
            setTeams2(null)
        }
    }, [selectedLeague, unMatchTeam])

    const addTeamPairs = async () => {
        if (selectedTeam1 && selectedTeam1 != "" && selectedTeam2 && selectedTeam2 != "") {
            const status = await UnMatchTeamStore.addLeaguePair(selectedTeam1, selectedTeam2)
            if (status == 200) {
                toast.success("Пара команд создана")
                setSelectedTeam1(null)
                setSelectedTeam2(null)
                // TODO: fast variant for hack
                if (teams1?.length == 1 || teams2?.length == 1) {
                    setSelectedLeague(null)
                }
            } else if (status == 409) {
                toast.warn("Пара уже была создана")
                setSelectedTeam1(null)
                setSelectedTeam2(null)
            } else {
                toast.error("Пара команд не создана")
            }
        }
    }

    return (
        <>
            <div className={styles.wrapper}>
                <Header header="Создание пары команд" />

                {sports && bookmakers && <MatchDropdown bookmakers={bookmakers} sports={sports} setSelect={setSelectDropdowns} />}

                <div className={styles.lists}>
                    {selectDropdowns && unMatchTeam && <ListMatchLeague elements={unMatchTeam} onSelect={handleLeagueSelect} size={10} center={true} />}
                </div>

                <div className={styles.lists}>
                    {teams1 && selectedLeague && <ListMatchTeam elements={teams1} onSelect={handleTeamSelect1} size={10} right={true} />}
                    {teams2 && selectedLeague && <ListMatchTeam elements={teams2} onSelect={handleTeamSelect2} size={10} />}
                </div>

                <div className={styles.buttonContainer}>
                    <button className={styles.button} onClick={addTeamPairs}>
                        Создать пару команд
                    </button>
                </div>
            </div>

        </>
    );
};