import { JSX, useState } from 'react';
import styles from './CheckList.module.css';
import { CheckListProps } from './CheckList.props';

export const CheckList = ({ label1, label2, second_label, elements, onSelect }: CheckListProps): JSX.Element => {
	const [selectedItems, setSelectedItems] = useState<Set<number>>(new Set());

	const handleCheckboxChange = (index: number) => {
		const newSelectedItems = new Set(selectedItems);
		if (selectedItems.has(index)) {
			newSelectedItems.delete(index);
		} else {
			newSelectedItems.add(index);
		}
		setSelectedItems(newSelectedItems);

		if (onSelect) {
			const selectedElements = elements.filter((_, i) => newSelectedItems.has(i));
			onSelect(selectedElements);
		}
	};

	return (
		<ul className={styles.elements}>
			{elements.map((element, index) => (
				<li key={index} className={styles.element}>
					<input
						type="checkbox"
						className={styles.checkbox}
						checked={selectedItems.has(index)}
						onChange={() => handleCheckboxChange(index)}
					/>

					<div className={styles.info}>
						<div>
							<span className={styles.label}>{label1} {second_label}:</span>
							<span className={styles.text}> {element.text1}</span>
							{element.second_text1 && <span className={styles.second_text}>({element.second_text1})</span>}
						</div>

						<div>
							<span className={styles.label}>{label2} {second_label}:</span>
							<span className={styles.text}> {element.text2}</span>
							{element.second_text2 && <span className={styles.second_text}>({element.second_text2})</span>}
						</div>
					</div>

					<div className={styles.similarity}>{element.similarity}%</div>
				</li>
			))}
		</ul>
	);
};