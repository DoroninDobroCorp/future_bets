import { GetServerSideProps } from 'next';
import fs from 'fs';
import path from 'path';
import {JSX} from "react";

function LiveMatchDataPage({ htmlContent }: { htmlContent: HTMLElement }): JSX.Element {
    return (
        <div>
            <div dangerouslySetInnerHTML={{ __html: htmlContent }} />
        </div>
    );
}

export const getServerSideProps: GetServerSideProps = async () => {
    // Путь к вашему HTML файлу
    const filePath = path.join(process.cwd(), 'public', 'html', 'live-match-data.html');

    // Чтение файла
    const htmlContent = fs.readFileSync(filePath, 'utf-8');

    return {
        props: {
            htmlContent,
        },
    };
};

export default LiveMatchDataPage;
