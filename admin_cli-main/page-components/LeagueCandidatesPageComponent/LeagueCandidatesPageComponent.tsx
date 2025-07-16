import { JSX, useState } from "react";
import styles from "./LeagueCandidatesPageComponent.module.css"
import { Candidate } from "@/interfaces/Candidate";
import { Header } from "@/components/Header/Header";
import { CheckList } from "@/components/CheckList/CheckList";
import { Dropdown } from "@/components/Dropdown/Dropdown";
import { ErrorMsg } from "@/components/Error/ErrorMsg";


export const LeagueCandidatesPageComponent = ({}): JSX.Element => { 
	let elements: Candidate[] = [
		{
			text1: "Лига1",
			text2: "Лига2",
			similarity: "80"
		},
		{
			text1: "Лига1",
			text2: "Лига2",
			similarity: "80"
		},
		{
			text1: "Лига1",
			text2: "Лига2",
			similarity: "80"
		},
		{
			text1: "Лига1",
			text2: "Лига2",
			similarity: "80"
		},
		{
			text1: "Лига1",
			text2: "Лига2",
			similarity: "80"
		},
		{
			text1: "Лига1",
			text2: "Лига2",
			similarity: "80"
		},
		{
			text1: "Лига1",
			text2: "Лига2",
			similarity: "80"
		},
		{
			text1: "Лига1",
			text2: "Лига2",
			similarity: "80"
		},
		{
			text1: "Лига1",
			text2: "Лига2",
			similarity: "80"
		},
		{
			text1: "Лига1",
			text2: "Лига2",
			similarity: "80"
		}
	];

	let bookmakers: string[] = ["Pinnacle", "Lobbet", "Sansabet"];
	let sports: string[] = ["Football", "Tennis", "Basketball"];

	const [selectedBookmaker1, setSelectedBookmaker1] = useState<string>("");
	const [selectedBookmaker2, setSelectedBookmaker2] = useState<string>("");
	const [selectedSport, setSelectedSport] = useState<string>("");

	const [selectedItems, setSelectedItems] = useState<Candidate[]>([])

	const [errorMessage, setErrorMessage] = useState<string | null>(null);


	const handleSelectedItemsChange = (items: Candidate[]) => {
		setSelectedItems(items);
	};

	const handleBookmaker1Select = (selected: string) => {
		if (selected === selectedBookmaker2) {
			setErrorMessage("Нельзя выбрать одинаковые букмейкеры");
			return;
		}
		setErrorMessage(null);
		setSelectedBookmaker1(selected);
	};

	const handleBookmaker2Select = (selected: string) => {
		if (selected === selectedBookmaker1) {
			setErrorMessage("Нельзя выбрать одинаковые букмейкеры");
			return;
		}
		setErrorMessage(null);
		setSelectedBookmaker2(selected);
	};

	const handleSportSelect = (selected: string) => {
		setSelectedSport(selected);
	}

	return (
		<>
			<Header header="Кандидаты лиг" />

			{errorMessage && (<ErrorMsg message={errorMessage}/>)}

			<div className={styles.dropdownsContainer}>
				<Dropdown elements={bookmakers} onSelect={handleBookmaker1Select}/>
				<Dropdown elements={bookmakers} onSelect={handleBookmaker2Select}/>
				<Dropdown elements={sports} onSelect={handleSportSelect}/>
			</div>

			<CheckList
				label1={selectedBookmaker1 || bookmakers[0]}
				label2={selectedBookmaker2 || bookmakers[0]}
				second_label={selectedSport || sports[0]}
				elements={elements}
				onSelect={handleSelectedItemsChange}
			/>

			<div className={styles.buttonContainer}>
				<button className={styles.button}>
					Подтвердить кандидатов
				</button>
			</div>
		</>
	);
}