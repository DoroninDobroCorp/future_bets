import { JSX, useEffect, useState } from "react";
import styles from "./OnlineLeaguesPrematchPageComponent.module.css";
import $sport from "../../stores/SportStore";
import $bookmaker from "../../stores/BookmakerStore";
import $onlineLeaguePrematch, { OnlineLeaguePrematchStore } from "../../stores/OnlineLeaguesPrematchStore";
import { Header } from "@/components/Header/Header";
import { ListLeague } from "@/components/List/ListLeague";
import { MatchDropdown } from "@/components/MatchDropdown/MatchDropdown";
import { IMatchDropdown } from "@/interfaces/MatchDropdown.interface";
import { useUnit } from "effector-react";
import { toast } from "react-toastify";

export const OnlineLeaguesPrematchPageComponent = ({ }): JSX.Element => {
	const bookmakers = useUnit($bookmaker)
	const sports = useUnit($sport)
	const unMatchLeague = useUnit($onlineLeaguePrematch)

	const [selectDropdowns, setSelectDropdowns] = useState<IMatchDropdown>();

	const [selectedLeague1, setSelectedLeague1] = useState<string | null>(null);
	const [selectedLeague2, setSelectedLeague2] = useState<string | null>(null);

	const handleLeagueSelect1 = (selected: string) => {
		setSelectedLeague1(selected);
	};

	const handleLeagueSelect2 = (selected: string) => {
		setSelectedLeague2(selected);
	};

	const addLeaguePairs = async () => {
		if (selectedLeague1 && selectedLeague1 != "" && selectedLeague2 && selectedLeague2 != "") {
			const status = await OnlineLeaguePrematchStore.addLeaguePair(selectedLeague1, selectedLeague2)
			if (status == 200) {
				toast.success("Пара создана")
				setSelectedLeague1(null)
				setSelectedLeague2(null)
			} else if (status == 409) {
				toast.warn("Пара уже была создана")
				setSelectedLeague1(null)
				setSelectedLeague2(null)
			} else {
				toast.error("Пара не создана")
			}
		}
	}

	useEffect(() => {
		const fetchUnMatchLeagues = async () => {
			if (selectDropdowns) {
				const [, status] = await OnlineLeaguePrematchStore.getLeagues(selectDropdowns.sport, selectDropdowns.bookmaker1, selectDropdowns.bookmaker2)
				if (status == 200) {
					toast.success("Лиги получены")
				} else if (status == 404) {
					toast.warning("Нет несопоставленных лиг")
				} else {
					toast.error("Лиги не получены")
				}
			}
		}

		fetchUnMatchLeagues()
	}, [selectDropdowns])

	return (
		<>
			<div className={styles.wrapper}>
				<Header header="Создание онлайн лиг PREMATCH" />

				{sports && bookmakers && <MatchDropdown bookmakers={bookmakers} sports={sports} setSelect={setSelectDropdowns} />}

				<div className={styles.lists}>
					{selectDropdowns && unMatchLeague && <ListLeague elements={unMatchLeague.filter(key => key.bookmakerName == selectDropdowns?.bookmaker1)} onSelect={handleLeagueSelect1} size={10} right={true} />}
					{selectDropdowns && unMatchLeague && <ListLeague elements={unMatchLeague.filter(key => key.bookmakerName == selectDropdowns?.bookmaker2)} onSelect={handleLeagueSelect2} size={10} />}
				</div>

				<div className={styles.buttonContainer}>
					<button className={styles.button} onClick={addLeaguePairs}>
						Создать пару лиг
					</button>
				</div>
			</div>
		</>
	);
};